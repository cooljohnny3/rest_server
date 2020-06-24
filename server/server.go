package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	// MySQL db driver
	_ "github.com/go-sql-driver/mysql"
)

// User represents user data
type User struct {
	ID       uuid.UUID
	Username string
	Age      int
	Password string
	Email    string
}

// Database is an iterface for dependency injection
type Database interface {
	Insert(User)
	GetUsers() []string
	GetUser(id uuid.UUID) User
	DeleteUser(id uuid.UUID)
}

// Serve starts and serves server on given port
func Serve(port int) {
	db := MySQLStruct{Username: Username, Password: Password}
	http.HandleFunc("/api/add", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			age, err := strconv.Atoi(r.Form.Get("age"))
			if err != nil {
				fmt.Fprintf(w, "error parsing age: %v", err)
				return
			}

			user := User{uuid.New(), r.Form.Get("username"), age, r.Form.Get("password"), r.Form.Get("email")}
			db.Insert(user)

		default:
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Forbidden"))
		}
	})

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, db.GetUsers())

		default:
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Forbidden"))
		}
	})

	http.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			id, err := uuid.Parse(r.Form.Get("id"))
			if err != nil {
				panic(err.Error())
			}

			j, err := json.Marshal(db.GetUser(id))
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(j))

		default:
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Forbidden"))
		}
	})

	http.HandleFunc("/api/remove", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			id, err := uuid.Parse(r.Form.Get("id"))
			if err != nil {
				panic(err.Error())
			}

			db.DeleteUser(id)

		default:
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Forbidden"))
		}
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
