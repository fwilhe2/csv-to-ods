package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"os"
	"strings"

	rb "github.com/fwilhe2/rechenbrett"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type CsvOptions struct {
	HeaderLines int      `json:"headerLines"`
	Comma       string   `json:"comma"`
	Types       []string `json:"types"`
}

func main() {
	flatPtr := flag.Bool("flat", false, "produce flat ods")
	inputFilePtr := flag.String("input", "input.csv", "input csv file")
	outputFilePtr := flag.String("output", "spreadsheet.ods", "output (flat-)ods file")

	flag.Parse()

	dat, err := os.ReadFile(*inputFilePtr)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	var csvOptions CsvOptions

	csvOptionsString, err := os.ReadFile(*inputFilePtr + ".options.json")
	if err != nil {
		println("no options file")
	} else {
		check(json.Unmarshal(csvOptionsString, &csvOptions))
	}

	body := bytes.TrimPrefix(dat, []byte("\xef\xbb\xbf"))
	r := csv.NewReader(strings.NewReader(string(body)))
	r.Comma = []rune(csvOptions.Comma)[0]
	r.FieldsPerRecord = -1

	records, err := r.ReadAll()
	check(err)

	cells := csvRecordsToOdtCells(records, csvOptions)

	spreadsheet := rb.MakeSpreadsheet(cells)

	if *flatPtr {
		if strings.HasSuffix(*outputFilePtr, ".ods") {
			*outputFilePtr = strings.Replace(*outputFilePtr, ".ods", ".fods", -1)
		}
		os.WriteFile(*outputFilePtr, []byte(rb.MakeFlatOds(spreadsheet)), 0o644)
	} else {
		buff := rb.MakeOds(spreadsheet)

		archive, err := os.Create(*outputFilePtr)
		if err != nil {
			panic(err)
		}

		archive.Write(buff.Bytes())
		archive.Close()
	}
}

func csvRecordsToOdtCells(records [][]string, csvOptions CsvOptions) [][]rb.Cell {
	var cells [][]rb.Cell
	for rowIndex, rows := range records {
		var xmlRow []rb.Cell
		for columnIndex, value := range rows {
			if rowIndex < csvOptions.HeaderLines {
				xmlRow = append(xmlRow, rb.MakeCell(value, "string"))
			} else {
				if csvOptions.Types != nil {
					xmlRow = append(xmlRow, rb.MakeCell(value, csvOptions.Types[columnIndex]))
				} else {
					xmlRow = append(xmlRow, rb.MakeCell(value, "string"))
				}
			}
		}
		cells = append(cells, xmlRow)
	}
	return cells
}
