package app

import (
	"errors"
	"log"

	"golang.org/x/sys/windows"
)

// The event is both the lifetime lock and the request to show the existing
// panel. Windows creates it atomically and releases it if the process dies.
type instanceGuard struct {
	event windows.Handle
	stop  chan struct{}
	done  chan struct{}
}

func claimInstance() (*instanceGuard, error) {
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, err
	}
	return claimEvent(`Local\io.github.hunlongyu.gus-gui.instance.` + user.User.Sid.String())
}

func claimEvent(name string) (*instanceGuard, error) {
	key, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	event, err := windows.CreateEventEx(nil, key, 0, windows.SYNCHRONIZE|windows.EVENT_MODIFY_STATE)
	if errors.Is(err, windows.ERROR_ALREADY_EXISTS) && event != 0 {
		defer windows.CloseHandle(event)
		// Pass on the foreground permission of the user's launch, as MyGo
		// does for its own second-instance notification.
		windows.NewLazySystemDLL("user32.dll").NewProc("AllowSetForegroundWindow").Call(0xffffffff)
		return nil, windows.SetEvent(event)
	}
	if err != nil {
		return nil, err
	}
	return &instanceGuard{event: event}, nil
}

func (g *instanceGuard) watch(show func()) {
	g.stop, g.done = make(chan struct{}), make(chan struct{})
	go func() {
		defer close(g.done)
		for {
			if _, err := windows.WaitForSingleObject(g.event, windows.INFINITE); err != nil {
				log.Printf("waiting for another launch: %v", err)
				return
			}
			select {
			case <-g.stop:
				return
			default:
				show()
			}
		}
	}()
}

func (g *instanceGuard) close() {
	if g.event == 0 {
		return
	}
	if g.stop != nil {
		close(g.stop)
		windows.SetEvent(g.event)
		<-g.done
	}
	windows.CloseHandle(g.event)
	g.event = 0
}
