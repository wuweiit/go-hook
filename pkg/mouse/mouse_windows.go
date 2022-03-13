//go:build windows
// +build windows

package mouse

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"github.com/wuweiit/go-hook/pkg/types"
	"github.com/wuweiit/go-hook/pkg/win32"
)

var hHook struct {
	sync.Mutex
	Pointer uintptr
}

// DefaultHookHandler used when calling keyboard.Register() without passing handler function.
func DefaultHookHandler(c chan<- types.MouseEvent) types.HOOKPROC {
	return func(code int32, wParam, lParam uintptr) uintptr {
		if lParam != 0 {
			c <- types.MouseEvent{
				Message:        types.Message(wParam),
				MSLLHOOKSTRUCT: *(*types.MSLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
			}
		}

		return win32.CallNextHookEx(0, code, wParam, lParam)
	}
}

func install(fn HookHandler, c chan<- types.MouseEvent) error {
	hHook.Lock()
	defer hHook.Unlock()

	if hHook.Pointer != 0 {
		return fmt.Errorf("mouse: hook function is already install")
	}
	if c == nil {
		return fmt.Errorf("mouse: chan must not be nil")
	}
	if fn == nil {
		fn = DefaultHookHandler
	}

	go func() {
		hhk := win32.SetWindowsHookEx(
			types.WH_MOUSE_LL,
			syscall.NewCallback(fn(c)),
			0,
			0)

		if hhk == 0 {
			panic("mouse: failed to install hook function")
		}

		hHook.Pointer = hhk

		var msg *types.MSG

		for {
			if hHook.Pointer == 0 {
				break
			}
			if result := win32.GetMessage(&msg, 0, 0, 0); result != 0 {
				if result < 0 {
					// We don't care what's went wrong, ignore the result value.
					continue
				} else {
					win32.TranslateMessage(&msg)
					win32.DispatchMessage(&msg)
				}
			}
		}
	}()

	return nil
}

func uninstall() error {
	hHook.Lock()
	defer hHook.Unlock()

	if !win32.UnhookWindowsHookEx(hHook.Pointer) {
		return fmt.Errorf("mouse: failed to uninstall hook function")
	}

	hHook.Pointer = 0

	return nil
}
