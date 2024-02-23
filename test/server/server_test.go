package server_test

import (
	"testing"

	"github.com/offlaneDefender/progress-tracker-go/internal/server"
)

func TestServer(t *testing.T) {
	t.Run("In memory progress tracker", func(t *testing.T) {
		t.Run("Happy cases", func(t *testing.T) {
			var pt server.ProgressTracker

			// test creation
			if len(pt.Goals) != 0 {
				t.Error("Failed to initialize test ProgressTracker.")
			}

			// test appending
			err := pt.AddGoal("testGoal", 10)
			if len(pt.Goals) != 1 || pt.Goals[0].Name != "testGoal" || err != nil {
				t.Error("Failed to add a new Goal.")
			}

			// test finding a Goal by name
			idx := pt.FindByName("testGoal")
			if idx == -1 {
				t.Error("Failed to find a Goal by name.")
			}

			// test deletion
			deleted, err := pt.DeleteGoal("testGoal")
			if !deleted || len(pt.Goals) != 0 || err != nil {
				t.Error("Failed to delete a Goal.")
			}

			// test ticking
			pt.AddGoal("testGoal", 10)
			prog, err := pt.TickProgress("testGoal")
			if prog == -1 || prog != 10 || err != nil {
				t.Error("Failed to tick the progress of a Goal.")
			}
		})
		t.Run("Error cases", func(t *testing.T) {
			var pt server.ProgressTracker

			// test creation
			if len(pt.Goals) != 0 {
				t.Error("Failed to initialize test ProgressTracker.")
			}

			// should fail to add a goal if MaxTicks is zero
			err := pt.AddGoal("testGoal", 0)
			if len(pt.Goals) == 1 || err == nil {
				t.Error("Failed to error adding on invalid input.")
			}

			// should fail to find a goal if no such goal exists
			idx := pt.FindByName("invalid")
			if idx != -1 {
				t.Error("Failed to error finding by non-present Goal name.")
			}

			// should fail to delete a goal if no such goal exists
			pt.AddGoal("testGoal", 10)
			_, err = pt.DeleteGoal("testGoal2")
			if err == nil {
				t.Error("Failed to error on deleting of a non-present Goal.")
			}

			// should fail to tick a goal if no such goal exists
			prog, err := pt.TickProgress("invalid")
			if prog != -1 || prog == 10 || err == nil {
				t.Error("Failed to error ticking the progress of a non-present Goal.")
			}
		})

	})
}
