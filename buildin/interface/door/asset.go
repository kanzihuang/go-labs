package door

import "log"

type Asset interface{}

type AssetService struct {
	assets []Asset
}

func (svc *AssetService) doDoorsAction(action func(door Door)) {
	for _, asset := range svc.assets {
		door, ok := asset.(Door)
		if ok {
			action(door)
		}
	}
}
func (svc *AssetService) OpenDoors() {
	log.Println("Open all door")
	svc.doDoorsAction(func(door Door) {
		door.Unlock()
		door.Open()
	})
}
func (svc *AssetService) CloseDoors() {
	log.Println("Close all door")
	svc.doDoorsAction(func(door Door) {
		door.Close()
		door.Lock()
	})
}

func (svc *AssetService) checkDoors(check func(door Door) bool) bool {
	for _, asset := range svc.assets {
		door, ok := asset.(Door)
		if ok {
			if check(door) != true {
				return false
			}
		}
	}
	return true
}
func (svc *AssetService) CheckDoorsOpened() bool {
	return svc.checkDoors(func(door Door) bool {
		return door.IsOpen()
	})
}

func (svc *AssetService) CheckDoorsClosed() bool {
	return svc.checkDoors(func(door Door) bool {
		return !door.IsOpen()
	})
}
func (svc *AssetService) CheckDoorsLocked() bool {
	return svc.checkDoors(func(door Door) bool {
		return door.IsLocked()
	})
}
func (svc *AssetService) CheckDoorsUnlocked() bool {
	return svc.checkDoors(func(door Door) bool {
		return !door.IsLocked()
	})
}
