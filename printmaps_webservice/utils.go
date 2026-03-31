package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	pip "github.com/JamesMilnerUK/pip-go"
	"github.com/printmaps/printmaps/pd"
)

/*
createMapOrder is a helper to create a (asynchronous) build order for the map defined in the meta data.
*/
func createMapOrder(pmData pd.PrintmapsData) error {
	file := filepath.Join(pd.PathWorkdir, pd.PathOrders, pmData.Data.ID) + ".json"

	// create directory if necessary
	if _, err := os.Stat(pd.PathOrders); os.IsNotExist(err) {
		if err := os.MkdirAll(pd.PathOrders, 0750); err != nil {
			log.Printf("error <%v> at os.MkdirAll(), path = <%s>", err, pd.PathOrders)
			return err
		}
	}

	data, err := json.MarshalIndent(pmData, pd.IndentPrefix, pd.IndexString)
	if err != nil {
		log.Printf("error <%v> at json.MarshalIndent()", err)
		return err
	}

	return os.WriteFile(file, data, 0600)
}

/*
readPolyfile reads the polygon file (osmosis poly(gon) format).
Spec: http://wiki.openstreetmap.org/wiki/Osmosis/Polygon_Filter_File_Format
0       none
1       1
2          6.388768E+00   5.187233E+01
3          6.389918E+00   5.187448E+01
...
n-3        6.385880E+00   5.186474E+01
n-2        6.388768E+00   5.187233E+01
n-1     END
n       END
Polyfiles with more than one individual polygon are not supported.
*/
func readPolyfile(filename string, pPolygon *pip.Polygon) error {
	var lon float64
	var lat float64
	var pP pip.Point

	lines, err := slurpFile(filename)
	if err != nil {
		log.Printf("error <%v> at slurpTextfile(), file = <%v>", err, filename)
		return err
	}

	for index, line := range lines {
		if index > 1 {
			if strings.TrimSpace(line) == "END" {
				break
			}
			_, err := fmt.Sscanf(line, "%f%f", &lon, &lat)
			if err != nil {
				log.Printf("error <%v> at fmt.Sscanf(), file = <%s>, line = <%s>", err, filename, line)
				return err
			}
			pP.X = lon
			pP.Y = lat
			pPolygon.Points = append(pPolygon.Points, pP)
		}
	}

	return nil
}

/*
slurpFile slurps all lines of a text file into a slice of strings.
*/
func slurpFile(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return lines, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines, nil
}

/*
readCapafile reads the capabilities file (json format).
*/
func readCapafile(filename string, pmFeature *PrintmapsFeature) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("error <%v> at os.ReadFile(), file = <%s>", err, filename)
		return err
	}

	err = json.Unmarshal(data, pmFeature)
	if err != nil {
		log.Printf("error <%v> at json.Unmarshal()", err)
		return err
	}

	return nil
}
