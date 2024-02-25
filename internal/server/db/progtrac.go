package db

import (
	"database/sql"
	"fmt"
	"os"
)

func ReadGoals(db *sql.DB) ([]Goal, error) {
	rows, err := db.Query("SELECT * FROM goals")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var goals []Goal

	for rows.Next() {
		var goal Goal

		if err := rows.Scan(&goal.ID, &goal.Name, &goal.Complete, &goal.MaxTicks, &goal.Progress); err != nil {
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
