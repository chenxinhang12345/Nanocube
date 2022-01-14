package parser

import (
	"fmt"
	"testing"
	"runtime"
	nc "../nanocube"
)
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
return b / 1024 / 1024
}


func TestReadCSV(t *testing.T) {
	fmt.Println(ReadCsvFile("test.csv"))
}

func TestParseObjects(t *testing.T) {
	fmt.Println(ParseObjects("test.csv", "Primary Type"))
}

func TestParseBig(t *testing.T) {
	fmt.Println(ParseObjects("crime2019.csv", "Primary Type"))
}

func TestNanoCubeFromSmallFile(t *testing.T) {
	fmt.Println(CreateNanoCubeFromCsvFile("test.csv", "Primary Type", 16))
}

func TestNanoCubeFromBigFile(t *testing.T) {
	n := CreateNanoCubeFromCsvFile("crime2020.csv", "Primary Type", 16)
	PrintMemUsage()
	fmt.Println(n.Root.Children[0].CatRoot.Summary)
	fmt.Println(n.Root.Children[1].CatRoot.Summary)
	fmt.Println(n.Root.Children[2].CatRoot.Summary)
	fmt.Println(n.Root.Children[3].CatRoot.Summary)
	fmt.Println(nc.Query(n.Root, nc.Bounds{-87.65,41.8,0.3,0.3},5))
}


