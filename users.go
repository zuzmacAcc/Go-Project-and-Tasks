package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store Store
}

func NewUserService(s Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
	r.HandleFunc("/users/login", s.handleUserLogin).Methods("POST")
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var userPayload *CreateUserPayload
	err = json.Unmarshal(body, &userPayload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateUserPayload(userPayload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hashedPassword, err := HashPassword(userPayload.Password)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating user"})
		return
	}
	userPayload.Password = hashedPassword

	u, err := s.store.CreateUser(userPayload)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating user"})
		return
	}

	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating user"})
		return
	}

	WriteJSON(w, http.StatusCreated, token)
}

func (s *UserService) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var loginPayload *LoginUserPayload
	err = json.Unmarshal(body, &loginPayload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid login request payload"})
		return
	}

	if err := validateLoginPayload(loginPayload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// 1. Find user in db by email
	user, err := s.store.GetUserByEmail(loginPayload.Email)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "User doesn't exist"})
		return
	}
	
	// 2. Compare password with hashed password
	if !user.validatePassword(loginPayload.Password){
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Not authenticated"})
		return
	}

	// 3. Create JWT and set it in a cookie
	token, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Not authenticatedr"})
		return
	}

	// 4. Return JWT in response
	WriteJSON(w, http.StatusCreated, token)
}

func validateUserPayload(user *CreateUserPayload) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.FirstName == "" {
		return errFirstNameRequired
	}

	if user.LastName == "" {
		return errLastNameRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func validateLoginPayload(user *LoginUserPayload) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func createAndSetAuthCookie(userID int64, w http.ResponseWriter) (string, error) {
	secret := []byte(Envs.JWTSecret)
	token, err := CreateJWT(secret, userID)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}

func (u *User) validatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}