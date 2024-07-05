package main

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateTaskPayload struct {
	Name         string `json:"name"`
	ProjectID    int64  `json:"projectID"`
	AssignedToID int64  `json:"assignedTo"`
}

type Task struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	ProjectID    int64     `json:"projectID"`
	AssignedToID int64     `json:"assignedTo"`
	CreatedAt    time.Time `json:"createdAt"`
}