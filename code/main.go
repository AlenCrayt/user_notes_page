package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type User struct 
{
	username string
	password string
}

type Note struct
{
	note_name string
	content string
}

func main() {
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:los_datos_caidos@localhost:5432/notes_page");
	
	if err != nil {
		fmt.Println("Error connecting to the database");
		return
	}

	//this endpoint is triggered when the user clicks the log-in button in the page that is served in /
	http.HandleFunc("POST /log-in", func(w http.ResponseWriter, r *http.Request) {
		log_in(connection, w, r)
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
	created_user, err := db.Query(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", "", "")
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
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing log in form")
		resp.WriteHeader(500)
		return
	}
	fmt.Println(r.PostFormValue("username") + r.PostFormValue("password"))
	//assuming the problem isn't in the front end then the problem would be as follows:
	//decoder.Decode() isn't correctly putting the values in the struct
	//instead of using the decoder, manually handle the form data by looping over it using http.NewServeMux()?
	var exists bool
	var session_id int
	db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 AND password = $2)", r.PostFormValue("username"), r.PostFormValue("password")).Scan(&exists)
	if exists {
		//store a randomly generated session ID
		//this header does seem to be received by the user
		//resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
		resp.WriteHeader(http.StatusOK)
		//directly causes DOM text to be swapped by it, response successfully received
		io.WriteString(resp, "Message about cookies trying to be sent")
		session_id = rand.Intn(2000000)
		fmt.Println("The session id for the user is: " + strconv.Itoa(session_id))
		//the client doesn't seem to be receiving the cookie or at least the values aren't visible currently
		cookie := http.Cookie {
			Name: "session",
			Value: "455",
			Domain: "192.168.1.29",
			Path: "/notes",
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
	db_note, err := db.Query(context.Background(), "INSERT INTO notes")
	if err != nil {
		fmt.Println("Error querying DB for a new note")
	}
	db_note.Close()
}

func delete_note(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {

}

func modify_note(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {

}

func read_note(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {

}