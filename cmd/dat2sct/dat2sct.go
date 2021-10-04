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
	"strings"

	"github.com/urfave/cli/v2"
	"hawton.dev/hygieia/internal/utils"
	"hawton.dev/hygieia/pkg/dat2parse"
	"hawton.dev/hygieia/pkg/sct2parse"
	"hawton.dev/log4g"
)

var log = log4g.Category("dat2sct")

func Command() *cli.Command {
	return &cli.Command{
		Name:      "dat2sct",
		Usage:     "Convert FAA .dat files to sct2",
		ArgsUsage: "[input file] [output file]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Name of map",
				Value:   "HYGIEIA_CONVERTED",
			},
			&cli.BoolFlag{
				Name:    "maponly",
				Usage:   "Output as just a map",
				Aliases: []string{"m"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "recursive",
				Usage:   "Recursively convert all .dat files in the directory",
				Value:   false,
				Aliases: []string{"r"},
			},
			&cli.BoolFlag{
				Name:    "recursive_split",
				Usage:   "Split the output into multiple files, one per map, in output directory (only with recursive)",
				Value:   false,
				Aliases: []string{"rs"},
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 2 {
				return cli.Exit("Missing required arguments", 1)
			}

			utils.GlobalRun(c)

			input := c.Args().Get(0)
			output := c.Args().Get(1)

			if _, err := os.Stat(input); os.IsNotExist(err) {
				return cli.Exit("Input file/directory does not exist", 1)
			}

			if c.Bool("recursive") && c.Bool("recursive_split") {
				if _, err := os.Stat(output); os.IsNotExist(err) {
					log.Info("Creating output directory")
					os.Mkdir(output, 0755)
				}
			}

			if c.Bool("recursive") && !strings.EqualFold(c.String("name"), "HYGIEIA_CONVERTED") {
				log.Warning("Ignoring name argument. In recursive mode, map name will be the dat filename.")
			}

			return Start(input, output, c.String("name"), c.Bool("maponly"))
		},
	}
}

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
		RawNameLine: fmt.Sprintf("%-26s\tN000.00.00.000 W000.00.00.000 N000.00.00.000 W000.00.00.000", mapname),
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
