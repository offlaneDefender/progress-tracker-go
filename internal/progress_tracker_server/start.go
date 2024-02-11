package progress_tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/offlaneDefender/progress-tracker-go/internal/common"
)

type goal = common.Goal

type ProgressTracker struct {
	Goals []goal
}

func (pt *ProgressTracker) AddGoal(name string) {
	pt.Goals = append(pt.Goals, goal{Name: name, Progress: 0})
}

func (pt *ProgressTracker) TickProgress(name string) int {
	index := pt.FindByName(name)
	if index == -1 {
		return -1
	}

	pt.Goals[index].Progress += 10

	return pt.Goals[index].Progress
}

func (pt *ProgressTracker) DeleteGoal(name string) bool {
	index := pt.FindByName(name)
	if index == -1 {
		return false
	}

	fmt.Println(pt.Goals)

	pt.Goals[index] = pt.Goals[len(pt.Goals)-1]
	pt.Goals = pt.Goals[:len(pt.Goals)-1]

	fmt.Println(pt.Goals)

	return true
}

func (pt *ProgressTracker) FindByName(name string) int {
	return slices.IndexFunc(pt.Goals, func(g goal) bool { return g.Name == name })
}

func Start() {
	pt := ProgressTracker{Goals: []goal{}}
	// Initialize a server
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		fmt.Fprintf(w, "Goals: %v", pt.Goals)
	})

	http.HandleFunc("PUT /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		decoder := json.NewDecoder(r.Body)

		var pb PostBody

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

		var pb PostBody

		err := decoder.Decode(&pb)

		if err != nil {
			panic(err)
		}

		pt.AddGoal(pb.Name)

		fmt.Fprintf(w, "Goals: %v", pt.Goals)
	})
	http.HandleFunc("DELETE /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		decoder := json.NewDecoder(r.Body)

		var pb PostBody

		err := decoder.Decode(&pb)

		if err != nil {
			panic(err)
		}

		didDelete := pt.DeleteGoal(pb.Name)

		if didDelete == false {
			fmt.Fprintf(w, "Could not delete %v", pb.Name)
		}

		fmt.Fprintf(w, "Deleted %v", pb.Name)
	})
	http.ListenAndServe(":8080", nil)
}

type PostBody struct {
	Name string `json:"name"`
}
