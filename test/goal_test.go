package common_test

import (
	"testing"

	"github.com/offlaneDefender/progress-tracker-go/internal/common"
)

func TestGoal(t *testing.T) {
	t.Run("Types test", func(t *testing.T) {
		g := common.Goal{Name: "test", Progress: 0, MaxTicks: 10, Complete: false}

		if g.Complete || g.Name != "test" || g.Progress != 0 || g.MaxTicks != 10 {
			t.Error("Somehow failed to even create a goal...")
		}
	})
}
