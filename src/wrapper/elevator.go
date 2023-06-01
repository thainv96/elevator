package wrapper

import (
	"sync"

	"elevator/src/constant"
	"elevator/src/library"
)

type Elevator interface {
	Move(direction constant.MoveDirection) error
	GetCurrentLevel() (int, error)
}

func NewElevator(thirdPartyElevator library.Elevator) Elevator {
	return &elevatorImpl{
		thirdPartyElevator: thirdPartyElevator,
	}
}

type elevatorImpl struct {
	thirdPartyElevator library.Elevator
	moveMutex          sync.Mutex
}

func (e *elevatorImpl) Move(direction constant.MoveDirection) error {
	e.moveMutex.Lock()
	defer e.moveMutex.Unlock()

	moveDirection := moveDirectionToThirdPartyDirection(direction)
	e.thirdPartyElevator.Move(moveDirection)

	return nil
}

func (e *elevatorImpl) GetCurrentLevel() (int, error) {
	e.moveMutex.Lock()
	defer e.moveMutex.Unlock()

	return e.thirdPartyElevator.GetCurrentLevel(), nil
}

func moveDirectionToThirdPartyDirection(direction constant.MoveDirection) library.MoveDirection {
	moveDirection := library.MoveDirectionDown
	if direction == constant.MoveDirectionUp {
		moveDirection = library.MoveDirectionUp
	}

	return moveDirection
}
