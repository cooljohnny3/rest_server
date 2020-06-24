package server

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// MySQLStruct is a struct for implimenting Database with MySQL
type MySQLStruct struct {
	Username string
	Password string
}

// Insert inserts a user to DB
func (m MySQLStruct) Insert(user User) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/mydb", m.Username, m.Password))
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtIns, err := db.Prepare("INSERT INTO users VALUES(?, ?, ?, ?, ?)") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(user.ID.String(), user.Username, user.Age, user.Password, user.Email)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

// GetUsers gets all users from DB
func (m MySQLStruct) GetUsers() []string {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/mydb", m.Username, m.Password))
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT * FROM users")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	results, err := stmtOut.Query()
	if err != nil {
		panic(err.Error())
	}

	var users []string

	for results.Next() {
		var user User
		err := results.Scan(&user.ID, &user.Username, &user.Age, &user.Password, &user.Email)
		if err != nil {
			panic(err.Error())
		}

		j, err := json.Marshal(user)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, string(j))
	}

	return users
}

// GetUser gets a user from DB by its ID
func (m MySQLStruct) GetUser(id uuid.UUID) User {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/mydb", m.Username, m.Password))
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT * FROM users WHERE user_id=?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	results, err := stmtOut.Query(id)
	if err != nil {
		panic(err.Error())
	}

	var user User
	results.Next()
	if err := results.Scan(&user.ID, &user.Username, &user.Age, &user.Password, &user.Email); err != nil {
		panic(err.Error())
	}

	return user
}

// DeleteUser deletes a user from DB by its ID
func (m MySQLStruct) DeleteUser(id uuid.UUID) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/mydb", m.Username, m.Password))
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("DELETE FROM users WHERE user_id=?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	_, err = stmtOut.Exec(id)
	if err != nil {
		panic(err.Error())
	}
}
