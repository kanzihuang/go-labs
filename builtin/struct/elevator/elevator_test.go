package elevator

import (
	"testing"
	"time"
)

const msPerMoveFloor = 100 * time.Millisecond
const msPerOpenDoor = 100 * time.Millisecond
const floorNumber int = 5

func waitElevatorReach(e *Elevator, targetFloor int) {
	d := 0 * time.Millisecond
	if e.doorOpen {
		d += e.msPerOpenDoor / 2
	}
	d += time.Duration.Abs(time.Duration(targetFloor-e.currentFloor) * e.msPerMoveFloor)
	d += e.msPerOpenDoor / 2
	time.Sleep(d)
}

func TestElevatorRest(t *testing.T) {
	currentFloor := 3
	elevator := CreateElevator(floorNumber, currentFloor, msPerMoveFloor, msPerOpenDoor)
	defer elevator.Terminate()
	waitElevatorReach(elevator, 1)
	if elevator.currentFloor != currentFloor {
		t.Errorf("没有人请求电梯，电梯应保持静止，不应该从 %d 层移动到 %d 层", currentFloor, elevator.currentFloor)
	}
}

func TestElevatorOneCall(t *testing.T) {
	currentFloor := 1
	elevator := CreateElevator(floorNumber, currentFloor, msPerMoveFloor, msPerOpenDoor)
	defer elevator.Terminate()

	elevator.callOnFloor(3)
	waitElevatorReach(elevator, 3)
	if elevator.currentFloor != 3 {
		t.Errorf("The elevator did not reach the 2th floor in time, currently on the %dth floor.", elevator.currentFloor)
	}
	waitElevatorReach(elevator, 2)
	if elevator.currentFloor != 3 {
		t.Errorf("The elevator did not stay on the 3th floor, currently on the %dth floor.", elevator.currentFloor)
	}
}

func TestElevatorTwoCall(t *testing.T) {
	currentFloor := 3
	elevator := CreateElevator(floorNumber, currentFloor, msPerMoveFloor, msPerOpenDoor)
	defer elevator.Terminate()

	elevator.callOnFloor(4)
	elevator.callOnFloor(2)
	waitElevatorReach(elevator, 4)
	if elevator.currentFloor != 4 {
		t.Errorf("The elevator did not reach the 2th floor in time, currently on the %dth floor.", elevator.currentFloor)
	}
	waitElevatorReach(elevator, 2)
	if elevator.currentFloor != 2 {
		t.Errorf("The elevator did not reach the 2th floor in time, currently on the %dth floor.", elevator.currentFloor)
	}
	waitElevatorReach(elevator, 3)
	if elevator.currentFloor != 2 {
		t.Errorf("The elevator did not stay on the 2th floor, currently on the %dth floor.", elevator.currentFloor)
	}
}

func TestElevatorThreeCall(t *testing.T) {
	currentFloor := 3
	elevator := CreateElevator(floorNumber, currentFloor, msPerMoveFloor, msPerOpenDoor)
	defer elevator.Terminate()

	elevator.callOnFloor(4)
	elevator.callOnFloor(5)
	elevator.callOnFloor(2)
	waitElevatorReach(elevator, 4)
	if elevator.currentFloor != 4 {
		t.Errorf("The elevator did not reach the 2th floor in time, currently on the %dth floor.", elevator.currentFloor)
	}
	waitElevatorReach(elevator, 5)
	if elevator.currentFloor != 5 {
		t.Errorf("The elevator did not reach the 5th floor in time, currently on the %dth floor.", elevator.currentFloor)
	}
	waitElevatorReach(elevator, 2)
	if elevator.currentFloor != 2 {
		t.Errorf("The elevator did not reach the 2th floor in time, currently on the %dth floor.", elevator.currentFloor)
	}
	waitElevatorReach(elevator, 3)
	if elevator.currentFloor != 2 {
		t.Errorf("The elevator did not stay on the 2th floor, currently on the %dth floor.", elevator.currentFloor)
	}
}
