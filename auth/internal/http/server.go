package http

import "net/http"

func RunServer(authController AuthController) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		authController.HandleSignup(w, r)
	})

	mux.HandleFunc("POST /auth/signin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		authController.HandleSignin(w, r)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}
