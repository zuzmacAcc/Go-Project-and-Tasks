package main

import "database/sql"

type Store interface {
	// Users
	CreateUser(u *User) (*User, error)
	GetUserByID(id string) (*User, error)
	//Project
	CreateProject(p *CreateProjectPayload) (*Project, error)
	//Tasks
	CreateTask(t *CreateTaskPayload) (*Task, error)
	GetTask(id string) (*Task, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec("INSERT INTO users (email, firstName, lastName, password) VALUES (?, ?, ?, ?)", u.Email, u.FirstName, u.LastName, u.Password)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

func (s *Storage) GetUserByID(id string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, createdAt FROM users WHERE id = ?", id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)
	return &u, err
}

func (s *Storage) CreateTask(taskPayload *CreateTaskPayload) (*Task, error) {
	rows, err := s.db.Exec("INSERT INTO tasks (name, status, projectId, assignedToId) VALUES (?, ?, ?, ?)", taskPayload.Name, taskPayload.Status, taskPayload.ProjectID, taskPayload.AssignedToID)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	task := &Task{
		ID: id,
		Name: taskPayload.Name,
		Status: taskPayload.Status,
		ProjectID: taskPayload.ProjectID,
		AssignedToID: taskPayload.AssignedToID,
	}
	return task, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task
	err := s.db.QueryRow("SELECT id, name, status, projectId, assignedToId, createdAt FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID, &t.CreatedAt)
	return &t, err
}

func (s *Storage) CreateProject(p *CreateProjectPayload) (*Project, error) {
	rows, err := s.db.Exec("INSERT INTO projects (name) VALUES (?)", p.Name)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	project := &Project {
		ID: id,
		Name: p.Name,
	}

	project.ID = id
	return project, nil
}
