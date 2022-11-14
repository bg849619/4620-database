package main

type MotifModel struct {
	Name                string `json:"Name" db:"Name"`
	Length              int    `json:"Length" db:"Length"`
	Quality             byte   `json:"Quality" db:"Quality"`
	UniprotID           string `json:"UniprotID" db:"UniprotID"`
	TranscriptionFactor string `json:"TranscriptionFactor" db:"TranscriptionFactor"`
	TFFamily            string `json:"TFFamily" db:"TFFamily"`
	EntrezGene          int    `json:"EntrezGene" db:"EntrezGene"`
}
