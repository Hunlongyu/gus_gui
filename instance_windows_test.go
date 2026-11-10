package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestInstanceElectionAndActivation(t *testing.T) {
	name := fmt.Sprintf(`Local\gus-test-%d-%s`, os.Getpid(), t.Name())
	var wg sync.WaitGroup
	owners := make(chan *instanceGuard, 16)
	for range cap(owners) {
		wg.Go(func() {
			g, err := claimEvent(name)
			if err != nil {
				t.Errorf("claim: %v", err)
				return
			}
			if g != nil {
				owners <- g
			}
		})
	}
	wg.Wait()
	if len(owners) != 1 {
		t.Fatalf("%d primary instances", len(owners))
	}
	primary := <-owners
	defer primary.close()
	woke := make(chan struct{}, 1)
	primary.watch(func() { woke <- struct{}{} })
	// The concurrent launches happened before the UI was ready. Their
	// activation request must still reach the panel.
	select {
	case <-woke:
	case <-time.After(3 * time.Second):
		t.Fatal("startup activation lost")
	}
	if other, err := claimEvent(name); err != nil || other != nil {
		t.Fatalf("second launch: %v %v", other, err)
	}
	select {
	case <-woke:
	case <-time.After(3 * time.Second):
		t.Fatal("second launch did not activate")
	}
	primary.close()
	next, err := claimEvent(name)
	if err != nil || next == nil {
		t.Fatalf("lock not released for relaunch: %v", err)
	}
	next.close()
}
