package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type CellType struct {
	Type string `json:"Type" db:"Type"`
}

/** SQL Helpers **/

func GetCellTypes() ([]CellType, error) {
	query := `SELECT * FROM CellTypes`
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]CellType, 0)

	for rows.Next() {
		var c CellType
		err = rows.StructScan(&c)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}

func GetCellType(t string) (CellType, error) {
	query := "SELECT * FROM CellTypes WHERE Type = ?"
	results, err := db.Queryx(query, t)
	if err != nil {
		return CellType{}, err
	}
	defer results.Close()

	if results.Next() {
		// If we have a result, return the first.
		var t CellType
		err = results.StructScan(&t)
		if err != nil {
			return t, err
		}
		return t, nil
	}

	// No results found.
	return CellType{}, errors.New("not_found")
}

/*
Saves changes to an existing CellType
*/
func (c CellType) Save(oldType string) (err error) {
	query := `UPDATE CellTypes SET Type = ? WHERE Type = ?;`
	_, err = db.Exec(query, c.Type, oldType)
	return
}

func (c CellType) Delete() (err error) {
	query := `DELETE FROM CellTypes WHERE Type = ?`
	_, err = db.Exec(query, c.Type)
	return
}

/** HTTP Routes **/

func handleCreateCellTypes(w http.ResponseWriter, r *http.Request) {
	// Make sure we can decode the endpoint into objects.
	result := make([]CellType, 0)
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Body must be an array of Cell Types.")
		return
	}

	// Store values needed to prepare an SQL statement.
	values := make([]string, len(result))
	args := make([]interface{}, len(result))

	for i := 0; i < len(result); i++ {
		values[i] = "(?)"
		args[i] = result[i].Type
	}

	// Create sql command. (Doing it in this method to filter for SQL injection.)
	query := fmt.Sprintf("INSERT INTO CellTypes VALUES %s", strings.Join(values, ", "))
	statement, err := db.Prepare(query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not create SQL statement.")
		return
	}

	_, err = statement.Exec(args...)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// If no errors until this point, everything was correct.
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Cell types created.")
}

func handleGetCellType(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Cell Type.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Unable to fetch Cell Type.")
		}
		return
	}

	err = json.NewEncoder(w).Encode(cell)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not format cell type for response.")
	}
}

func handleEditCellType(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Cell Type.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Unable to fetch Cell Types.")
		}
		return
	}

	var newCell CellType
	err = json.NewDecoder(r.Body).Decode(&newCell)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Body must be a Cell Type")
		return
	}

	// Update cell with values in body, provided the type of the old.
	err = newCell.Save(cell.Type)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not update Cell Type.")
		fmt.Println(err.Error())
	}
}

func handleDeleteCellType(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Could not find Cell Type.")
		return
	}

	err = cell.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not delete Cell Type.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleGetMotifInstancesInCellType(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Could not find Cell Type.")
		return
	}

	instances, err := cell.GetMotifInstances()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Motif Instances.")
		return
	}

	json.NewEncoder(w).Encode(instances)
}

func handleGetInteractionsInCellType(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Could not find Cell Type.")
		return
	}

	interactions, err := cell.GetInteractions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Interactions")
		return
	}

	json.NewEncoder(w).Encode(interactions)
}

func handleCreateInteraction(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Could not find Cell Type.")
		return
	}

	// No need for JSON, just use the cell type and create an interaction.
	// Function will return the new ID.
	it := Interaction{CellType: cell.Type}
	newid, err := it.Create()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not create Interaction.\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	it.ID = newid
	json.NewEncoder(w).Encode(it)
}

func handleDeleteCellTypeMotif(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Could not find Cell Type.")
		return
	}

	start, err := strconv.ParseInt(v["start"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Start field is invalid. Must be an integer.")
		return
	}

	inst, err := GetMotifInstance(cell.Type, v["chr"], int(start))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Motif Instances.")
		return
	}

	err = inst.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not delete Motif Instance.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleGetGeneExpressions(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Could not find Cell Type.")
		return
	}

	expresesions, err := cell.GeneExpressions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Gene Expressions.")
		return
	}

	json.NewEncoder(w).Encode(expresesions)
}

func handleEditGeneExpression(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Could not find Cell Type.")
		return
	}

	_, err = cell.GeneExpression(v["name"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Gene Expression.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Gene Expressions.")
		}
		return
	}

	var newExpression GeneExpression
	err = json.NewDecoder(r.Body).Decode(&newExpression)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body should be a Gene Expression")
		return
	}

	newExpression.Save()
}

func handleDeleteGeneExpression(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	cell, err := GetCellType(v["type"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Could not find Cell Type.")
		return
	}

	expression, err := cell.GeneExpression(v["name"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Gene Expression.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Gene Expressions.")
		}
		return
	}

	err = expression.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not delete Gene Expression.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func init() {
	registerRoute(Route{"/celltypes", handleGetGeneric(GetCellTypes), "GET"})
	registerRoute(Route{"/celltypes", handleCreateCellTypes, "POST"})
	registerRoute(Route{"/celltypes/{type}", handleGetCellType, "GET"})
	registerRoute(Route{"/celltypes/{type}", handleEditCellType, "PUT"})
	registerRoute(Route{"/celltypes/{type}", handleDeleteCellType, "DELETE"})
	registerRoute(Route{"/celltypes/{type}/motifs", handleGetMotifInstancesInCellType, "GET"})
	registerRoute(Route{"/celltypes/{type}/motifs/{chr}/{start}", handleDeleteCellTypeMotif, "DELETE"})
	registerRoute(Route{"/celltypes/{type}/interactions", handleGetInteractionsInCellType, "GET"})
	registerRoute(Route{"/celltypes/{type}/interactions", handleCreateInteraction, "POST"})
	registerRoute(Route{"/celltypes/{type}/genes", handleGetGeneExpressions, "GET"})
	registerRoute(Route{"/celltypes/{type}/genes/{name}", handleEditGeneExpression, "PUT"})
	registerRoute(Route{"/celltypes/{tpye}/genes/{name}", handleDeleteGeneExpression, "DELETE"})
}
