package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/offlaneDefender/progress-tracker-go/internal/common"
	"github.com/offlaneDefender/progress-tracker-go/internal/server/repo"
)

type goal = common.Goal

func Start() {
	pt := repo.CreateProgressTracker()

	// Initialize a server
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		fmt.Fprintf(w, "Goals: %v", pt.ReadGoals())
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

		fmt.Fprintf(w, "Goals: %v", pt.ReadGoals())
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
