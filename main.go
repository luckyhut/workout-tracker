package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

var (
	tpl *template.Template
	db  *sql.DB
)

func main() {
	// db setup
	var err error
	db, err = sql.Open("sqlite3", "data/workouts.db")
	if err != nil {
		log.Fatal("Can't open initial DB connection: ", err)
	}
	defer db.Close()
	initWorkoutsTbl()

	// server
	tpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Error parsing index.html")
	}
	sm := http.NewServeMux()
	sm.HandleFunc("/", handlerIndex)
	sm.HandleFunc("POST /workouts", handlerWorkout)
	sm.HandleFunc("POST /reset-db", handlerResetDB)
	http.ListenAndServe(":4000", sm)
}

func initWorkoutsTbl() {
	createTblSQL := `
	CREATE TABLE IF NOT EXISTS workouts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		exercise TEXT NOT NULL,
		weight INTEGER NOT NULL,
		reps INTEGER NOT NULL,
		date DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(createTblSQL)
	if err != nil {
		log.Fatal("Error setting up Workouts DB: ", err)
	}
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func handlerWorkout(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	log.Println("logged:", r.FormValue("exercise"), r.FormValue("weight"), r.FormValue("reps"))
	// send to db
	stmt, err := db.Prepare(
		"INSERT INTO workouts(exercise, weight, reps) VALUES (?, ?, ?)",
	)
	if err != nil {
		log.Println("Error preparing db Statement: ", err)
		return
	}
	defer stmt.Close()
	_, _ = stmt.Exec(
		r.FormValue("exercise"),
		r.FormValue("weight"),
		r.FormValue("reps"),
	)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlerResetDB(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(
		"DELETE FROM workouts",
	)
	if err != nil {
		log.Println("Error preparing delete statement: ", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("Error deleting records from DB: ", err)
	}
	log.Println("Database records deleted succesfully.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
