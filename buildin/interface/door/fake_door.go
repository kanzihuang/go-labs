package door

import (
	"log"
)

var _ Door = (*fakeDoor)(nil)

type fakeDoor struct {
	isOpen   bool
	isLocked bool
}

func (door *fakeDoor) IsLocked() bool {
	return door.isLocked
}

func (door *fakeDoor) Lock() {
	log.Println("Fake door locked")
	door.isLocked = true
}

func (door *fakeDoor) Unlock() {
	log.Println("Fake door unlocked")
	door.isLocked = false
}

func (door *fakeDoor) IsOpen() bool {
	return door.isOpen
}

func (door *fakeDoor) Open() {
	log.Println("Fake door opened")
	door.isOpen = true
}

func (door *fakeDoor) Close() {
	log.Println("Fake door closed")
	door.isOpen = false
}
