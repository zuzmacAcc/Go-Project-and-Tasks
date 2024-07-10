package main

import (
	"encoding/json"
	"io"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	StatusTODO      = "TODO"
	StatusInProgress = "IN_PROGRESS"
	StatusInTesting  = "IN_TESTING"
	StatusDone       = "DONE"
)

var validStatuses = map[string]bool{
	StatusTODO:      true,
	StatusInProgress: true,
	StatusInTesting:  true,
	StatusDone:       true,
}

type TasksService struct {
	store Store
}

func NewTasksService(s Store) *TasksService {
	return &TasksService{store: s}
}

func (s *TasksService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", WithJWTAuth(s.handleCreateTask, s.store)).Methods("POST")
	r.HandleFunc("/tasks/{id}", WithJWTAuth(s.handleGetTask, s.store)).Methods("GET")
	r.HandleFunc("/tasks/{id}", WithJWTAuth(s.handleDeleteTask, s.store)).Methods("DELETE")
	r.HandleFunc("/tasks/{id}", WithJWTAuth(s.handleEditTask, s.store)).Methods("PUT")
}

func (s *TasksService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	defer r.Body.Close()

	var taskPayload *CreateTaskPayload
	err = json.Unmarshal(body, &taskPayload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateTaskPayload(taskPayload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	t, err := s.store.CreateTask(taskPayload)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error while creating the task"})
		return
	}

	WriteJSON(w, http.StatusCreated, t)
}

func (s *TasksService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	task, err := s.store.GetTask(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "task not found"})
		return
	}

	WriteJSON(w, http.StatusOK, task)
}

func (s *TasksService) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := s.store.DeleteTask(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error deleting task"})
		return
	}

	WriteJSON(w, http.StatusNoContent, nil)
}


func (s *TasksService) handleEditTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	_, err := s.store.GetTask(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("task with ID %s not found", id)})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	defer r.Body.Close()

	var taskPayload *EditTaskPayload
	err = json.Unmarshal(body, &taskPayload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateEditTaskPayload(taskPayload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	t, err := s.store.EditTask(id, taskPayload)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error while updating the task"})
		return
	}

	WriteJSON(w, http.StatusCreated, t)

}

func validateTaskPayload(task *CreateTaskPayload) error {
	if task.Status == "" {
		task.Status = "TODO"
	}

	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}

	return nil
}

func validateEditTaskPayload(task *EditTaskPayload) error {
	if task.Status == "" {
		return errStatusRequired
	}

	if err := validateStatus(task.Status); err != nil {
		return err
	}

	if task.Name == "" {
		return errNameRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}

	return nil
}

func validateStatus(status string) error {
	if !validStatuses[status] {
		return invalidStatus
	}
	return nil
}