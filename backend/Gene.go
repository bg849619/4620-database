package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Gene struct {
	Name  string `json:"Name" db:"Name"`
	Chr   string `json:"Chr" db:"Chr"`
	Start int    `json:"Start" db:"Start"`
	End   int    `json:"End" db:"End"`
}

func (g Gene) Save(oldName string) (err error) {
	query := `UPDATE Genes SET Name = ?, Chr = ?, Start = ?, End = ? WHERE Name = ?`
	_, err = db.Exec(query, g.Name, g.Chr, g.Start, g.End, oldName)
	return
}

func (g Gene) Delete() (err error) {
	query := `DELETE FROM Genes WHERE Name = ?`
	_, err = db.Exec(query, g.Name)
	return
}

func GetGene(name string) (gene Gene, err error) {
	query := `SELECT * FROM Genes WHERE Name = ?`
	rows, err := db.Queryx(query, name)
	if err != nil {
		return Gene{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var g Gene
		rows.StructScan(&g)
		return g, nil
	}
	// Not found.
	return Gene{}, errors.New("not_found")
}

func GetGenes() (genes []Gene, err error) {
	query := `SELECT * FROM Genes`
	rows, err := db.Queryx(query)
	if err != nil {
		return []Gene{}, err
	}
	defer rows.Close()

	genes = make([]Gene, 0)

	for rows.Next() {
		var g Gene
		rows.StructScan(&g)
		genes = append(genes, g)
	}

	return
}

func handleCreateGenes(w http.ResponseWriter, r *http.Request) {
	result := make([]Gene, 0)
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body must be an array of Genes")
		return
	}

	values := make([]string, len(result))
	args := make([]interface{}, (len(result) * 4))
	point := 0
	for i := 0; i < len(result); i++ {
		values[i] = "(?, ?, ?, ?)"
		args[point] = result[i].Name
		args[point+1] = result[i].Chr
		args[point+2] = result[i].Start
		args[point+3] = result[i].End
		point += 4
	}

	query := fmt.Sprintf("INSERT INTO Genes VALUES %s", strings.Join(values, ", "))
	_, err = db.Exec(query, args...)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		fmt.Fprint(w, "Could not create new Genes.\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Created new Genes.")
}

func handleGetGene(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	gene, err := GetGene(v["name"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Gene.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Gene.")
		}
		return
	}

	err = json.NewEncoder(w).Encode(gene)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not fetch Gene.")
	}
}

func handleEditGene(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	gene, err := GetGene(v["name"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Gene.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Gene.")
		}
		return
	}

	var newGene Gene
	err = json.NewDecoder(r.Body).Decode(&newGene)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body must be a Gene.")
		return
	}

	err = newGene.Save(gene.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not save Gene.")
	}
}

func handleDeleteGene(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	gene, err := GetGene(v["name"])
	if err != nil {
		if err.Error() == "not_found" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Could not find Gene.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch Genes.")
		}
		return
	}

	err = gene.Delete()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not delete Gene.")
	}
}

func init() {
	registerRoute(Route{"/genes", handleGetGeneric(GetGenes), "GET"})
	registerRoute(Route{"/genes", handleCreateGenes, "POST"})
	registerRoute(Route{"/genes/{name}", handleGetGene, "GET"})
	registerRoute(Route{"/genes/{name}", handleEditGene, "PUT"})
	registerRoute(Route{"/genes/{name}", handleDeleteGene, "DELETE"})
}
