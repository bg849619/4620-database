package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

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

func (c CellType) GetMotifInstances() (instances []MotifInstance, err error) {
	query := "SELECT * FROM MotifInstance WHERE CellType=?"
	rows, err := db.Queryx(query, c.Type)

	if err != nil {
		return []MotifInstance{}, err
	}
	defer rows.Close()

	instances = make([]MotifInstance, 0)
	for rows.Next() {
		var mi MotifInstance
		rows.StructScan(&mi)
		instances = append(instances, mi)
	}

	return
}

func CreateMotifInstances(instances []MotifInstance) (err error) {
	values := make([]string, len(instances))
	args := make([]interface{}, (len(instances) * 7))
	point := 0
	for i := 0; i < len(instances); i++ {
		values[i] = "(?, ?, ?, ?, ?, ?, ?)"
		args[point] = instances[i].CellType
		args[point+1] = instances[i].Chr
		args[point+2] = instances[i].Start
		args[point+3] = instances[i].Forward
		args[point+4] = instances[i].ThresholdScore
		args[point+5] = instances[i].LocusID
		args[point+6] = instances[i].Model
		point += 7
	}

	query := fmt.Sprintf("INSERT INTO MotifInstances VALUES %s", strings.Join(values, ", "))
	_, err = db.Exec(query, args...)

	return
}

func GetMotifInstance(celltype string, chr string, start int) (instance MotifInstance, err error) {
	query := "SELECT * FROM MotifInstance WHERE CellType=? AND Chr=? AND Start=?"
	rows, err := db.Queryx(query, celltype, chr, start)
	if err != nil {
		return MotifInstance{}, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.StructScan(&instance)
		return
	}

	return MotifInstance{}, errors.New("not_found")
}

func (m MotifInstance) Delete() (err error) {
	query := "DELETE FROM MotifModels WHERE CellType=? AND Chr=? AND Start=?"
	_, err = db.Exec(query, m.CellType, m.Chr, m.Start)
	return
}

func handleCreateMotifInstance(w http.ResponseWriter, r *http.Request) {
	instances := make([]MotifInstance, 0)
	err := json.NewDecoder(r.Body).Decode(&instances)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body should be an array of Motif Instances")
		return
	}

	err = CreateMotifInstances(instances)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not create Motif Instances.")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Created Motif Instances.")
}

func init() {
	registerRoute(Route{"/motifinstances", handleCreateMotifInstance, "POST"})
	// Need to figure out how to implement edits.
}
