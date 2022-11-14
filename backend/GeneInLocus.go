package main

type GeneInLocus struct {
	Locus string `json:"Locus" db:"Locus"`
	Gene  string `json:"Gene" db:"Gene"`
}
