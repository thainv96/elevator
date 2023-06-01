package usecase

import (
	"fmt"
	"sync"

	"elevator/src/constant"
)

const (
	moveUpStepCount   = 1
	moveDownStepCount = -1
)

type MoveRequestUsecase interface {
	SetRequest(direction constant.MoveDirection, currentLevel int) error
	RemoveRequest(direction constant.MoveDirection, currentLevel int) error
	MaxDistanceToMove(direction constant.MoveDirection, currentLevel int) int
	CheckIfCurrentLevelHasRequestingUsers(direction constant.MoveDirection, currentLevel int) bool
	GetMaxLevel() int
	GetMinLevel() int

	// ExposeLevelData only for testing purpose
	ExposeLevelData() map[constant.MoveDirection][]bool
}

func NewMoveRequestUsecase(maxLevel int, minLevel int) MoveRequestUsecase {
	statusArrayCapacity := maxLevel + 1
	if minLevel < 0 {
		statusArrayCapacity = maxLevel + abs(minLevel) + 1
	}

	requestUpStatusArray := make([]bool, statusArrayCapacity)
	requestDownStatusArray := make([]bool, statusArrayCapacity)
	mapDirectionToStatusArray := make(map[constant.MoveDirection][]bool)
	mapDirectionToStatusArray[constant.MoveDirectionUp] = requestUpStatusArray
	mapDirectionToStatusArray[constant.MoveDirectionDown] = requestDownStatusArray

	return &moveRequestUsecaseImpl{
		maxLevel:                  maxLevel,
		minLevel:                  minLevel,
		requestUpStatusArray:      requestUpStatusArray,
		requestDownStatusArray:    requestDownStatusArray,
		mapDirectionToStatusArray: mapDirectionToStatusArray,
	}
}

type moveRequestUsecaseImpl struct {
	maxLevel                  int
	minLevel                  int
	requestUpStatusArray      []bool
	requestDownStatusArray    []bool
	mapDirectionToStatusArray map[constant.MoveDirection][]bool

	/*
		note for changeLevelStatusMutex:
		can separate mutex for up/down direction for performance.
		but let's trade it off for simplicity as the request rate is low
	*/
	changeLevelStatusMutex sync.Mutex
}

func (c *moveRequestUsecaseImpl) RemoveRequest(direction constant.MoveDirection, currentLevel int) error {
	if !c.isValidLevel(currentLevel) {
		return fmt.Errorf("invalid level")
	}
	c.changeLevelStatusMutex.Lock()
	defer c.changeLevelStatusMutex.Unlock()

	statusArray := c.mapDirectionToStatusArray[direction]
	if statusArray == nil {
		return fmt.Errorf("invalid direction")
	}

	currentLevelIndex := c.getLevelIndex(currentLevel)
	statusArray[currentLevelIndex] = false
	return nil
}

func (c *moveRequestUsecaseImpl) CheckIfCurrentLevelHasRequestingUsers(direction constant.MoveDirection, currentLevel int) bool {
	c.changeLevelStatusMutex.Lock()
	defer c.changeLevelStatusMutex.Unlock()

	statusArray := c.mapDirectionToStatusArray[direction]
	if statusArray == nil {
		return false
	}

	currentLevelIndex := c.getLevelIndex(currentLevel)
	return statusArray[currentLevelIndex]
}

func (c *moveRequestUsecaseImpl) SetRequest(direction constant.MoveDirection, currentLevel int) error {
	if !c.isValidLevel(currentLevel) || !c.isValidDirection(direction, currentLevel) {
		return fmt.Errorf("invalid request")
	}
	c.changeLevelStatusMutex.Lock()
	defer c.changeLevelStatusMutex.Unlock()

	statusArray := c.mapDirectionToStatusArray[direction]
	if statusArray == nil {
		return fmt.Errorf("invalid direction")
	}

	currentLevelIndex := c.getLevelIndex(currentLevel)
	statusArray[currentLevelIndex] = true
	return nil
}

func (c *moveRequestUsecaseImpl) MaxDistanceToMove(direction constant.MoveDirection, currentLevel int) int {
	c.changeLevelStatusMutex.Lock()
	defer c.changeLevelStatusMutex.Unlock()

	currentLevelIndex := c.getLevelIndex(currentLevel)
	if c.mapDirectionToStatusArray[direction] == nil {
		return 0
	}

	maxDistanceToMove := 0
	moveDirectionStep := moveUpStepCount
	if direction == constant.MoveDirectionDown {
		moveDirectionStep = moveDownStepCount
	}

	statusArrayCapacity := len(c.requestUpStatusArray)
	for levelIndex := currentLevelIndex; levelIndex < statusArrayCapacity && levelIndex >= 0; levelIndex = levelIndex + moveDirectionStep {
		if c.requestUpStatusArray[levelIndex] == true || c.requestDownStatusArray[levelIndex] {
			maxDistanceToMove = abs(levelIndex - currentLevelIndex)
		}
	}

	return maxDistanceToMove
}

func (c *moveRequestUsecaseImpl) GetMaxLevel() int {
	return c.maxLevel
}

func (c *moveRequestUsecaseImpl) GetMinLevel() int {
	return c.minLevel
}

// ExposeLevelData fpr testing purpose
func (c *moveRequestUsecaseImpl) ExposeLevelData() map[constant.MoveDirection][]bool {
	return c.mapDirectionToStatusArray
}

func (c *moveRequestUsecaseImpl) getLevelIndex(level int) int {
	levelIndex := level
	if c.minLevel < 0 {
		levelIndex = level + abs(c.minLevel)
	}

	return levelIndex
}

func (c *moveRequestUsecaseImpl) isValidLevel(level int) bool {
	return level >= c.minLevel && level <= c.maxLevel
}

func (c *moveRequestUsecaseImpl) isValidDirection(direction constant.MoveDirection, currentLevel int) bool {
	if direction == constant.MoveDirectionUp {
		return currentLevel < c.maxLevel
	}
	if direction == constant.MoveDirectionDown {
		return currentLevel > c.minLevel
	}
	return direction != constant.MoveDirectionIdle
}

func abs(number int) int {
	if number < 0 {
		return -number
	}

	return number
}
