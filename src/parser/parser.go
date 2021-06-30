package parser

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	nc "../nanocube"
)

//ReadCsvFile return records for a csv file with given filename
func ReadCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

//ParseObjects parse a csv files into objects with Lat Lon time and type
func ParseObjects(filename string, typeHead string) []nc.Object {
	LngIndex := -1 //column index for Lng
	LatIndex := -1 //column index for Lat
	TypeIndex := -1
	records := ReadCsvFile(filename)
	res := []nc.Object{}
	for i := 0; i < len(records[0]); i++ {
		if records[0][i] == "Longitude" {
			LngIndex = i
		} else if records[0][i] == "Latitude" {
			LatIndex = i
		} else if records[0][i] == typeHead {
			TypeIndex = i
		}
	}
	if LngIndex == -1 || LatIndex == -1 || TypeIndex == -1 {
		return res
	}
	for i := 1; i < len(records); i++ {
		lngstr := records[i][LngIndex]
		if lngstr == "" {
			continue
		}
		latstr := records[i][LatIndex]
		if latstr == "" {
			continue
		}
		lng, err := strconv.ParseFloat(records[i][LngIndex], 64)
		if err != nil {
			log.Fatal("Longitude is not a valid float", err)
		}
		lat, err1 := strconv.ParseFloat(records[i][LatIndex], 64)
		if err1 != nil {
			log.Fatal("Latitude is not a valid float", err)
		}
		ty := records[i][TypeIndex]
		res = append(res, nc.Object{Lng: lng, Lat: lat, Type: ty})
	}
	return res
}
