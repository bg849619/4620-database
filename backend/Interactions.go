package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type InteractionParticipation struct {
	Locus       string `json:"Locus" db:"Locus"`
	Interaction string `json:"Interaction" db:"Interaction"`
}

type Interaction struct {
	CellType string `json:"CellType" db:"CellType"`
	ID       int64  `json:"ID" db:"ID"`
}

func (c CellType) GetInteractions() (its []Interaction, err error) {
	query := "SELECT * FROM Interactions WHERE CellType=?"
	rows, err := db.Queryx(query, c.Type)
	if err != nil {
		return []Interaction{}, err
	}
	defer rows.Close()

	its = make([]Interaction, 0)
	for rows.Next() {
		var it Interaction
		rows.StructScan(&it)
		its = append(its, it)
	}

	return
}

func (it Interaction) GetLoci() (loci []Locus, err error) {
	query := "SELECT L.* FROM Loci AS L INNER JOIN InteractionParticipation AS I ON I.Locus=L.ID WHERE I.Interaction=?"
	rows, err := db.Queryx(query, it.ID)
	if err != nil {
		return []Locus{}, err
	}
	defer rows.Close()

	loci = make([]Locus, 0)
	for rows.Next() {
		var loc Locus
		rows.StructScan(&loc)
		loci = append(loci, loc)
	}

	return
}

// TODO: This is a good spot to implement interesting interaction logic.
// TODO: Implement interactions with Nested structs, as a flat struct makes no sense outside the database.

// We won't implement a function to create multiple interactions simultaneously until we have the Nested struct implemented.
func (it Interaction) Create() (newid int64, err error) {
	query := "INSERT INTO Interactions (CellType) VALUES (?)"
	result, err := db.Exec(query, it.CellType)
	if err != nil {
		return 0, err
	}

	newid, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return
}

func (it Interaction) Delete() (err error) {
	query := "DELETE FROM Interactions WHERE ID=?"
	_, err = db.Exec(query, it.ID)
	return
}

func (it Interaction) AddLocus(ID string) (err error) {
	query := "INSERT INTO InteractionParticipation VALUES (?, ?)"
	_, err = db.Exec(query, ID, it.ID)
	return
}

func (it Interaction) RemoveLocus(ID string) (err error) {
	query := "DELETE FROM InteractionParticipation WHERE Locus=? AND Interaction=?"
	_, err = db.Exec(query, ID, it.ID)
	return
}

func GetInteraction(ID int64) (it Interaction, err error) {
	query := "SELECT * FROM Interactions WHERE ID=?"
	rows, err := db.Queryx(query, ID)
	if err != nil {
		return Interaction{}, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.StructScan(&it)
		return
	}

	return Interaction{}, errors.New("not_found")
}

// Creating interactions is handled by the CellType routes.
func handleAddLocusToInteraction(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	key, err := strconv.ParseInt(v["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid ID.")
		return
	}

	it, err := GetInteraction(key)
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Interaction.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Interactions.")
		}
		return
	}

	var locid string
	err = json.NewDecoder(r.Body).Decode(&locid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Body should be a string representing a Locus ID")
		return
	}

	err = it.AddLocus(locid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not add Locus to Interaction.\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleRemoveLocusFromInteraction(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	key, err := strconv.ParseInt(v["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid ID.")
		return
	}

	it, err := GetInteraction(key)
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Interaction.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Interactions.")
		}
		return
	}

	err = it.RemoveLocus(v["loc"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Could not delete Locus from Interaction.")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted Locus from Interaction.")
}

func handleDeleteInteraction(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	key, err := strconv.ParseInt(v["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid ID.")
		return
	}

	it, err := GetInteraction(key)
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Interaction.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Interactions.")
		}
		return
	}

	err = it.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not delete Interaction.")
		return
	}

	fmt.Fprint(w, "Interaction deleted.")
}

func handleGetInteractionLoci(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	key, err := strconv.ParseInt(v["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid ID.")
		return
	}

	it, err := GetInteraction(key)
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Interaction.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Interactions.")
		}
		return
	}

	loc, err := it.GetLoci()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Loci.")
		return
	}

	json.NewEncoder(w).Encode(loc)
}

func init() {
	registerRoute(Route{"/interactions/{id}/loci", handleGetInteractionLoci, "GET"})
	registerRoute(Route{"/interactions/{id}/loci", handleAddLocusToInteraction, "POST"})
	registerRoute(Route{"/interactions/{id}/loci/{loc}", handleRemoveLocusFromInteraction, "DELETE"})
	registerRoute(Route{"/interactions/{id}", handleDeleteInteraction, "DELETE"})
}
