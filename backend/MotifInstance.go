package main

type MotifInstance struct {
	CellType       string  `json:"CellType" db:"CellType"`
	Chr            string  `json:"Chr" db:"Chr"`
	Start          int     `json:"Start" db:"Start"`
	Forward        bool    `json:"Forward" db:"End"`
	ThresholdScore float64 `json:"ThresholdScore" db:"ThresholdScore"`
	LocusID        string  `json:"LocusID" db:"LocusID"`
	Model          string  `json:"Model" db:"Model"`
}

func (m MotifInstance) GetModel() (model MotifModel, err error) {
	model, err = GetMotifModel(m.Model)
	return
}

func (m MotifModel) GetInstances() (instances []MotifInstance, err error) {
	query := "SELECT * FROM MotifInstance WHERE Model=?"
	rows, err := db.Queryx(query, m.Name)
	if err != nil {
		return []MotifInstance{}, err
	}

	instances = make([]MotifInstance, 0)
	for rows.Next() {
		var mi MotifInstance
		rows.StructScan(&mi)
		instances = append(instances, mi)
	}

	return
}
