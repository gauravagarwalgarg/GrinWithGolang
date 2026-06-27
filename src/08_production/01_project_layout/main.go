/*
Module 8: Production - Standard Go Project Layout

Demonstrates through code:
  - /cmd       → Application entry points (main packages)
  - /internal  → Private packages (compiler-enforced encapsulation)
  - /pkg       → Public libraries safe for external import
  - Dependency injection wiring in main()
  - Interface-based decoupling between layers
  - Clean separation: handler → service → repository

Key insight: Go enforces /internal privacy at the compiler level.
No build tags or access modifiers needed. The directory name IS the access control.

Layout:
  /cmd/api/main.go          ← this file (entry point, DI wiring)
  /internal/user/service.go ← business logic (private to module)
  /internal/user/repo.go    ← data access (private to module)
  /pkg/response/json.go     ← shared utilities (public)

Run: go run main.go
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// --- /pkg/response (public utility) ---

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeResponse(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// --- /internal/user/model.go ---

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// --- /internal/user/repo.go (repository interface + implementation) ---

type UserRepository interface {
	FindByID(id string) (*User, error)
	FindAll() ([]*User, error)
}

type inMemoryUserRepo struct {
	users map[string]*User
}

func newInMemoryUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{
		users: map[string]*User{
			"1": {ID: "1", Name: "Alice", Email: "alice@example.com"},
			"2": {ID: "2", Name: "Bob", Email: "bob@example.com"},
		},
	}
}

func (r *inMemoryUserRepo) FindByID(id string) (*User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %s not found", id)
	}
	return user, nil
}

func (r *inMemoryUserRepo) FindAll() ([]*User, error) {
	result := make([]*User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

// --- /internal/user/service.go (business logic) ---

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id string) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) ListUsers() ([]*User, error) {
	return s.repo.FindAll()
}

// --- /internal/handler/user.go (HTTP handler) ---

type UserHandler struct {
	service *UserService
}

func NewUserHandler(svc *UserService) *UserHandler {
	return &UserHandler{service: svc}
}

func (h *UserHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, APIResponse{Error: err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, APIResponse{Data: users})
}

// --- /cmd/api/main.go (DI wiring) ---

func main() {
	// Dependency injection: wire dependencies bottom-up
	repo := newInMemoryUserRepo()       // data layer
	service := NewUserService(repo)     // business logic
	handler := NewUserHandler(service)  // HTTP layer

	http.HandleFunc("/users", handler.HandleList)

	log.Println("Project layout demo server on :8080")
	log.Println("Try: curl http://localhost:8080/users")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
