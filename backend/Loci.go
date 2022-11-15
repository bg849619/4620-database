// TODO: Implement automatic relationships of Genes and Loci upon creation or edit of Loci.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Locus struct {
	ID    string `json:"ID" db:"ID"`
	Chr   string `json:"Chr" db:"Chr"`
	Start int    `json:"Start" db:"Start"`
	End   int    `json:"End" db:"End"`
}

func (l Locus) Save(oldId string) (err error) {
	query := `UPDATE Loci SET ID=?, Chr=?, Start=?, End=? WHERE ID=?`
	_, err = db.Exec(query, l.ID, l.Chr, l.Start, l.End, oldId)
	return
}

func (l Locus) Delete() (err error) {
	query := `DELETE FROM Loci WHERE ID=?`
	_, err = db.Exec(query, l.ID)
	return
}

func GetLoci() ([]Locus, error) {
	query := `SELECT * FROM Loci`
	rows, err := db.Queryx(query)
	if err != nil {
		return []Locus{}, err
	}
	defer rows.Close()

	results := make([]Locus, 0)
	for rows.Next() {
		var l Locus
		rows.StructScan(&l)
		results = append(results, l)
	}

	return results, nil
}

func GetLocus(ID string) (Locus, error) {
	query := `SELECT * FROM Loci WHERE ID=?`
	rows, err := db.Queryx(query, ID)
	if err != nil {
		return Locus{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var l Locus
		rows.StructScan(&l)
		return l, nil
	}

	return Locus{}, errors.New("not_found")
}

func handleGetLoci(w http.ResponseWriter, r *http.Request) {
	loci, err := GetLoci()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Loci.")
		return
	}

	err = json.NewEncoder(w).Encode(loci)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Loci.")
		return
	}
}

func handleGetLocus(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	locus, err := GetLocus(v["id"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Locus.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Loci.")
		}
		return
	}

	err = json.NewEncoder(w).Encode(locus)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Loci.")
	}
}

func handleDeleteLocus(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	locus, err := GetLocus(v["id"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Locus.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Loci.")
		}
		return
	}

	err = locus.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not delete Locus.")
	}
}

func handleCreateLoci(w http.ResponseWriter, r *http.Request) {
	result := make([]Locus, 0)
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body must be an array of Loci.")
		return
	}

	values := make([]string, len(result))
	args := make([]interface{}, (len(result) * 4))
	point := 0
	for i := 0; i < len(result); i++ {
		values[i] = "(?, ?, ?, ?)"
		args[point] = result[i].ID
		args[point+1] = result[i].Chr
		args[point+2] = result[i].Start
		args[point+3] = result[i].End
		point += 4
	}

	query := fmt.Sprintf("INSERT INTO Loci VALUES %s", strings.Join(values, ", "))
	_, err = db.Exec(query, args...)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Could not create Loci.\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Created new Loci.")
}

func handleGetLocusGenes(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	locus, err := GetLocus(v["id"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Locus.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Loci.")
		}
		return
	}

	genes, err := locus.GetGenes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Genes of Locus.")
		return
	}

	json.NewEncoder(w).Encode(genes)
}

func init() {
	registerRoute(Route{"/loci", handleGetLoci, "GET"})
	registerRoute(Route{"/loci/{id}", handleGetLocus, "GET"})
	registerRoute(Route{"/loci/{id}", handleDeleteLocus, "DELETE"})
	registerRoute(Route{"/loci", handleCreateLoci, "POST"})
	registerRoute(Route{"/loci/{id}/genes", handleGetLocusGenes, "GET"})
}
