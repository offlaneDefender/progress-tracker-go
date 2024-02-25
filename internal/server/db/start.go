package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/offlaneDefender/progress-tracker-go/internal/common"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Goal common.Goal

func Start() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("env: no .env file")
	}
	torsoUrl := os.Getenv("TURSO_URL")
	if torsoUrl == "" {
		log.Fatal("env: TURSO_URL not set")
	}
	torsoToken := os.Getenv("TURSO_AUTH_TOKEN")
	if torsoToken == "" {
		log.Fatal("env: TURSO_AUTH_TOKEN not set")
	}
	dbUrl := torsoUrl + "?authToken=" + torsoToken

	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatalf("failed to open db %s: %s", torsoUrl, err)
	}

	err = CreateTableIfNotPresent(db)
	if err != nil {
		log.Fatalf("can't create goals table: %s", err)
	}

	http.HandleFunc("GET /goals", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		goals, err := ReadGoals(db)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Error reading goals: %s", err), 500)
			return
		}

		fmt.Fprintf(w, "Goals: %v", goals)
	})

	http.HandleFunc("POST /goals", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)
		decoder := json.NewDecoder(r.Body)
		var pb common.GoalPostBody
		err := decoder.Decode(&pb)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding goal: %s", err), 500)
			return
		}
		err = AddGoal(db, pb.Name, pb.MaxTicks)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding goal: %s", err), 500)
			return
		}
		fmt.Fprintf(w, "Added goal %s", pb.Name)
	})

	http.HandleFunc("PUT /goals/{name}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)
		name := r.PathValue("name")
		prg, err := TickProgress(db, name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error ticking %s : %s", name, err), 500)
			return
		}
		fmt.Fprintln(w, prg)
	})

	http.HandleFunc("DELETE /goals/{name}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)
		name := r.PathValue("name")
		_, err := DeleteGoal(db, name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while deleting %s: %s", name, err), 500)
			return
		}
		fmt.Fprintf(w, "Deleted %s\n", name)
	})

	http.ListenAndServe(":8080", nil)

	defer db.Close()
}
