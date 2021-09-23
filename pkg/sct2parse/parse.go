/*
 * Hygieia - sct2 cleaner
 * Copyright (C) <year> <name of author
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

package sct2parse

import (
	"regexp"

	"hawton.dev/hygieia/pkg/utils"
)

var (
	commentRegex  = regexp.MustCompile(`(?i)^\s*;`)
	newBlock      = regexp.MustCompile(`(?i)^\[`)
	sidRegex      = regexp.MustCompile(`(?i)^\[SID\]`)
	starRegex     = regexp.MustCompile(`(?i)^\[STAR\]`)
	fixBlockRegex = regexp.MustCompile(`(?i)^\[FIXES\]`)
	mapRegex      = regexp.MustCompile(`(?i)^(\S+)\s+`)
	lineRegex     = regexp.MustCompile(`(?i)^\s+([NS])([0-9.]+)\s+([EW])([0-9.]+)\s+([NS])([0-9.]+)\s+([EW])([0-9.]+)`)
	lineFixRegex  = regexp.MustCompile(`(?i)^\s+([A-Z0-9]{5})\s+([A-Z0-9]{5})\s+([A-Z0-9]{5})\s+([A-Z0-9]{5})`)
	fixRegex      = regexp.MustCompile(`(?i)([A-Z0-9]{5})\s+([NS])([0-9.]+)\s*([EW])([0-9.]+)`)
)

func Parse(file string) (*Sct2, error) {
	sct2 := &Sct2{}
	sct2.Fixes = make(map[string]Sct2Point)
	lines, err := utils.ReadFileSlice(file)
	if err != nil {
		return nil, err
	}

	workingMap := Sct2Map{}
	isInSlice := false
	inMaps := false
	inFixes := false

	for _, line := range lines {
		if len(line) <= 0 || commentRegex.MatchString(line) {
			sct2.RawLines = append(sct2.RawLines, line)
			continue
		}

		if !inFixes && !inMaps && !sidRegex.MatchString(line) && !starRegex.MatchString(line) && !fixBlockRegex.MatchString(line) {
			sct2.RawLines = append(sct2.RawLines, line)
			continue
		}

		if fixBlockRegex.MatchString(line) {
			inMaps = false
			inFixes = true
			continue
		} else if sidRegex.MatchString(line) || starRegex.MatchString(line) {
			if isInSlice {
				sct2.Maps = append(sct2.Maps, workingMap)
				isInSlice = false
			}
			inMaps = true
			inFixes = false
			continue
		} else if inMaps && newBlock.MatchString(line) {
			if isInSlice {
				sct2.Maps = append(sct2.Maps, workingMap)
				isInSlice = false
			}
			sct2.RawLines = append(sct2.RawLines, line)
			inMaps = false
			inFixes = false
			continue
		}

		if inFixes {
			if fixRegex.MatchString(line) {
				res := fixRegex.FindAllStringSubmatch(line, -1)
				pt, err := ConvertSct2Point(res[0][2], res[0][3], res[0][4], res[0][5])
				if err != nil {
					return nil, err
				}
				sct2.Fixes[res[0][1]] = pt
			}
		}

		if inMaps {
			if mapRegex.MatchString(line) {
				if isInSlice {
					sct2.Maps = append(sct2.Maps, workingMap)
					isInSlice = false
				}

				res := mapRegex.FindAllStringSubmatch(line, -1)
				workingMap = Sct2Map{
					Name:        res[0][1],
					RawNameLine: line,
				}
				isInSlice = true
			} else if isInSlice && lineFixRegex.MatchString(line) {
				res := lineFixRegex.FindAllStringSubmatch(line, -1)
				l := Sct2Line{
					Start: sct2.Fixes[res[0][1]],
					End:   sct2.Fixes[res[0][3]],
				}
				workingMap.Lines = append(workingMap.Lines, l)
			} else if isInSlice && lineRegex.MatchString(line) {
				res := lineRegex.FindAllStringSubmatch(line, -1)

				point1, err := ConvertSct2Point(res[0][1], res[0][2], res[0][3], res[0][4])
				if err != nil {
					return nil, err
				}
				point2, err := ConvertSct2Point(res[0][5], res[0][6], res[0][7], res[0][8])
				if err != nil {
					return nil, err
				}

				workingMap.Lines = append(workingMap.Lines, Sct2Line{
					Start: point1,
					End:   point2,
				})

			}
		}
	}

	return sct2, nil
}
