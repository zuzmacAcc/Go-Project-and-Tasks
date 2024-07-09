package main

// Mocks

type MockStore struct{}

func (s *MockStore) CreateUser(u *CreateUserPayload) (*User, error) {
	return &User{}, nil
}

func (s *MockStore) CreateProject(p *CreateProjectPayload) (*Project, error) {
	return &Project{}, nil
}


func (s *MockStore) GetProject(id string) (*Project, error) {
	return &Project{}, nil
}

func (s *MockStore)GetProjects() ([]*Project, error){
	return []*Project{}, nil
}

func (s *MockStore) CreateTask(t *CreateTaskPayload) (*Task, error) {
	return &Task{}, nil
}

func (s *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}

func (s *MockStore) GetUserByID(id string) (*User, error) {
	return &User{}, nil
}