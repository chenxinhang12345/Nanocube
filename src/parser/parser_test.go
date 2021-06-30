package parser

import (
	"fmt"
	"testing"
)

func TestReadCSV(t *testing.T) {
	fmt.Println(ReadCsvFile("test.csv"))
}

func TestParseObjects(t *testing.T) {
	fmt.Println(ParseObjects("test.csv", "Primary Type"))
}

func TestParseBig(t *testing.T) {
	fmt.Println(ParseObjects("crime2019.csv", "Primary Type"))
}
