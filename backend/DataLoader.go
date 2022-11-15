package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func LocusExist(ID string) bool {
	_, err := GetLocus(ID)
	return (err == nil)
}

func GeneExists(Name string) bool {
	_, err := GetGene(Name)
	return (err == nil)
}

func LocusFromID(ID string) Locus {
	colonPos := strings.Index(ID, ":")
	dashPos := strings.Index(ID, "-")

	chr := ID[:colonPos-1]
	startStr := ID[colonPos+1 : dashPos-1]
	endStr := ID[dashPos+1:]

	start, err := strconv.ParseInt(startStr, 10, 32)
	if err != nil {
		log.Panic(err)
	}
	end, err := strconv.ParseInt(endStr, 10, 32)
	if err != nil {
		log.Panic(err)
	}

	return Locus{
		ID:    ID,
		Chr:   chr,
		Start: int(start),
		End:   int(end),
	}
}

func ImportInteractions(ct string, Filename string) {
	f, err := os.Open(Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t' // Tab delimmtted
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Now have rec, which contains the data from the line.
		tempint := Interaction{CellType: ct}
		newid, err := tempint.Create()
		tempint.ID = newid
		if err != nil {
			log.Fatal(err) // Shouldn't be any reason this fails that isn't fatal.
		}

		// Support for n-wise interactions.
		for i := 0; i < len(rec); i++ {
			if !LocusExist(rec[i]) {
				LocusFromID(rec[i]).Create() // Create the Locus if it doesn't already exist.
			}
			tempint.AddLocus(rec[i]) // Add the locus to the new interaction.
		}
	}
}

// This is currently hardcoded for the type of file we have.
func ImportGeneExpressions(Filename string) {
	f, err := os.Open(Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Column 1 is the gene name. Test for existence.
		name := rec[0]
		if !GeneExists(name) {
			start, err := strconv.ParseInt(rec[2], 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			end, err := strconv.ParseInt(rec[3], 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			Gene{Name: name, Chr: rec[0], Start: int(start), End: int(end)}.Create()
		}

		// Now we know the gene exists, create expressions for either cell type.
		dn_exp, err := strconv.ParseFloat(rec[4], 64)
		if err != nil {
			log.Fatal(err)
		}
		pgn_exp, err := strconv.ParseFloat(rec[6], 64)
		if err != nil {
			log.Fatal(err)
		}
		GeneExpression{CellType: "DN", Gene: name, ExpressionLevel: dn_exp}.Create()
		GeneExpression{CellType: "PGN", Gene: name, ExpressionLevel: pgn_exp}.Create()
	}
}

func strIsForw(test string) bool {
	return test == "+"
}

func ImportMotifInstances(ct string, Filename string) {
	f, err := os.Open(Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Column 3 is the Locus ID.
		if !LocusExist(rec[3]) {
			LocusFromID(rec[3]).Create() // Create if non-existent.
		}

		// Otherwise, just parse everything.
		start, err := strconv.ParseInt(rec[5], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		threshold, err := strconv.ParseFloat(rec[8], 64)
		if err != nil {
			log.Fatal(err)
		}
		forward := strIsForw(rec[9])

		MotifInstance{CellType: ct, Model: rec[7], Start: int(start), ThresholdScore: threshold, Chr: rec[4], Forward: forward, LocusID: rec[3]}.Create()
	}
}

func ImportMotifModels(Filename string) {
	f, err := os.Open(Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		length, err := strconv.ParseInt(rec[3], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		entrezGene, err := strconv.ParseInt(rec[6], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		MotifModel{Name: rec[0], TranscriptionFactor: rec[2], Length: int(length), Quality: rec[4][0], TFFamily: rec[5], EntrezGene: int(entrezGene), UniprotID: rec[7]}.CreateMotifModel()
	}
}

func RunDataLoader() {
	// This assumes the cell types DN and PGN have already been imported through other means.
	ImportMotifModels("./motif_models.tsv")
	ImportInteractions("DN", "./dn_interactions.tsv")
	ImportInteractions("PGN", "./pgn_interactions.tsv")
	ImportMotifInstances("DN", "./dn_motifs.bed")
	ImportMotifInstances("PGN", "./pgn_motifs.bed")
	ImportGeneExpressions("./gene_expressions.csv")
}
