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

package dat2parse

import (
	"fmt"
	"regexp"

	"hawton.dev/hygieia/pkg/sct2parse"
	"hawton.dev/hygieia/pkg/utils"
)

var (
	commentRegex = regexp.MustCompile(`^!`)
	newLine      = regexp.MustCompile(`^LINE !`)
	point        = regexp.MustCompile(`GP (\d+) (\d+) ([0-9.]+)\s+(\d+) (\d+) ([0-9.]+)\s+!`)
)

func Parse(sct2 *sct2parse.Sct2, file string) error {
	lines, err := utils.ReadFileSlice(file)
	if err != nil {
		return err
	}

	var lastPoint sct2parse.Sct2Point

	for _, line := range lines {
		if commentRegex.MatchString(line) {
			continue
		}

		if newLine.MatchString(line) {
			lastPoint = sct2parse.Sct2Point{}
			continue
		}

		if point.MatchString(line) {
			match := point.FindAllStringSubmatch(line, -1)
			point, err := sct2parse.ConvertSct2Point("N", fmt.Sprintf("%s.%s.%s", match[0][1], match[0][2], match[0][3]), "W", fmt.Sprintf("%s.%s.%s", match[0][4], match[0][5], match[0][6]))
			if err != nil {
				return err
			}
			if (sct2parse.Sct2Point{}) != lastPoint {
				sct2line := sct2parse.Sct2Line{
					Start: lastPoint,
					End:   point,
				}
				sct2.Maps[0].Lines = append(sct2.Maps[0].Lines, sct2line)
			}
			lastPoint = point
		}
	}

	return nil
}
