package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type MotifModel struct {
	Name                string `json:"Name" db:"Name"`
	Length              int    `json:"Length" db:"Length"`
	Quality             byte   `json:"Quality" db:"Quality"`
	UniprotID           string `json:"UniprotID" db:"UniprotID"`
	TranscriptionFactor string `json:"TranscriptionFactor" db:"TranscriptionFactor"`
	TFFamily            string `json:"TFFamily" db:"TFFamily"`
	EntrezGene          int    `json:"EntrezGene" db:"EntrezGene"`
}

func (mm MotifModel) Save(oldName string) (err error) {
	query := "UPDATE MotifModels SET Name=? Length=? Quality=? UniprotID=? TranscriptionFactor=? TFFamily=? EntrezGene=? WHERE Name=?"
	_, err = db.Exec(query, mm.Name, mm.Length, mm.Quality, mm.UniprotID, mm.TranscriptionFactor, mm.TFFamily, mm.EntrezGene, oldName)
	return
}

func (mm MotifModel) Delete() (err error) {
	query := "DELETE FROM MotifModels WHERE Name=?"
	_, err = db.Exec(query, mm.Name)
	return
}

func (mm MotifModel) Create() (err error) {
	query := "INSERT INTO MotifModels VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, mm.Name, mm.Length, mm.Quality, mm.UniprotID, mm.TranscriptionFactor, mm.TFFamily, mm.EntrezGene)
	return
}

func Create(mms []MotifModel) (err error) {
	// We could just call Create() on each, but that's not efficient.
	values := make([]string, len(mms))
	args := make([]interface{}, (len(mms) * 7))
	point := 0
	for i := 0; i < len(mms); i++ {
		values[i] = "(?, ?, ?, ?, ?, ?, ?)"
		args[point] = mms[i].Name
		args[point+1] = mms[i].Length
		args[point+2] = mms[i].Quality
		args[point+3] = mms[i].UniprotID
		args[point+4] = mms[i].TranscriptionFactor
		args[point+5] = mms[i].TFFamily
		args[point+6] = mms[i].EntrezGene
		point += 7
	}

	query := fmt.Sprintf("INSERT INTO MotifModels VALUES %s", strings.Join(values, ", "))
	_, err = db.Exec(query, args...)

	return
}

func GetMotifModels() ([]MotifModel, error) {
	query := "SELECT * FROM MotifModels"
	rows, err := db.Queryx(query)
	if err != nil {
		return []MotifModel{}, err
	}
	defer rows.Close()

	results := make([]MotifModel, 0)
	for rows.Next() {
		var m MotifModel
		rows.StructScan(&m)
		results = append(results, m)
	}

	return results, nil
}

func GetMotifModel(Name string) (MotifModel, error) {
	query := "SELECT * FROM MotifModels WHERE Name=?"
	rows, err := db.Queryx(query, Name)
	if err != nil {
		return MotifModel{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var m MotifModel
		rows.StructScan(&m)
		return m, nil
	}

	return MotifModel{}, errors.New("not_found")
}

func handleGetMotifModels(w http.ResponseWriter, r *http.Request) {
	mms, err := GetMotifModels()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Motif Models.")
		return
	}

	json.NewEncoder(w).Encode(mms)
}

func handleGetMotifModel(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	mm, err := GetMotifModel(v["name"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Motif Model.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Motif Models.")
		}
		return
	}

	json.NewEncoder(w).Encode(mm)
}

func handleEditMotifModel(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	mm, err := GetMotifModel(v["name"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Motif Model.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Motif Models.")
		}
		return
	}

	var newMM MotifModel
	err = json.NewDecoder(r.Body).Decode(&newMM)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body should be a Motif Model.")
		return
	}

	newMM.Save(mm.Name)
}

func handleDeleteMotifModel(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	mm, err := GetMotifModel(v["name"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Motif Model.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Motif Models.")
		}
		return
	}

	err = mm.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not delete Motif Model.\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleCreateMotifModels(w http.ResponseWriter, r *http.Request) {
	var models []MotifModel
	err := json.NewDecoder(r.Body).Decode(&models)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body must be an array of Motif Models.")
		return
	}

	err = Create(models)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not create Motif Models.\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleGetInstancesOfMotif(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	mm, err := GetMotifModel(v["name"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Motif Model.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Motif Models.")
		}
		return
	}

	instances, err := mm.GetInstances()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Motif Instances.")
		return
	}

	json.NewEncoder(w).Encode(instances)
}

func init() {
	registerRoute(Route{"/motifmodels", handleGetMotifModels, "GET"})
	registerRoute(Route{"/motifmodels/{name}", handleGetMotifModel, "GET"})
	registerRoute(Route{"/motifmodels/{name}", handleEditMotifModel, "PUT"})
	registerRoute(Route{"/motifmodels/{name}", handleDeleteMotifModel, "DELETE"})
	registerRoute(Route{"/motifmodels", handleCreateMotifModels, "POST"})
	registerRoute(Route{"/motifmodels/{name}/instances", handleGetInstancesOfMotif, "GET"})
}
