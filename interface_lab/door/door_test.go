package door

import (
	"testing"
)

func TestDoorOpen(t *testing.T) {
	var asset Asset = &fakeDoor{isOpen: false}
	door, ok := asset.(Door)
	if ok {
		door.Open()
		if door.IsOpen() != true {
			t.Error("Fake door expected: opened, real: closed")
		}
	}
}
func TestDoorClose(t *testing.T) {
	var asset Asset = &fakeDoor{isOpen: true}
	door, ok := asset.(Door)
	if ok {
		door.Close()
		if door.IsOpen() {
			t.Error("expected: door is opened, real: door is closed")
		}
	}
}

func TestDoorDaily(t *testing.T) {
	service := AssetService{assets: []Asset{
		&fakeDoor{isOpen: false, isLocked: true},
	}}

	if service.CheckDoorsOpened() != false {
		t.Error("expected: all door are open, real: not all doors are open")
	}

	// Open all doors in the morning
	service.OpenDoors()
	if service.CheckDoorsUnlocked() != true {
		t.Error("expected: all door are unlocked, real: not all doors are unlocked")
	}
	if service.CheckDoorsOpened() != true {
		t.Error("expected: all door are open, real: not all doors are open")
	}
	// Close all doors at night
	service.CloseDoors()
	if service.CheckDoorsClosed() != true {
		t.Error("expected: all door are close, real: not all doors are close")
	}
	if service.CheckDoorsLocked() != true {
		t.Error("expected: all door are locked, real: not all doors are locked")
	}
}
