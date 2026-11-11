//go:build !windows

package app

func migrateStartupPath() error { return nil }

// Other platforms use MyGo's single-instance lock.
type instanceGuard struct{}

func claimInstance() (*instanceGuard, error) { return &instanceGuard{}, nil }
func (*instanceGuard) watch(func())          {}
func (*instanceGuard) close()                {}
