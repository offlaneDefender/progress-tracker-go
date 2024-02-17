package server

import (
	"encoding/json"
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

func (pt *ProgressTracker) AddGoal(name string, maxTicks int) {
	if maxTicks == 0 {
		return
	}
	pt.Goals = append(pt.Goals, goal{Name: name, Progress: 0, MaxTicks: maxTicks})
}

func (pt *ProgressTracker) TickProgress(name string) float64 {
	index := pt.FindByName(name)
	if index == -1 {
		return -1
	}

	foundGoal := pt.Goals[index]

	if foundGoal.MaxTicks == 0 {
		return -1
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
		return 100
	}

	pt.Goals[index].Progress += tickrate

	return pt.Goals[index].Progress
}

func (pt *ProgressTracker) DeleteGoal(name string) bool {
	index := pt.FindByName(name)
	if index == -1 {
		return false
	}

	pt.Goals[index] = pt.Goals[len(pt.Goals)-1]
	pt.Goals = pt.Goals[:len(pt.Goals)-1]

	return true
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

		prog := pt.TickProgress(pb.Name)

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

		if pb.MaxTicks <= 0 {
			http.Error(w, "MaxTicks should not be a non-zero positive integer", 400)
			return
		}

		pt.AddGoal(pb.Name, pb.MaxTicks)

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

		didDelete := pt.DeleteGoal(pb.Name)

		if !didDelete {
			err := fmt.Sprintf("Could not delete %v", pb.Name)
			http.Error(w, err, 400)
			return
		}

		fmt.Fprintf(w, "Deleted %v", pb.Name)
	})
	http.ListenAndServe(":8080", nil)
}
