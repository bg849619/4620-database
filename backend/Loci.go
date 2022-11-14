package main

type Locus struct {
	ID    string `json:"ID" db:"ID"`
	Chr   string `json:"Chr" db:"Chr"`
	Start string `json:"Start" db:"Start"`
	End   string `json:"End" db:"End"`
}
