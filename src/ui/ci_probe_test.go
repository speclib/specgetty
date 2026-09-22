package ui

import "testing"

// A deliberate failure, pushed to verify the Check workflow reports one.
// Removed before this branch is deleted.
func TestCIProbeDeliberateFailure(t *testing.T) {
	t.Fatal("deliberate: verifying that CI reports a failing test")
}
