package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// danh sách các API không cần xác thực bằng token
		notAuth := []string{"/api/user/new", "/api/user/login"}
		requestPath := r.URL.Path
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")
		// thiếu jwt token, trả về lỗi
		if tokenHeader == "" {
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		// thông thường chuỗi token có định dạng: Bearer {token-body}, nên cần tách phần token ra
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		// chuỗi jwt token trong phần header của request
		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid {
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		// tiếp tục thực hiện request
		next.ServeHTTP(w, r)
	})
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	// ...
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	// ...
}

func CreateContact(w http.ResponseWriter, r *http.Request) {
	// ...
}

func GetContactsFor(w http.ResponseWriter, r *http.Request) {
	// ...
}

func main() {

	router := mux.NewRouter()

	router.Use(JwtAuthentication) //attach JWT auh middleware

	router.HandleFunc("/api/user/new", CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts", GetContactsFor).Methods("GET") //  user/2/contacts

	handler := c.Handler(router)

	err := http.ListenAndServe(":8080", handler) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
