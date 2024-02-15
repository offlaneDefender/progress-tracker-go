package common

type Goal struct {
	Name     string
	Progress int
}

type GoalPostBody struct {
	Name string
}

type GoalPutBody = Goal

type GoalDeleteBody struct {
	Name string
}
