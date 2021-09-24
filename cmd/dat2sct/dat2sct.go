/*
 * Hygieia - sct2 cleaner
 * Copyright (C) 2021 Daniel A. Hawton <daniel@hawton.com>, Raaj Patel
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option)
 * any later version.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of  MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for
 * more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package dat2sct

import (
	"fmt"
	"os"

	"hawton.dev/hygieia/pkg/dat2parse"
	"hawton.dev/hygieia/pkg/sct2parse"
	"hawton.dev/log4g"
)

var log = log4g.Category("dat2sct")

func Start(input string, output string, mapname string, maponly bool) error {
	sct2 := sct2parse.Sct2{}

	if !maponly {
		log.Info("Building sct2 template")
		// Prefill a basic sct2
		sct2.RawLines = []string{
			"[INFO]",
			"ZXX",
			"ZXX_CTR",
			"ZXX",
			"N041.19.51.676",
			"W080.40.28.760",
			"60",
			"46",
			"9.2",
			"1",
			"[VOR]",
			"[NDB]",
			"[AIRPORT]",
			"KBED 118.500 N042.28.11.831 W071.17.20.507 D",
			"[RUNWAY]",
			"[FIXES]",
			"[ARTCC]",
			"[ARTCC HIGH]",
			"[ARTCC LOW]",
			"[SID]",
			"[STAR]",
		}
	} else {
		sct2parse.MapOnly = true
	}

	log.Info("Defining map details")
	sct2.Maps = append(sct2.Maps, sct2parse.Sct2Map{
		Name:        mapname,
		RawNameLine: fmt.Sprintf("%s\tN000.00.00.000 W000.00.00.000 N000.00.00.000 W000.00.00.000", mapname),
	})
	log.Info("Parsing and converting dat")
	err := dat2parse.Parse(&sct2, input)
	if err != nil {
		log.Error("Error converting dat file: %s", err.Error())
		return err
	}

	log.Info("Writing sct2")
	data, err := sct2.ToSct2()
	if err != nil {
		log.Error("Error Converting sct2: %s", err.Error())
		return err
	}

	f, err := os.Create(output)
	if err != nil {
		log.Error("Error creating output file: %s", err.Error())
		return err
	}
	defer f.Close()
	for _, value := range data {
		fmt.Fprintf(f, "%s\n", value)
	}
	log.Info("Done")

	return nil
}
