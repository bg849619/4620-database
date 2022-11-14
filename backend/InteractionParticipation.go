package main

type InteractionParticipation struct {
	Locus       string `json:"Locus" db:"Locus"`
	Interaction string `json:"Interaction" db:"Interaction"`
}
