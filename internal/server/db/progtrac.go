package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
)

func FindByName(db *sql.DB, name string) (Goal, error) {
	if name == "" {
		return Goal{}, errors.New("cannot find goal with empty name")
	}

	const query = `SELECT * FROM goals WHERE name = ?`

	row := db.QueryRow(query, name)

	var goal Goal

	if err := row.Scan(&goal.ID, &goal.Name, &goal.MaxTicks, &goal.Progress, &goal.Complete); err == sql.ErrNoRows {
		return Goal{}, err
	}

	return goal, nil
}

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

	_, err = res.LastInsertId()

	if err != nil {
		return err
	}

	return nil
}

func TickProgress(db *sql.DB, name string) (float64, error) {
	if name == "" {
		return 0, errors.New("cannot tick progress for goal with empty name")
	}

	goal, err := FindByName(db, name)

	if err != nil {
		return 0, err
	}

	tickrate := float64(100 / goal.MaxTicks)
	complete := 0

	if goal.Complete || math.Abs(goal.Progress-100) < 0.005 ||
		math.Abs(goal.Progress+tickrate-100) < 0.005 {
		if goal.Progress != 100 {
			goal.Progress = 100
		}
		complete = 1
	} else {
		goal.Progress += tickrate
	}

	stmt, err := db.Prepare(`UPDATE goals SET progress = ?, complete = ? WHERE id = ?`)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(goal.Progress, complete, goal.ID)
	if err != nil {
		return 0, err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		return 0, errors.New("could not tick progress")
	}

	return goal.Progress, nil
}

func DeleteGoal(db *sql.DB, name string) (bool, error) {
	if name == "" {
		return false, errors.New("cannot delete goal with emtpy name")
	}

	goal, err := FindByName(db, name)
	if err != nil {
		return false, err
	}

	stmt, err := db.Prepare(`DELETE FROM goals WHERE id = ?`)
	if err != nil {
		return false, err
	}

	res, err := stmt.Exec(goal.ID)
	if err != nil {
		return false, err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		return false, errors.New("could not delete goal")
	}

	return true, nil
}

func CreateTableIfNotPresent(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS goals(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		maxTicks INTEGER NOT NULL DEFAULT 1,
		progress REAL DEFAULT 0,
		complete INTEGER DEFAULT 0
	);`

	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
