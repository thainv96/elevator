package library

import "sync"

type MoveDirection int

const (
	MoveDirectionUp   MoveDirection = 1
	MoveDirectionDown MoveDirection = -1

	moveUpStepCount   = 1
	moveDownStepCount = -1
)

type Elevator interface {
	Move(direction MoveDirection)
	GetCurrentLevel() int
}

func NewMockElevator(initLevel int) Elevator {
	return &mockElevatorImpl{
		currentLevel: initLevel,
	}
}

type mockElevatorImpl struct {
	currentLevel int
	moveMutex    sync.Mutex
}

func (e *mockElevatorImpl) Move(direction MoveDirection) {
	e.moveMutex.Lock()
	defer e.moveMutex.Unlock()

	moveStepCount := moveDownStepCount
	if direction == MoveDirectionUp {
		moveStepCount = moveUpStepCount
	}

	e.currentLevel = e.currentLevel + moveStepCount
}

func (e *mockElevatorImpl) GetCurrentLevel() int {
	e.moveMutex.Lock()
	defer e.moveMutex.Unlock()

	return e.currentLevel
}
