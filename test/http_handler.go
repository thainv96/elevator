package test

import (
	"log"
	"net/http"

	"elevator/src/usecase"

	"elevator/src/constant"
	"elevator/src/library"
)

type TestHandler struct {
	mockElevator       library.Elevator
	moveRequestUsecase usecase.MoveRequestUsecase
}

func NewHttpTestHandler(mockElevator library.Elevator,
	moveRequestUsecase usecase.MoveRequestUsecase) *TestHandler {
	return &TestHandler{
		mockElevator:       mockElevator,
		moveRequestUsecase: moveRequestUsecase,
	}
}

func (t *TestHandler) Run() {
	http.HandleFunc("/api/view-elevator-status", t.OnViewElevatorRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (h *TestHandler) OnViewElevatorRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	upRequestInfoString := ""
	downRequestInfoString := ""
	elevatorInfoString := ""
	for _, status := range h.moveRequestUsecase.ExposeLevelData()[constant.MoveDirectionUp] {
		if status == false {
			upRequestInfoString = upRequestInfoString + "_"
		} else {
			upRequestInfoString = upRequestInfoString + "x"
		}
	}
	for _, status := range h.moveRequestUsecase.ExposeLevelData()[constant.MoveDirectionDown] {
		if status == false {
			downRequestInfoString = downRequestInfoString + "_"
		} else {
			downRequestInfoString = downRequestInfoString + "x"
		}
	}
	for i := h.moveRequestUsecase.GetMinLevel(); i <= h.moveRequestUsecase.GetMaxLevel(); i++ {
		if h.mockElevator.GetCurrentLevel() == i {
			elevatorInfoString = elevatorInfoString + "x"
		} else {
			elevatorInfoString = elevatorInfoString + "_"
		}
	}
	w.Write([]byte(upRequestInfoString + "up request\n"))
	w.Write([]byte(downRequestInfoString + "down request\n"))
	w.Write([]byte(elevatorInfoString + "elevator\n"))
}
