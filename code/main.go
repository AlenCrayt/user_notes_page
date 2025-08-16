package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5"
)

var decoder = schema.NewDecoder()

type User struct {
	username string `schema:"username"`
	password string `schema:"password"`
}

type Note struct
{
	note_name string `schema:"name"`
	content string `schema:"content"`
}

func main() {
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:los_datos_caidos@localhost:5432/notes_page");
	
	if err != nil {
		fmt.Println("Error connecting to the database");
		return
	}

	//this endpoint is triggered when the user clicks the log-in button in the page that is served in /
	http.HandleFunc("/log-in", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			fmt.Println("Received a request of type POST")
			log_in(connection, w, r)
		}
	})

	http.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
		new_user(connection, w, r)
	})

	http.HandleFunc("/new-note", func(w http.ResponseWriter, r *http.Request) {
		new_note(connection, w, r)
	})

	http.HandleFunc("/mod-note", func(w http.ResponseWriter, r *http.Request) {
		modify_note(connection, w, r)
	})

	http.HandleFunc("/delete-note", func(w http.ResponseWriter, r *http.Request) {
		delete_note(connection, w, r)
	})

	http.HandleFunc("/read-note", func(w http.ResponseWriter, r *http.Request) {
		read_note(connection, w, r)
	})
	//initial web page at / is the log in screen
	webpage := http.FileServer(http.Dir("./page"))
	http.Handle("/", webpage)
	http.ListenAndServe(":8080", nil)
}

func new_user(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {
	decoder := schema.NewDecoder()
	var sent_user User
	err := req.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form for new user")
		resp.WriteHeader(500)
		return
	}
	err = decoder.Decode(&sent_user, req.PostForm)
	if err != nil {
		fmt.Println("Error decoding user")
		resp.WriteHeader(500)
		return
	}
	created_user, err := db.Query(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", sent_user.username, sent_user.password)
	if err != nil {
		fmt.Println("Error querying DB for new user")
		resp.WriteHeader(500)
		return
	}
	defer created_user.Close()
	resp.WriteHeader(200)
}

func log_in(db *pgx.Conn, resp http.ResponseWriter, r *http.Request) {
	//need to create a session ID for the user, this lasts for the duration of the session which is refreshed with further actions from the user
	var sent_log_in_data User
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing log in form")
		resp.WriteHeader(500)
		return
	}
	fmt.Println(r.PostForm)
	//assuming the problem isn't in the front end then the problem would be as follows:
	//decoder.Decode() isn't correctly putting the values in the struct
	//instead of using the decoder, manually handle the form data by looping over it using http.NewServeMux()?
	err = decoder.Decode(&sent_log_in_data, r.PostForm)
	if err != nil {
		fmt.Println("Error decoding log in data")
		resp.WriteHeader(500)
	}
	var exists bool
	var session_id int
	//the User struct is empty
	fmt.Println(sent_log_in_data)
	fmt.Println(sent_log_in_data.password)
	db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 AND password = $2)", sent_log_in_data.username, sent_log_in_data.password).Scan(&exists)
	if exists {
		//store a randomly generated session ID
		session_id = rand.Intn(2000000)
		cookie := http.Cookie {
			Name: "session",
			Value: strconv.Itoa(session_id),
			Domain: "192.168.1.29",
			Path: "/",
			MaxAge: 60 * 60,
			HttpOnly: true,
		}
		http.SetCookie(resp, &cookie)
		//don't store it in a database table, use variables in memory instead
		//send it to the user
	} else {
		//send a response saying the user doesn't exist
		fmt.Println("Log in data doesn't correspond to an existing user")
	}
}

func log_out(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {

}

func new_note(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {
	decoder := schema.NewDecoder()
	err := req.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form")
	}
	var note Note
	err = decoder.Decode(&note, req.PostForm)
	if err != nil {
		fmt.Println("Error decoding form")
	}
	db_note, err := db.Query(context.Background(), "INSERT INTO notes")
	if err != nil {
		fmt.Println("Error querying DB for a new note")
	}
	db_note.Close()
	fmt.Println(note.note_name)
	fmt.Println(note.content)
}

func delete_note(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {

}

func modify_note(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {

}

func read_note(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {

}