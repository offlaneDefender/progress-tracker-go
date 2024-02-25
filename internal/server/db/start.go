package db

import (
	"database/sql"
	"fmt"
	"log"
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

	res, err := ReadGoals(db)

	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Println(res)

	err = AddGoal(db, "TestInsert", 10)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted goal")

	defer db.Close()
}
