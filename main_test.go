package main

import (
	"testing"
)

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
