package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"slices"

	"github.com/offlaneDefender/progress-tracker-go/internal/common"
)

type goal = common.Goal

type ProgressTracker struct {
	Goals []goal
}

/*
	TODO:
		1. Add error to return value of functions for unit testing
*/

func (pt *ProgressTracker) AddGoal(name string, maxTicks int) error {
	if maxTicks <= 0 {
		return errors.New("MaxTicks cannot be less than 1")
	}
	pt.Goals = append(pt.Goals, goal{Name: name, Progress: 0, MaxTicks: maxTicks})

	return nil
}

func (pt *ProgressTracker) TickProgress(name string) (float64, error) {
	index := pt.FindByName(name)
	if index == -1 {
		err := fmt.Sprintf("Cannot find goal %v", name)
		return -1, errors.New(err)
	}

	foundGoal := pt.Goals[index]

	if foundGoal.MaxTicks == 0 {
		return -1, errors.New("malformed data, MaxTicks 0")
	}

	tickrate := float64(100 / foundGoal.MaxTicks)

	/*	floating point numbers are fun!
		chose 0.005 as close-enough to zero
		might 'increase' precision ( add more zeros ) later if needed
	*/
	if foundGoal.Complete ||
		math.Abs(foundGoal.Progress-100) < 0.005 ||
		math.Abs(foundGoal.Progress+tickrate-100) < 0.005 {
		if foundGoal.Progress != 100 {
			pt.Goals[index].Progress = 100
		}
		pt.Goals[index].Complete = true
		return 100, nil
	}

	pt.Goals[index].Progress += tickrate

	return pt.Goals[index].Progress, nil
}

func (pt *ProgressTracker) DeleteGoal(name string) (bool, error) {
	index := pt.FindByName(name)
	if index == -1 {
		err := fmt.Sprintf("Cannot find goal %v", name)
		return false, errors.New(err)
	}

	pt.Goals[index] = pt.Goals[len(pt.Goals)-1]
	pt.Goals = pt.Goals[:len(pt.Goals)-1]

	return true, nil
}

func (pt *ProgressTracker) FindByName(name string) int {
	return slices.IndexFunc(pt.Goals, func(g goal) bool { return g.Name == name })
}

func Start() {
	pt := ProgressTracker{Goals: []goal{}}

	// connect to the db
	progtracdb, err := ConnectToDBAndInit()

	if err != nil {
		panic(err)
	}

	id, err := InsertTestData(progtracdb)

	if err != nil {
		panic(err)
	}

	switch id {
	case 0:
		fmt.Println("Data already present, skipping test data insertion")
	default:
		fmt.Printf("Inserted Goal with id %v", id)
	}

	// Initialize a server
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		fmt.Fprintf(w, "Goals: %v", pt.Goals)
	})

	http.HandleFunc("PUT /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		decoder := json.NewDecoder(r.Body)

		var pb common.GoalPutBody

		err := decoder.Decode(&pb)

		if err != nil {
			panic(err)
		}

		prog, err := pt.TickProgress(pb.Name)

		if err != nil {
			http.Error(w, "Failed to tick progress", 500)
			return
		}

		fmt.Fprintf(w, "Goal ticked! Progress: %v", prog)
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		decoder := json.NewDecoder(r.Body)

		var pb common.GoalPostBody

		err := decoder.Decode(&pb)

		if err != nil {
			panic(err)
		}

		err = pt.AddGoal(pb.Name, pb.MaxTicks)

		if err != nil {
			http.Error(w, err.Error(), 400)
		}

		fmt.Fprintf(w, "Goals: %v", pt.Goals)
	})
	http.HandleFunc("DELETE /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		decoder := json.NewDecoder(r.Body)

		var pb common.GoalPostBody

		err := decoder.Decode(&pb)

		if err != nil {
			panic(err)
		}

		didDelete, deleteErr := pt.DeleteGoal(pb.Name)

		if !didDelete || deleteErr != nil {
			err := fmt.Sprintf("Could not delete %v", pb.Name)
			http.Error(w, err, 400)
			return
		}

		fmt.Fprintf(w, "Deleted %v", pb.Name)
	})
	http.ListenAndServe(":8080", nil)
}
