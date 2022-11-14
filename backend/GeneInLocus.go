package main

type GeneInLocus struct {
	Locus string `json:"Locus" db:"Locus"`
	Gene  string `json:"Gene" db:"Gene"`
}

func (g Gene) GetLoci() (result []Locus, err error) {
	query := `SELECT L.* FROM Loci AS L INNER JOIN GeneInLocus AS G ON L.ID=G.Locus WHERE G.Gene=?`
	rows, err := db.Queryx(query, g.Name)
	if err != nil {
		return []Locus{}, err
	}
	defer rows.Close()

	result = make([]Locus, 0)
	for rows.Next() {
		var l Locus
		err = rows.StructScan(&l)
		if err != nil {
			return result, err
		}
		result = append(result, l)
	}

	return
}
