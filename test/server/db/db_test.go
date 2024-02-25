package db_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	progtrac "github.com/offlaneDefender/progress-tracker-go/internal/server/db"
)

func TestDb(t *testing.T) {
	t.Run("Happy cases", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "libsql-*")
		if err != nil {
			t.Fatal(err)
		}
		dbPath := dir + "/test.db"

		db, err := sql.Open("libsql", "file:"+dbPath)
		if err != nil {
			t.Errorf("failed to open local db %s", err)
		}

		err = progtrac.CreateTableIfNotPresent(db)

		if err != nil {
			t.Errorf("can't create goals table: %s", err)
		}

		res, err := progtrac.ReadGoals(db)
		if err != nil {
			t.Errorf("%s", err)
		}

		if len(res) != 0 {
			t.Errorf("non-empty goals table")
		}

		testName := "TestInsert"

		_, err = progtrac.FindByName(db, testName)
		if err != nil {
			if err != sql.ErrNoRows {
				t.Errorf("test goal already present %s", err)
			} else {
				err = progtrac.AddGoal(db, testName, 10)
				if err != nil {
					t.Errorf("failed to insert goal: %s", err)
				}
			}
		}

		prg, err := progtrac.TickProgress(db, testName)
		if err != nil {
			t.Errorf("failed to tick progress: %s", err)
		}

		if prg != 10 {
			t.Error("error ticking progress")
		}

		didDelete, err := progtrac.DeleteGoal(db, testName)
		if err != nil {
			t.Errorf("failed to delete goal: %s", err)
		}
		if !didDelete {
			t.Error("error deleting test goal")
		}

		t.Cleanup(func() {
			db.Close()
			defer os.RemoveAll(dir)
		})
	})

}
