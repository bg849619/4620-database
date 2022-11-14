package main

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
