//go:build !linux

package nothp

// relaunch is a no-op on platforms without THP/prctl support.
func relaunch() {}
