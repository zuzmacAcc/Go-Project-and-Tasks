package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type ProjectService struct {
	store Store
}

func NewProjectService(s Store) *ProjectService {
	return &ProjectService{store: s}
}

func (s *ProjectService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/projects", WithJWTAuth(s.handleCreateProject, s.store)).Methods("POST")
	r.HandleFunc("/projects/{id}", WithJWTAuth(s.handleGetProject, s.store)).Methods("GET")
	r.HandleFunc("/projects", WithJWTAuth(s.handleGetProjects, s.store)).Methods("GET")
	// r.HandleFunc("/projects/{id}", WithJWTAuth(s.handleDeleteProject, s.store)).Methods("DELETE")
}

func (s *ProjectService) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	defer r.Body.Close()

	var project *CreateProjectPayload
	err = json.Unmarshal(body, &project)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateProjectPayload(project); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	p, err := s.store.CreateProject(project)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error while creating the project"})
		return
	}

	WriteJSON(w, http.StatusCreated, p)
}

func validateProjectPayload(project *CreateProjectPayload) error {
	if project.Name == "" {
		return errNameRequired
	}

	return nil
}

func (s *ProjectService) handleGetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	project, err := s.store.GetProject(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "project not found"})
		return
	}

	WriteJSON(w, http.StatusOK, project)
}

func (s *ProjectService) handleGetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := s.store.GetProjects()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "projects not found"})
		return
	}

	WriteJSON(w, http.StatusOK, projects)
}