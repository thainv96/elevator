package test

import (
	"fmt"

	"elevator/src/constant"
	"elevator/src/controller"
	"elevator/src/library"
	"elevator/src/usecase"
	"elevator/src/wrapper"
)

func FullFlowTest() {
	elevatorInitLevel := 0
	mockElevator := library.NewMockElevator(elevatorInitLevel)
	minElevatorLevel := -2
	maxElevatorLevel := 20
	levelController := usecase.NewMoveRequestUsecase(maxElevatorLevel, minElevatorLevel)
	elevatorWrapper := wrapper.NewElevator(mockElevator)
	elevatorController := controller.NewElevatorController(maxElevatorLevel,
		minElevatorLevel,
		levelController,
		elevatorWrapper)

	err := elevatorController.StartMoveRunner()
	if err != nil {
		return
	}

	t := NewHttpTestHandler(mockElevator, levelController)
	go t.Run()

	for {
		fmt.Println("please input current level and direction")
		fmt.Println("example: 12 up")
		fmt.Println("example: 15 down")
		inputLevel := getIntFromConsole()
		inputDirection := getStringFromConsole()

		direction := constant.MoveDirectionUp
		if inputDirection == "up" {
			direction = constant.MoveDirectionUp
		} else if inputDirection == "down" {
			direction = constant.MoveDirectionDown
		} else {
			continue
		}

		elevatorController.RequestMove(direction, inputLevel)
	}
}

func getIntFromConsole() int {
	var number int
	fmt.Scan(&number)
	return number
}

func getStringFromConsole() string {
	var number string
	fmt.Scan(&number)
	return number
}
