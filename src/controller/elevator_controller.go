package controller

import (
	"fmt"
	"time"

	"elevator/src/usecase"

	"elevator/src/constant"
	"elevator/src/wrapper"
)

const (
	runnerWorkInterval = time.Second
)

type ElevatorController interface {
	RequestMove(direction constant.MoveDirection, currentLevel int) error
	StartMoveRunner() error
}

func NewElevatorController(minLevel int,
	maxLevel int,
	levelController usecase.MoveRequestUsecase,
	elevator wrapper.Elevator) ElevatorController {
	return &elevatorControlImpl{
		maxLevel:                maxLevel,
		minLevel:                minLevel,
		elevatorMovingDirection: constant.MoveDirectionIdle,
		levelController:         levelController,
		elevator:                elevator,
	}
}

type elevatorControlImpl struct {
	maxLevel                int
	minLevel                int
	elevatorMovingDirection constant.MoveDirection
	levelController         usecase.MoveRequestUsecase
	elevator                wrapper.Elevator
}

func (c *elevatorControlImpl) StartMoveRunner() error {
	go func() {
		for {
			err := c.processMove()
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(runnerWorkInterval)
		}
	}()

	return nil
}

func (c *elevatorControlImpl) RequestMove(direction constant.MoveDirection, currentLevel int) error {
	err := c.levelController.SetRequest(direction, currentLevel)
	fmt.Println(err)
	return err
}

func (c *elevatorControlImpl) processMove() error {
	elevatorCurrentLevel, err := c.elevator.GetCurrentLevel()
	if err != nil {
		return fmt.Errorf("failed to get elevator current level")
	}

	nextMoveDirection := c.pickNextMoveDirection(elevatorCurrentLevel)

	if c.elevatorMovingDirection != constant.MoveDirectionIdle && c.elevatorMovingDirection == reverseDirection(nextMoveDirection) {
		/*
			this block is a workaround because we lack interface for user to choose which level to want to go to
			in reality, the user would choose their level to go. and then they go up or down. and their request will be removed.
		*/
		err := c.levelController.RemoveRequest(c.elevatorMovingDirection, elevatorCurrentLevel)
		if err != nil {
			return err
		}
	}
	c.elevatorMovingDirection = nextMoveDirection
	if c.haveGoingForwardUsers(nextMoveDirection, elevatorCurrentLevel) {
		// now we should proceed opening the door for users
		// c.elevator.OpenDoor()
		// c.elevator.CloseDoor()

		err := c.levelController.RemoveRequest(nextMoveDirection, elevatorCurrentLevel)
		if err != nil {
			return err
		}
	}
	return c.moveElevator(nextMoveDirection)
}

func (c *elevatorControlImpl) pickNextMoveDirection(elevatorCurrentLevel int) constant.MoveDirection {
	nextMoveDirection := constant.MoveDirectionIdle
	if c.elevatorMovingDirection == constant.MoveDirectionIdle {
		maxMoveUpDistance := c.levelController.MaxDistanceToMove(constant.MoveDirectionUp, elevatorCurrentLevel)
		maxMoveDownDistance := c.levelController.MaxDistanceToMove(constant.MoveDirectionDown, elevatorCurrentLevel)

		if maxMoveUpDistance == 0 && maxMoveDownDistance == 0 {
			nextMoveDirection = constant.MoveDirectionIdle
		} else if maxMoveDownDistance != 0 && maxMoveUpDistance == 0 {
			nextMoveDirection = constant.MoveDirectionDown
		} else if maxMoveDownDistance == 0 && maxMoveUpDistance != 0 {
			nextMoveDirection = constant.MoveDirectionUp
		} else if maxMoveDownDistance < maxMoveUpDistance {
			nextMoveDirection = constant.MoveDirectionDown
		} else {
			nextMoveDirection = constant.MoveDirectionUp
		}
	} else {
		if c.levelController.MaxDistanceToMove(c.elevatorMovingDirection, elevatorCurrentLevel) != 0 {
			nextMoveDirection = c.elevatorMovingDirection
		} else {
			reversedDirection := reverseDirection(c.elevatorMovingDirection)
			if c.levelController.MaxDistanceToMove(reversedDirection, elevatorCurrentLevel) != 0 {
				nextMoveDirection = reversedDirection
			}
		}
	}

	return nextMoveDirection
}

func (c *elevatorControlImpl) moveElevator(direction constant.MoveDirection) error {
	if direction == constant.MoveDirectionIdle {
		return nil
	}

	err := c.elevator.Move(direction)
	if err != nil {
		return err
	}

	return nil
}

func (c *elevatorControlImpl) haveGoingForwardUsers(direction constant.MoveDirection, currentLevel int) bool {
	return c.levelController.CheckIfCurrentLevelHasRequestingUsers(direction, currentLevel)
}

func reverseDirection(direction constant.MoveDirection) constant.MoveDirection {
	if direction == constant.MoveDirectionDown {
		return constant.MoveDirectionUp
	}
	if direction == constant.MoveDirectionUp {
		return constant.MoveDirectionDown
	}
	return constant.MoveDirectionIdle
}
