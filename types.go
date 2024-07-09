package main

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateTaskPayload struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	ProjectID    int64  `json:"projectId"`
	AssignedToID int64  `json:"assignedToId"`
}

type CreateProjectPayload struct {
	Name         string    `json:"name"`
}

type CreateUserPayload struct {
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"password"`
}

type Project struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Task struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	ProjectID    int64     `json:"projectId"`
	AssignedToID int64     `json:"assignedToId"`
	CreatedAt    time.Time `json:"createdAt"`
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}
