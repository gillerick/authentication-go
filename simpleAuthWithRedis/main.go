package main

import (
	"gopkg.in/boj/redistore.v1"
	"log"
	"net/http"
	"os"
	"time"
)
import (
	"github.com/gorilla/mux"
)

var store, _ = redistore.NewRediStore(10, "tcp", ":6379", "", []byte(os.Getenv("SESSION_SECRET")))
var users = map[string]string{"naren": "passme", "admin": "password"}

// HealthCheckHandler returns the date and time
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		w.Write([]byte(time.Now().String()))
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

}

// LoginHandler validates the user credentials
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")
	{
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
			return
		}
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")
		if originalPassword, ok := users[username]; ok {
			if originalPassword == password {
				session.Values["authenticated"] = true
				session.Save(r, w)
			} else {
				http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "User is not found", http.StatusNotFound)
			return
		}
		w.Write([]byte("Logged in successfully"))
	}

}

// LogoutHandler removes the session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")
	// Here, instead of setting session value to false, we remove the session
	session.Options.MaxAge = -1
	session.Save(r, w)
	w.Write([]byte(""))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/healthcheck", HealthCheckHandler)
	r.HandleFunc("/logout", LogoutHandler)
	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}
