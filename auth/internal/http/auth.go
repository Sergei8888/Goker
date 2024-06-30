package http

import (
	"auth-goker/internal/service"
	"encoding/json"
	"net/http"
)

type authDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthController struct {
	As service.AuthService
}

func (c AuthController) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var signup authDto

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

func (c AuthController) HandleSignin(w http.ResponseWriter, r *http.Request) {
	var signin authDto

	err := json.NewDecoder(r.Body).Decode(&signin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(signin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := c.As.Signin(signin.Email, signin.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "goker_session",
		Value:  session.Token,
		MaxAge: 60 * 60 * 24 * 7,
		Path:   "/",
	})
	w.WriteHeader(http.StatusOK)
}

func (c AuthController) HandleMe(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("goker_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := c.As.GetInfoAboutUserBySessionToken(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(user)
}
