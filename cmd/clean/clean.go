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

	"hawton.dev/hygieia/internal/config"
	"hawton.dev/hygieia/pkg/geo"
	"hawton.dev/hygieia/pkg/sct2parse"
	"hawton.dev/hygieia/pkg/utils"
	"hawton.dev/log4g"
)

var log = log4g.Category("clean")

func Start(input string, output string, cfg config.Config) error {
	if utils.StringEquals(cfg.Filter.Type, "polygon") {
		log.Info("Building polygon")

		points := []geo.Point{}
		for _, p := range cfg.Points {
			points = append(points, geo.Point{X: p.Lat, Y: p.Lon})
		}

		poly := geo.Polygon{Points: points}
		log.Debug("Polygon is: %q", poly)
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

func CleanSCT2(sct2 *sct2parse.Sct2, cfg config.Config) {
	for i, m := range sct2.Maps {
		for j, line := range m.Lines {
			if !shouldInclude(line, geo.Polygon{}, cfg) {
				sct2.Maps[i].Lines[j].Remove = true
			}
		}
	}
}

func shouldInclude(line sct2parse.Sct2Line, poly geo.Polygon, config config.Config) bool {
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
