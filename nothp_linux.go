//go:build linux

package nothp

import (
	"log/slog"
	"os"
	"runtime"
	"syscall"

	"golang.org/x/sys/unix"
)

// prSetThpDisable is a prctl(2) option value (see /usr/include/linux/prctl.h),
// safe to hardcode since it's architecture-independent.
const prSetThpDisable = 41

// envRelaunched marks a process as already re-exec'd by Relaunch, to avoid
// an infinite re-exec loop.
const envRelaunched = "NOTHP_RELAUNCHED"

// relaunch disables THP, then re-execs /proc/self/exe with the same
// argv/env plus envRelaunched. On the re-exec'd process it just returns.
func relaunch() {
	if os.Getenv(envRelaunched) == "1" {
		return
	}

	disableTHP()

	// Resolve the real path: exec'ing "/proc/self/exe" literally sets
	// /proc/pid/comm to "exe" instead of the binary's actual name.
	exePath, err := os.Executable()
	if err != nil {
		slog.Warn("nothp: failed to resolve executable path, continuing without re-exec", "err", err)
		return
	}

	env := append(os.Environ(), envRelaunched+"=1")
	if err := syscall.Exec(exePath, os.Args, env); err != nil { //nolint:gosec // argv/env are our own, not attacker-controlled.
		slog.Warn("nothp: failed to re-exec self, continuing without THP disabled", "err", err)
	}
}

// disableTHP disables Transparent Huge Pages via prctl(2). The flag is
// tied to the calling thread, so avoid an OS thread switch before exec.
func disableTHP() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := unix.Prctl(prSetThpDisable, 1, 0, 0, 0); err != nil {
		slog.Warn("nothp: failed to disable transparent huge pages, continuing without it", "err", err)
	}
}
