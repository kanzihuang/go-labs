package elevator

import (
	"log"
	"time"
)

type Direction int

const DirectionUp Direction = 1
const DirectionRest Direction = 0
const DirectionDown Direction = -1

type Elevator struct {
	floorNumber      int
	currentFloor     int
	currentDirection Direction
	msPerMoveFloor   time.Duration
	msPerOpenDoor    time.Duration
	floors           []bool
	callChannel      chan int
	terminated       bool
	doorOpen         bool
}

func sumBool(data []bool) (total int) {
	for _, item := range data {
		if item {
			total++
		}
	}
	return
}

func CreateElevator(floorNumber int, currentFloor int, msPerMoveFloor time.Duration, msPerOpenDoor time.Duration) *Elevator {
	elevator := Elevator{
		floorNumber:      floorNumber,
		currentFloor:     currentFloor,
		currentDirection: 0,
		msPerMoveFloor:   msPerMoveFloor,
		msPerOpenDoor:    msPerOpenDoor,
		floors:           make([]bool, floorNumber),
		callChannel:      make(chan int, 10),
	}
	log.Println("Elevator created.")
	go elevator.run()
	return &elevator
}

func (e *Elevator) Terminate() {
	if e == nil {
		log.Println("电梯未创建")
		return
	}
	e.terminated = true
	for sumBool(e.floors) > 0 {
		time.Sleep(time.Millisecond * 100)
	}
	log.Println("电梯关停")
}

func (e *Elevator) run() {
	for {
		<-e.callChannel
		e.Move()
	}
}

func (e *Elevator) getFloorCalling(floor int) bool {
	return e.floors[floor-1]
}

func (e *Elevator) setFloorCalling(floor int, calling bool) {
	e.floors[floor-1] = calling
	if e.currentDirection == DirectionRest {
		if floor > e.currentFloor {
			e.currentDirection = DirectionUp
		} else if floor < e.currentFloor {
			e.currentDirection = DirectionDown
		}
	}
}

func (e *Elevator) Move() {
	for {
		if e.getFloorCalling(e.currentFloor) {
			e.OpenDoor()
			time.Sleep(e.msPerOpenDoor)
			e.CloseDoor()
			e.setFloorCalling(e.currentFloor, false)
		}
		e.CalcDirection()
		if e.currentDirection == DirectionRest {
			break
		}
		time.Sleep(e.msPerMoveFloor)
		e.currentFloor += int(e.currentDirection)
		log.Printf("电梯到达 %d 层", e.currentFloor)
	}
}

func (e *Elevator) callOnFloor(i int) {
	log.Printf("%d 层呼叫电梯\n", i)
	if e.terminated {
		log.Println("==========电梯已关停===========")
		return
	}
	e.setFloorCalling(i, true)
	e.callChannel <- i
}

func (e *Elevator) getCallingFloorNumber(direction Direction) int {
	switch direction {
	case DirectionUp:
		return sumBool(e.floors[e.currentFloor:])
	case DirectionDown:
		return sumBool(e.floors[:e.currentFloor-1])
	default:
		return 0
	}
}

func (e *Elevator) CalcDirection() {
	upFloors := e.getCallingFloorNumber(DirectionUp)
	downFloors := e.getCallingFloorNumber(DirectionDown)

	if e.currentDirection == DirectionRest {
		if upFloors > downFloors {
			log.Println("电梯向上")
			e.currentDirection = DirectionUp
		} else if upFloors+downFloors > 0 {
			log.Println("电梯向下")
			e.currentDirection = DirectionDown
		}
	} else if upFloors+downFloors == 0 {
		log.Println("电梯停止")
		e.currentDirection = DirectionRest
	} else if e.currentDirection == DirectionUp && upFloors == 0 {
		log.Println("电梯向下")
		e.currentDirection = DirectionDown
	} else if e.currentDirection == DirectionDown && downFloors == 0 {
		log.Println("电梯向上")
		e.currentDirection = DirectionUp
	}
}

func (e *Elevator) OpenDoor() {
	log.Println("电梯开门")
	e.doorOpen = true
}

func (e *Elevator) CloseDoor() {
	log.Println("电梯关门")
	e.doorOpen = false
}
