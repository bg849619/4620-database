package main

type MotifInstance struct {
	CellType       string  `json:"CellType" db:"CellType"`
	Chr            string  `json:"Chr" db:"Chr"`
	Start          int     `json:"Start" db:"Start"`
	Forward        bool    `json:"Forward" db:"End"`
	ThresholdScore float64 `json:"ThresholdScore" db:"ThresholdScore"`
}
