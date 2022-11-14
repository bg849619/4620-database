package main

type Gene struct {
	Name  string `json:"Name" db:"Name"`
	Chr   string `json:"Chr" db:"Chr"`
	Start string `json:"Start" db:"Start"`
	End   string `json:"End" db:"End"`
}
