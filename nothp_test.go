package nothp

import "testing"

// TestRelaunchNoLoop asserts Relaunch returns (rather than re-exec'ing)
// once the re-exec sentinel env var is already set.
func TestRelaunchNoLoop(t *testing.T) {
	t.Setenv("NOTHP_RELAUNCHED", "1")
	Relaunch() // must return; a hang or process replacement fails the test run.
}
