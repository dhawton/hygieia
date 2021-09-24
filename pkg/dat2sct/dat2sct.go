package dat2sct

import (
	"bufio"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func openDat(filename string) string {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	star := ""
	lastline := ""
	line := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text()[0] != '!' {
			if strings.Count(line, "N") == 2 {
				star += line + "\n                          "
				line = lastline
			}
			if scanner.Text()[0] == 'L' {
				if strings.Count(line, "N") == 1 {
					line += " " + lastline
					line += "\n                          "
					star += line
					line = ""
				}
			} else {
				lineArray := strings.Fields(scanner.Text())
				lastline = "N" + lineArray[1] + "." + lineArray[2] + "." + lineArray[3] + " " +
					"W" + lineArray[4] + "." + lineArray[5] + "." + lineArray[6]
				line += " " + lastline
			}
		}
	}
	return star
}

func Convert(filename string) {
	stardata := openDat(filename)
	mapname := strings.TrimSuffix(filename, ".dat")
	if filename[0:2] == "./" || filename[0:2] == ".\\" {
		mapname = mapname[2:]
	}
	data :=
		"[INFO]\n" +
			"ZXX\n" +
			"ZXX_CTR\n" +
			"ZXX\n" +
			"N041.19.51.676\n" +
			"W080.40.28.760\n" +
			"60\n" +
			"46\n" +
			"9.2\n" +
			"1\n\n" +
			"[VOR]\n\n" +
			"[NDB]\n\n" +
			"[AIRPORT]\n" +
			"KBED 118.500 N042.28.11.831 W071.17.20.507 D\n\n" +
			"[RUNWAY]\n\n" +
			"[FIXES]\n\n" +
			"[ARTCC]\n\n" +
			"[ARTCC HIGH]\n\n" +
			"[ARTCC LOW]\n\n" +
			"[SID]\n\n" +
			"[STAR]\n" +
			mapname + "                 N099.00.00.000 E099.00.00.000 N099.00.00.000 E099.00.00.000\n                          " +
			stardata +
			"\n[LOW AIRWAY]\n\n" +
			"[HIGH AIRWAY]\n\n" +
			"[GEO]\n\n" +
			"[REGIONS]\n\n" +
			"[LABELS]\n\n"

	err := os.WriteFile("output.sct2", []byte(data), 0644)
	check(err)

}
