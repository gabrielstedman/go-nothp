// Package nothp disables Transparent Huge Pages (THP) for the calling
// process before it maps its heap, by re-exec'ing itself if needed.
package nothp

// Relaunch disables Transparent Huge Pages for the current process,
// re-executing the running binary if needed so the setting takes effect
// before the runtime maps its heap. Call this as the first line of main().
func Relaunch() {
	relaunch()
}
