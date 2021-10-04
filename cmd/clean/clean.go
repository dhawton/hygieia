/*
 * Hygieia - sct2 cleaner
 * Copyright (C) 2021 Daniel A. Hawton <daniel@hawton.com>
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

package clean

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	internalConfig "hawton.dev/hygieia/internal/config"
	internalUtils "hawton.dev/hygieia/internal/utils"
	"hawton.dev/hygieia/pkg/config"
	"hawton.dev/hygieia/pkg/geo"
	"hawton.dev/hygieia/pkg/sct2parse"
	"hawton.dev/hygieia/pkg/utils"
	"hawton.dev/log4g"
)

var log = log4g.Category("clean")

func Command() *cli.Command {
	return &cli.Command{
		Name:      "clean",
		Usage:     "Clean your sct2 maps of unneeded information",
		ArgsUsage: "[input file] [output file]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config file",
				Value:   "config.yaml",
			},
			&cli.BoolFlag{
				Name:    "maponly",
				Usage:   "Input is not a full sct2, just a map",
				Aliases: []string{"m"},
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 2 {
				return cli.Exit("Missing required arguments", 1)
			}

			internalUtils.GlobalRun(c)

			input := c.Args().Get(0)
			output := c.Args().Get(1)
			cfg := c.String("config")

			if _, err := os.Stat(cfg); os.IsNotExist(err) {
				return cli.Exit("Config "+cfg+" file does not exist", 1)
			}

			if _, err := os.Stat(input); os.IsNotExist(err) {
				return cli.Exit("Input file does not exist", 1)
			}

			yml := internalConfig.Config{}
			err := config.LoadConfigYaml(cfg, &yml)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := internalConfig.ValidateConfig(&yml); err != nil {
				fmt.Printf("Error processing config: %s", err.Error())
				return cli.Exit("Config file is invalid", 1)
			}

			if c.Bool("maponly") {
				yml.MapOnly = true
			}

			return Start(input, output, yml)
		},
	}
}

func Start(input string, output string, cfg internalConfig.Config) error {
	if utils.StringEquals(cfg.Filter.Type, "polygon") {
		log.Info("Building polygon")

		points := []geo.Point{}
		for _, p := range cfg.Points {
			points = append(points, geo.Point{X: p.Lat, Y: p.Lon})
		}

		poly := geo.Polygon{Points: points}
		log.Debug("Polygon is: %q", poly)
	}

	if utils.StringEquals(cfg.Filter.Type, "radius") {
		if utils.StringInSlice(cfg.Radius.Unit, []string{"mi", "sm", "nm"}) {
			if utils.StringEquals(cfg.Radius.Unit, "nm") {
				cfg.Radius.KMRadius = geo.ConvertNMToKM(cfg.Radius.Radius)
			} else {
				cfg.Radius.KMRadius = geo.ConvertSMToKM(cfg.Radius.Radius)
			}
		} else {
			cfg.Radius.KMRadius = cfg.Radius.Radius
		}
		log.Debug("KMRadius set to %f", cfg.Radius.KMRadius)
	}

	if cfg.MapOnly {
		log.Info("Map Only mode enabled")
		sct2parse.MapOnly = true
	}

	log.Info("Parsing sct2 file")
	sct2, err := sct2parse.Parse(input)
	if err != nil {
		log.Error("Error parsing sct2: %s", err.Error())
		return err
	}

	dat, _ := json.MarshalIndent(sct2, "", "  ")
	log.Debug("Sct2: %s", string(dat))

	log.Info("Checking for lines to filter")
	CleanSCT2(sct2, cfg)

	log.Info("Converting back to sct2")
	lines, err := sct2.ToSct2()
	if err != nil {
		log.Error("Error converting to sct2: %s", err.Error())
		return err
	}

	f, err := os.Create(output)
	if err != nil {
		log.Error("Error creating output file: %s", err.Error())
		return err
	}
	defer f.Close()
	for _, value := range lines {
		fmt.Fprintf(f, "%s\n", value)
	}
	log.Info("Done")

	return nil
}

func CleanSCT2(sct2 *sct2parse.Sct2, cfg internalConfig.Config) {
	for i, m := range sct2.Maps {
		for j, line := range m.Lines {
			if !shouldInclude(line, geo.Polygon{}, cfg) {
				sct2.Maps[i].Lines[j].Remove = true
			}
		}
	}
}

func shouldInclude(line sct2parse.Sct2Line, poly geo.Polygon, config internalConfig.Config) bool {
	var containsStart bool
	var containsEnd bool
	filter := config.Filter

	if strings.EqualFold(strings.ToLower(filter.Type), "polygon") {
		containsStart = geo.PointInPolygon(geo.Point{X: line.Start.Lat, Y: line.Start.Lon}, poly)
		containsEnd = geo.PointInPolygon(geo.Point{X: line.End.Lat, Y: line.End.Lon}, poly)
	} else if strings.EqualFold(strings.ToLower(filter.Type), "radius") {
		containsStart = geo.CalcGreatCircleDistance(line.Start.Lat, line.Start.Lon, config.Radius.Center.Lat, config.Radius.Center.Lon) <= config.Radius.KMRadius
		containsEnd = geo.CalcGreatCircleDistance(line.End.Lat, line.End.Lon, config.Radius.Center.Lat, config.Radius.Center.Lon) <= config.Radius.KMRadius
	}

	if strings.EqualFold(strings.ToLower(filter.Direction), "inside") {
		if containsStart || containsEnd {
			log.Debug("Filtering line: %v", line)
			return false
		}
	}

	if strings.EqualFold(strings.ToLower(filter.Direction), "outside") {
		if !containsStart || !containsEnd {
			log.Debug("Filtering line: %v", line)
			return false
		}
	}

	return true
}
