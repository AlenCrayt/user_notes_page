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

type User struct 
{
	username string
	password string
}

type Note struct
{
	note_name string `schema:"name"`
	content string `schema:"content"`
}

func main() {
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:los_datos_caidos@192.168.1.29:5432/notes_page");
	
	if err != nil {
		fmt.Println("Error connecting to the database");
		return
	}

	http.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
		new_user(connection, w, r)
	})

	http.HandleFunc("/log-in", func(w http.ResponseWriter, r *http.Request) {
		log_in(connection, w, r)
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
	http.ListenAndServe(":8080", nil)
		var name string;
		fmt.Scan(&name);
		create_row, err := connection.Query(context.Background(), "INSERT INTO users (name) VALUES ($1)", name);
		if err != nil {
			fmt.Println("error querying db")
			return
		}
		create_row.Close();
		//funciona y esta enviando el valor a la base de datos, el problema es que no se esta creando un arreglo o leyendo bien los row retornados por Query()
		returned_row, err := connection.Query(context.Background(), "SELECT name FROM users");
		if err != nil {
			fmt.Println("error querying db with SELECT")
			return
		}
		names, err := pgx.CollectRows(returned_row, pgx.RowTo[string])
		if err != nil {
			fmt.Println("error collecting rows")
			return
		}
		fmt.Println(names);
		returned_row.Close();
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

func log_in(db *pgx.Conn, resp http.ResponseWriter, req *http.Request) {
	//need to create a session ID for the user, this lasts for the duration of the session which is refreshed with further actions from the user
	decoder := schema.NewDecoder()
	var sent_log_in_data User
	err := req.ParseForm()
	if err != nil {
		fmt.Println("Error parsing log in form")
		resp.WriteHeader(500)
		return
	}
	err = decoder.Decode(&sent_log_in_data, req.PostForm)
	if err != nil {
		fmt.Println("Error decoding log in data")
		resp.WriteHeader(500)
		return
	}
	var exists bool
	db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 AND password = $2)", sent_log_in_data.username, sent_log_in_data.password).Scan(&exists)
	if exists {
		//create a session ID
		session_id := rand.Intn(2000000)
		cookie := http.Cookie {
			Name: "session",
			Value: strconv.Itoa(session_id),
			Domain: "192.168.1.29",
			Path: "/",
			MaxAge: 60 * 60,
			HttpOnly: true,
		}
		http.SetCookie(resp, &cookie)
		var usid int
		db.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1 AND password = $2", sent_log_in_data.username, sent_log_in_data.password).Scan(&usid)
		//store it in a database table?
		db.QueryRow(context.Background(), "INSERT INTO sessions (id, user_id) VALUES ($1, $2)", session_id, usid)
		//send it to the user
	} else {
		//send a response saying the user doesn't exist
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