package http

import (
	"auth-goker/internal/service"
	"encoding/json"
	"net/http"
)

type signupDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthController struct {
	As service.AuthService
}

func (c AuthController) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var signup signupDto

	err := json.NewDecoder(r.Body).Decode(&signup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(signup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.As.Signup(signup.Email, signup.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c AuthController) HandleSignin(w http.ResponseWriter, r *http.Request) {}

func (c AuthController) HandleSessionValidation(w http.ResponseWriter, r *http.Request) {}
