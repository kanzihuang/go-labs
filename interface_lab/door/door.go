package door

type Door interface {
	IsOpen() bool
	Open()
	Close()
	IsLocked() bool
	Lock()
	Unlock()
}
