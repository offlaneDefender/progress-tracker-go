package repo

import (
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/offlaneDefender/progress-tracker-go/internal/common"
)

type goal common.Goal

type ProgressTracker struct {
	goals []goal
}

func CreateProgressTracker() *ProgressTracker {
	return &ProgressTracker{goals: []goal{}}
}

func (pt *ProgressTracker) ReadGoals() []goal {
	return pt.goals
}

func (pt *ProgressTracker) AddGoal(name string, maxTicks int) error {
	if maxTicks <= 0 {
		return errors.New("MaxTicks cannot be less than 1")
	}
	pt.goals = append(pt.goals, goal{Name: name, Progress: 0, MaxTicks: maxTicks})

	return nil
}

func (pt *ProgressTracker) TickProgress(name string) (float64, error) {
	index := pt.FindByName(name)
	if index == -1 {
		err := fmt.Sprintf("Cannot find goal %v", name)
		return -1, errors.New(err)
	}

	foundGoal := pt.goals[index]

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
			pt.goals[index].Progress = 100
		}
		pt.goals[index].Complete = true
		return 100, nil
	}

	pt.goals[index].Progress += tickrate

	return pt.goals[index].Progress, nil
}

func (pt *ProgressTracker) DeleteGoal(name string) (bool, error) {
	index := pt.FindByName(name)
	if index == -1 {
		err := fmt.Sprintf("Cannot find goal %v", name)
		return false, errors.New(err)
	}

	pt.goals[index] = pt.goals[len(pt.goals)-1]
	pt.goals = pt.goals[:len(pt.goals)-1]

	return true, nil
}

func (pt *ProgressTracker) FindByName(name string) int {
	return slices.IndexFunc(pt.goals, func(g goal) bool { return g.Name == name })
}
