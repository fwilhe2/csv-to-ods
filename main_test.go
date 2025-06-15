package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestParseCsv(t *testing.T) {
	var expected [][]string
	expected = append(expected, []string{"string", "int", "float", "currency"})
	expected = append(expected, []string{"foo", "42", "23.12", "4.33"})
	expected = append(expected, []string{"bar", "23", "2", "3"})

	runTestParseCsv(t, "sample", expected)
}

func TestParseCsvSemicolon(t *testing.T) {
	var expected [][]string
	expected = append(expected, []string{"string", "int", "float", "currency"})
	expected = append(expected, []string{"foo", "42", "23.12", "4.33"})
	expected = append(expected, []string{"bar", "23", "2", "3"})

	runTestParseCsv(t, "sample-semicolon", expected)
}

func TestParseCsvBankTransactions(t *testing.T) {
	var expected [][]string
	expected = append(expected, []string{"Date", "Time", "Merchant", "Description", "Amount", "Currency"})
	expected = append(expected, []string{"2025-06-01", "10:30", "Bäcker Schmidt", "Einkauf von Brot und Brötchen", "12.50", "EUR"})
	expected = append(expected, []string{"2025-06-01", "14:15", "Kiosk Müller", "Zeitschrift, Süßigkeiten und Getränke", "7.80", "EUR"})
	expected = append(expected, []string{"", "", "", "", "20,30", ""})

	runTestParseCsv(t, "bank-transactions", expected)
}

func runTestParseCsv(t *testing.T, testCase string, expected [][]string) {
	csvInputString, err := os.ReadFile(fmt.Sprintf("test-files/%s.csv", testCase))
	if err != nil {
		t.Fail()
	}

	csvOptionsString, err := os.ReadFile(fmt.Sprintf("test-files/%s.csv.options.json", testCase))
	var csvOptions CsvOptions
	if err != nil {
		t.Fail()
	}

	check(json.Unmarshal(csvOptionsString, &csvOptions))

	actual, err := parseCsv(csvInputString, csvOptions)
	if err != nil {
		t.Fail()
	}

	for rowIndex, rowValue := range actual {
		for columnIndex, cellValue := range rowValue {
			if cellValue != expected[rowIndex][columnIndex] {
				t.Errorf("Expected %s to equal %s", cellValue, expected[rowIndex][columnIndex])
			}
		}
	}
}

func TestBuildCells(t *testing.T) {
	var records [][]string
	records = append(records, []string{"string", "int", "float", "currency"})
	records = append(records, []string{"foo", "42", "23.12", "4.33"})

	var csvOptions CsvOptions
	csvOptions.HeaderLines = 1
	csvOptions.Types = []string{"string", "int", "float", "currency"}

	actual := csvRecordsToOdtCells(records, csvOptions)

	if actual[1][0].Text != "foo" {
		t.Fail()
	}

	if actual[1][2].Value != "23.12" {
		t.Fail()
	}
}
