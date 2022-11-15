package main

import "errors"

type GeneExpression struct {
	CellType        string  `json:"CellType" db:"CellType"`
	Gene            string  `json:"Gene" db:"Gene"`
	ExpressionLevel float64 `json:"ExpressionLevel" db:"ExpressionLevel"`
}

func (c CellType) GeneExpressions() ([]GeneExpression, error) {
	query := `SELECT * FROM GeneExpression WHERE CellType=?`
	results, err := db.Queryx(query, c.Type)
	if err != nil {
		return []GeneExpression{}, err
	}
	defer results.Close()

	expressions := make([]GeneExpression, 0)
	for results.Next() {
		var e GeneExpression
		results.StructScan(&e)
		expressions = append(expressions, e)
	}

	return expressions, nil
}

func (c CellType) GeneExpression(name string) (GeneExpression, error) {
	query := `SELECT * FROM GeneExpression WHERE CellType=? AND Gene=?`
	rows, err := db.Queryx(query, c.Type, name)
	if err != nil {
		return GeneExpression{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var exp GeneExpression
		rows.StructScan(&exp)
		return exp, nil
	}

	return GeneExpression{}, errors.New("not_found")
}

func (g Gene) GeneExpressions() ([]GeneExpression, error) {
	query := `SELECT * FROM GeneExpression WHERE Gene=?`
	results, err := db.Queryx(query, g.Name)
	if err != nil {
		return []GeneExpression{}, err
	}
	defer results.Close()

	expressions := make([]GeneExpression, 0)
	for results.Next() {
		var e GeneExpression
		results.StructScan(&e)
		expressions = append(expressions, e)
	}

	return expressions, nil
}

func (g GeneExpression) Save() (err error) {
	query := `UPDATE GeneExpression WHERE CellType=? AND Gene=?`
	_, err = db.Exec(query, g.CellType, g.Gene)
	return
}

func (g GeneExpression) Delete() (err error) {
	query := `DELETE FROM GeneExpression WHERE CellType=? AND Gene=?`
	_, err = db.Exec(query, g.CellType, g.Gene)
	return
}
