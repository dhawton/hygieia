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
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"hawton.dev/hygieia/pkg/geo"
	"hawton.dev/hygieia/pkg/sct2parse"
	"hawton.dev/log4g"
)

var log = log4g.Category("clean")

func Start(input string, output string, cfg Config) error {
	intro := figure.NewFigure("Hygieia", "", false).Slicify()
	for i := 0; i < len(intro); i++ {
		log.Info(intro[i])
	}

	log.Info("Thanks for using Hygieia")
	log.Info("")
	log.Info("Copyright (C) 2021 Daniel A. Hawton <daniel@hawton.com>")
	log.Info("")
	log.Info("This program is free software: you can redistribute it and/or modify it")
	log.Info("under the terms of the GNU General Public License as published by the Free")
	log.Info("Software Foundation, either version 3 of the License, or (at your option)")
	log.Info("any later version.")
	log.Info("")
	log.Info("This program is distributed in the hope that it will be useful, but WITHOUT")
	log.Info("ANY WARRANTY; without even the implied warranty of  MERCHANTABILITY or")
	log.Info("FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for")
	log.Info("more details.")
	log.Info("")
	log.Info("You should have received a copy of the GNU General Public License along with")
	log.Info("this program.  If not, see <http://www.gnu.org/licenses/>.")
	log.Info("")

	log.Info("Building polygon")

	points := []geo.Point{}
	for _, p := range cfg.Points {
		points = append(points, geo.Point{X: p.Lat, Y: p.Lon})
	}

	poly := geo.Polygon{Points: points}
	log.Debug("Polygon is: %q", poly)

	log.Info("Parsing sct2 file")
	sct2, err := sct2parse.Parse(input)
	if err != nil {
		log.Error("Error parsing sct2: %s", err.Error())
		return err
	}

	log.Info("Checking for lines to filter")
	for i, m := range sct2.Maps {
		for j, line := range m.Lines {
			if !shouldInclude(line, poly, cfg.Filter) {
				sct2.Maps[i].Lines[j].Remove = true
			}
		}
	}

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

func shouldInclude(line sct2parse.Sct2Line, poly geo.Polygon, filter string) bool {
	//log.Debug("Checking line %f, %f", line.Start.Lat, line.Start.Lon)
	containsStart := geo.PointInPolygon(geo.Point{X: line.Start.Lat, Y: line.Start.Lon}, poly)
	containsEnd := geo.PointInPolygon(geo.Point{X: line.End.Lat, Y: line.End.Lon}, poly)

	if strings.EqualFold(strings.ToLower(filter), "inside") {
		if containsStart || containsEnd {
			log.Debug("Filtering line: %v", line)
			return false
		}
	}

	if strings.EqualFold(strings.ToLower(filter), "outside") {
		if !containsStart || !containsEnd {
			log.Debug("Filtering line: %v", line)
			return false
		}
	}

	return true
}
