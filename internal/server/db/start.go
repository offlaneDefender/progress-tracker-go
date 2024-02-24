package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/offlaneDefender/progress-tracker-go/internal/common"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Goal common.Goal

func Start() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	dbUrl := os.Getenv("TURSO_URL")
	if dbUrl == "" {
		log.Fatal("env: TURSO_URL not set")
	}

	dbAuthToken := os.Getenv("TURSO_AUTH_TOKEN")
	if dbAuthToken == "" {
		log.Fatal("env: TUROS_AUTH_TOKEN not set")
	}

	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatalf("failed to open db %s: %s", dbUrl, err)
	}

	res, err := ReadGoals(db)

	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Println(res)

	defer db.Close()
}
