//go:build !windows
// +build !windows

package mouse

import (
	"fmt"

	"github.com/wuweiit/go-hook/pkg/types"
)

func install(fn HookHandler, c chan<- types.MouseEvent) error {
	return fmt.Errorf("mouse: not supported")
}

func uninstall() error {
	return fmt.Errorf("mouse: not supported")
}
