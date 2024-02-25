package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func ReadGoals(db *sql.DB) ([]Goal, error) {
	rows, err := db.Query("SELECT * FROM goals")
	if err != nil {
		log.Fatalf("failed to execute query: %v\n", err)
	}
	defer rows.Close()

	var goals []Goal

	for rows.Next() {
		var goal Goal

		if err := rows.Scan(&goal.ID, &goal.Name, &goal.MaxTicks, &goal.Progress, &goal.Complete); err != nil {
			fmt.Println("Error scanning row:", err)
			return make([]Goal, 0), err
		}

		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return make([]Goal, 0), err
	}

	return goals, nil
}

func AddGoal(db *sql.DB, name string, maxTicks int) error {
	if name == "" {
		return errors.New("cannot insert goal with empty name")
	}

	query := `INSERT INTO goals(name, maxTicks, complete, progress) VALUES(
		?,
		?,
		0,
		0
	);`

	res, err := db.Exec(query, name, maxTicks)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	fmt.Printf("Inserted goal with id %d", id)

	return nil
}

func CreateTableIfNotPresent(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS goals(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		maxTicks INTEGER NOT NULL,
		progress REAL,
		complete INTEGER
	);`

	_, err := db.Query(query)

	if err != nil {
		return err
	}

	return nil
}
