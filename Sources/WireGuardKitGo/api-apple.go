/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2018-2019 Jason A. Donenfeld <Jason@zx2c4.com>. All Rights Reserved.
 */

package main

// #include <stdlib.h>
// #include <sys/types.h>
// static void callLogger(void *func, void *ctx, int level, const char *msg)
// {
// 	((void(*)(void *, int, const char *))func)(ctx, level, msg);
// }
// static void callPacketCallback(void *func, void *ctx, const void *buf, int len)
// {
// 	((void(*)(void *, const void *, int))func)(ctx, buf, len);
// }
import "C"

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
)

var loggerFunc unsafe.Pointer
var loggerCtx unsafe.Pointer

// ChannelTun implements tun.Device using Go channels instead of file descriptors.
// Swift pushes packets in via wgReceivePacket (→ inbound channel → RoutineReadFromTUN).
// Go pushes packets out via a Swift callback (RoutineSequentialReceiver → tun.Write → callback).
type ChannelTun struct {
	inbound  chan []byte // Swift → Go (packets from the app, outbound to VPN)
	closed   chan struct{}
	events   chan tun.Event
	mtu      int

	callbackFunc unsafe.Pointer // Swift function pointer for inbound packets (Go → Swift)
	callbackCtx  unsafe.Pointer // Swift context for the callback
	callbackMu   sync.RWMutex
}

// tunnelHandles maps handle IDs to running tunnels.
var tunnelHandles = make(map[int32]tunnelHandle)

// channelTuns maps handle IDs to their ChannelTun, so wgReceivePacket can find them.
var channelTuns = make(map[int32]*ChannelTun)

func CreateChannelTun() *ChannelTun {
	t := &ChannelTun{
		inbound: make(chan []byte, 1024),
		closed:  make(chan struct{}),
		events:  make(chan tun.Event, 10),
		mtu:     1420,
	}
	t.events <- tun.EventUp
	return t
}

func (t *ChannelTun) File() *os.File          { return nil }
func (t *ChannelTun) Name() (string, error)   { return "channel0", nil }
func (t *ChannelTun) Events() <-chan tun.Event { return t.events }
func (t *ChannelTun) MTU() (int, error)       { return t.mtu, nil }
func (t *ChannelTun) Flush() error            { return nil }

func (t *ChannelTun) Close() error {
	select {
	case <-t.closed:
		// already closed
	default:
		close(t.closed)
	}
	close(t.events)
	return nil
}

// Read blocks until a packet is available from Swift (via wgReceivePacket).
// Called by WireGuard's RoutineReadFromTUN.
func (t *ChannelTun) Read(buf []byte, offset int) (int, error) {
	select {
	case <-t.closed:
		return 0, errors.New("channel tun closed")
	case pkt := <-t.inbound:
		n := copy(buf[offset:], pkt)
		return n, nil
	}
}

// Write delivers a decrypted packet from WireGuard back to Swift via callback.
// Called by WireGuard's RoutineSequentialReceiver.
func (t *ChannelTun) Write(buf []byte, offset int) (int, error) {
	t.callbackMu.RLock()
	fn := t.callbackFunc
	ctx := t.callbackCtx
	t.callbackMu.RUnlock()

	if uintptr(fn) == 0 {
		return 0, errors.New("no packet callback registered")
	}

	pkt := buf[offset:]
	if len(pkt) == 0 {
		return 0, nil
	}

	C.callPacketCallback(fn, ctx, unsafe.Pointer(&pkt[0]), C.int(len(pkt)))
	return len(pkt), nil
}

type CLogger int

func cstring(s string) *C.char {
	b, err := unix.BytePtrFromString(s)
	if err != nil {
		b := [1]C.char{}
		return &b[0]
	}
	return (*C.char)(unsafe.Pointer(b))
}

func (l CLogger) Printf(format string, args ...interface{}) {
	if uintptr(loggerFunc) == 0 {
		return
	}
	C.callLogger(loggerFunc, loggerCtx, C.int(l), cstring(fmt.Sprintf(format, args...)))
}

type tunnelHandle struct {
	*device.Device
	*device.Logger
}

func init() {
	signals := make(chan os.Signal)
	signal.Notify(signals, unix.SIGUSR2)
	go func() {
		buf := make([]byte, os.Getpagesize())
		for {
			select {
			case <-signals:
				n := runtime.Stack(buf, true)
				buf[n] = 0
				if uintptr(loggerFunc) != 0 {
					C.callLogger(loggerFunc, loggerCtx, 0, (*C.char)(unsafe.Pointer(&buf[0])))
				}
			}
		}
	}()
}

//export wgSetLogger
func wgSetLogger(context, loggerFn uintptr) {
	loggerCtx = unsafe.Pointer(context)
	loggerFunc = unsafe.Pointer(loggerFn)
}

//export wgTurnOn
func wgTurnOn(settings *C.char, tunFd int32) int32 {
	logger := &device.Logger{
		Verbosef: CLogger(0).Printf,
		Errorf:   CLogger(1).Printf,
	}

	tunDev := CreateChannelTun()
	logger.Verbosef("Attaching to interface (channel tun)")
	dev := device.NewDevice(tunDev, conn.NewStdNetBind(), logger)

	err := dev.IpcSet(C.GoString(settings))
	if err != nil {
		logger.Errorf("Unable to set IPC settings: %v", err)
		return -1
	}

	dev.Up()
	logger.Verbosef("Device started")

	var i int32
	for i = 0; i < math.MaxInt32; i++ {
		if _, exists := tunnelHandles[i]; !exists {
			break
		}
	}
	if i == math.MaxInt32 {
		return -1
	}
	tunnelHandles[i] = tunnelHandle{dev, logger}
	channelTuns[i] = tunDev
	return i
}

//export wgTurnOff
func wgTurnOff(tunnelHandle int32) {
	dev, ok := tunnelHandles[tunnelHandle]
	if !ok {
		return
	}
	delete(tunnelHandles, tunnelHandle)
	delete(channelTuns, tunnelHandle)
	dev.Close()
}

//export wgSetConfig
func wgSetConfig(tunnelHandle int32, settings *C.char) int64 {
	dev, ok := tunnelHandles[tunnelHandle]
	if !ok {
		return 0
	}
	err := dev.IpcSet(C.GoString(settings))
	if err != nil {
		dev.Errorf("Unable to set IPC settings: %v", err)
		if ipcErr, ok := err.(*device.IPCError); ok {
			return ipcErr.ErrorCode()
		}
		return -1
	}
	return 0
}

//export wgGetConfig
func wgGetConfig(tunnelHandle int32) *C.char {
	device, ok := tunnelHandles[tunnelHandle]
	if !ok {
		return nil
	}
	settings, err := device.IpcGet()
	if err != nil {
		return nil
	}
	return C.CString(settings)
}

//export wgBumpSockets
func wgBumpSockets(tunnelHandle int32) {
	dev, ok := tunnelHandles[tunnelHandle]
	if !ok {
		return
	}
	go func() {
		for i := 0; i < 10; i++ {
			err := dev.BindUpdate()
			if err == nil {
				dev.SendKeepalivesToPeersWithCurrentKeypair()
				return
			}
			dev.Errorf("Unable to update bind, try %d: %v", i+1, err)
			time.Sleep(time.Second / 2)
		}
		dev.Errorf("Gave up trying to update bind; tunnel is likely dysfunctional")
	}()
}

//export wgDisableSomeRoamingForBrokenMobileSemantics
func wgDisableSomeRoamingForBrokenMobileSemantics(tunnelHandle int32) {
	dev, ok := tunnelHandles[tunnelHandle]
	if !ok {
		return
	}
	dev.DisableSomeRoamingForBrokenMobileSemantics()
}

// wgReceivePacket is called by Swift to push a packet into WireGuard.
// The packet data is copied into Go memory before returning.
// Returns 0 on success, -1 if the tunnel is not found or closed.
//
//export wgReceivePacket
func wgReceivePacket(handle int32, buf unsafe.Pointer, pktLen int32) int32 {
	tunDev, ok := channelTuns[handle]
	if !ok {
		return -1
	}

	// Copy the packet from Swift memory into Go-managed memory
	goPacket := make([]byte, pktLen)
	copy(goPacket, (*[1 << 30]byte)(buf)[:pktLen:pktLen])

	select {
	case <-tunDev.closed:
		return -1
	case tunDev.inbound <- goPacket:
		return 0
	}
}

// wgSetPacketCallback registers a Swift callback that receives packets from WireGuard.
// The callback signature is: void callback(void *ctx, const void *buf, int len)
// Called from Go's RoutineSequentialReceiver goroutine — the callback must be thread-safe.
//
//export wgSetPacketCallback
func wgSetPacketCallback(handle int32, ctx unsafe.Pointer, fn unsafe.Pointer) {
	tunDev, ok := channelTuns[handle]
	if !ok {
		return
	}
	tunDev.callbackMu.Lock()
	tunDev.callbackCtx = ctx
	tunDev.callbackFunc = fn
	tunDev.callbackMu.Unlock()
}

//export wgVersion
func wgVersion() *C.char {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return C.CString("unknown")
	}
	for _, dep := range info.Deps {
		if dep.Path == "golang.zx2c4.com/wireguard" {
			parts := strings.Split(dep.Version, "-")
			if len(parts) == 3 && len(parts[2]) == 12 {
				return C.CString(parts[2][:7])
			}
			return C.CString(dep.Version)
		}
	}
	return C.CString("unknown")
}

func main() {}
