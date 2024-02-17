package common

type Goal struct {
	Name     string
	Progress float64
	MaxTicks int
	Complete bool
}

type GoalPostBody struct {
	Name     string
	MaxTicks int
}

type GoalPutBody = Goal

type GoalDeleteBody struct {
	Name string
}
