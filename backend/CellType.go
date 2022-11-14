package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not decode cell types.")
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

func init() {
	registerRoute(Route{"/celltypes", handleGetGeneric(GetCellTypes), "GET"})
	registerRoute(Route{"/celltypes", handleCreateCellTypes, "POST"})
	registerRoute(Route{"/celltypes/{type}", handleGetCellType, "GET"})
	registerRoute(Route{"/celltypes/{type}", handleEditCellType, "PUT"})
	registerRoute(Route{"/celltypes/{type}", handleDeleteCellType, "DELETE"})
}
