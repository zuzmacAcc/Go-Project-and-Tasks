package main

import "database/sql"

type Store interface {
	// Users
	CreateUser(u *CreateUserPayload) (*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	//Project
	CreateProject(p *CreateProjectPayload) (*Project, error)
	GetProject(id string) (*Project, error)
	GetProjects() ([]*Project, error)
	DeleteProject(id string) error
	//Tasks
	CreateTask(t *CreateTaskPayload) (*Task, error)
	GetTask(id string) (*Task, error)
	DeleteTask(id string) error
	EditTask(id string, t *EditTaskPayload) (*Task, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(userPayload *CreateUserPayload) (*User, error) {
	rows, err := s.db.Exec("INSERT INTO users (email, firstName, lastName, password) VALUES (?, ?, ?, ?)", userPayload.Email, userPayload.FirstName, userPayload.LastName, userPayload.Password)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        id,
		Email:     userPayload.Email,
		FirstName: userPayload.FirstName,
		LastName:  userPayload.LastName,
		Password:  userPayload.Password,
	}
	return user, nil
}

func (s *Storage) GetUserByID(id string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, createdAt FROM users WHERE id = ?", id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)
	return &u, err
}

func (s *Storage) GetUserByEmail(email string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, password, createdAt FROM users WHERE email = ?", email).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Password, &u.CreatedAt)
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
		ID:           id,
		Name:         taskPayload.Name,
		Status:       taskPayload.Status,
		ProjectID:    taskPayload.ProjectID,
		AssignedToID: taskPayload.AssignedToID,
	}
	return task, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task
	err := s.db.QueryRow("SELECT id, name, status, projectId, assignedToId, createdAt FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID, &t.CreatedAt)
	return &t, err
}

func (s *Storage) DeleteTask(id string) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
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

	project := &Project{
		ID:   id,
		Name: p.Name,
	}

	project.ID = id
	return project, nil
}

func (s *Storage) GetProject(id string) (*Project, error) {
	var p Project
	err := s.db.QueryRow("SELECT id, name, createdAt FROM projects WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.CreatedAt)
	return &p, err
}

func (s *Storage) GetProjects() ([]*Project, error) {
	rows, err := s.db.Query("SELECT id, name, createdAt FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []*Project{}

	for rows.Next() {
		// p := new(Project)
		var p Project
		err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		// projects = append(projects, p)
		projects = append(projects, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *Storage) DeleteProject(id string) error {
	_, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) EditTask(id string, t *EditTaskPayload) (*Task, error) {
	query := "UPDATE tasks SET name = ?, status = ?, AssignedToID = ? WHERE id = ?"

	_, err := s.db.Exec(query, t.Name, t.Status, t.AssignedToID, id)
	if err != nil {
		return nil, err
	}

	var updatedTask Task
	err = s.db.QueryRow("SELECT id, name, status, AssignedToID, createdAt FROM tasks WHERE id = ?", id).Scan(
		&updatedTask.ID,
		&updatedTask.Name,
		&updatedTask.Status,
		&updatedTask.AssignedToID,
		&updatedTask.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedTask, nil
}
