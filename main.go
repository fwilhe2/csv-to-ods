package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"strings"

	rb "github.com/fwilhe2/rechenbrett"
)

var version = "dev"

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
	debugPtr := flag.Bool("debug", false, "output debug logs")
	versionPtr := flag.Bool("version", false, "print version and exit")

	flag.Parse()

	if *versionPtr {
		println("csv-to-ods", version)
		os.Exit(0)
	}

	opts := &slog.HandlerOptions{}
	if *debugPtr {
		opts.Level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	logger.Debug(*inputFilePtr)
	csvInputString, err := os.ReadFile(*inputFilePtr)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	var csvOptions CsvOptions

	csvOptionsString, err := os.ReadFile(*inputFilePtr + ".options.json")
	if err != nil {
		logger.Debug("no options file")
	} else {
		check(json.Unmarshal(csvOptionsString, &csvOptions))
	}

	records, err := parseCsv(csvInputString, csvOptions)
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

func parseCsv(csvInputString []byte, csvOptions CsvOptions) ([][]string, error) {
	body := bytes.TrimPrefix(csvInputString, []byte("\xef\xbb\xbf"))
	csvReader := csv.NewReader(strings.NewReader(string(body)))
	if len(csvOptions.Comma) > 0 {
		csvReader.Comma = []rune(csvOptions.Comma)[0]
	}
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	return records, err
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
