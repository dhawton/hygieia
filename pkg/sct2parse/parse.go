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

package sct2parse

import (
	"fmt"
	"regexp"

	"hawton.dev/hygieia/pkg/utils"
)

var (
	commentRegex  = regexp.MustCompile(`^\s*;`)
	newBlock      = regexp.MustCompile(`^\[`)
	sidRegex      = regexp.MustCompile(`(?i)^\[SID\]`)
	starRegex     = regexp.MustCompile(`(?i)^\[STAR\]`)
	fixBlockRegex = regexp.MustCompile(`(?i)^\[FIXES\]`)
	mapRegex      = regexp.MustCompile(`^(\S+)\s+`)
	lineRegex     = regexp.MustCompile(`(?i)^\s+([NS])([0-9.]+)\s+([EW])([0-9.]+)\s+([NS])([0-9.]+)\s+([EW])([0-9.]+)\s+(.+)$`)
	lineFixRegex  = regexp.MustCompile(`(?i)^\s+([A-Z0-9]{5})\s+([A-Z0-9]{5})\s+([A-Z0-9]{5})\s+([A-Z0-9]{5})\s+(.+)$`)
	fixRegex      = regexp.MustCompile(`(?i)([A-Z0-9]{5})\s+([NS])([0-9.]+)\s*([EW])([0-9.]+)`)
)

var MapOnly bool

func Parse(file string) (*Sct2, error) {
	sct2 := &Sct2{}
	sct2.Fixes = make(map[string]Sct2Point)
	lines, err := utils.ReadFileSlice(file)
	if err != nil {
		return nil, err
	}

	workingMap := Sct2Map{}
	isInSlice := false
	isSID := false
	inMaps := false
	inFixes := false

	for _, line := range lines {
		if len(line) <= 0 || commentRegex.MatchString(line) {
			sct2.RawLines = append(sct2.RawLines, line)
			continue
		}

		if !MapOnly && !inFixes && !inMaps && !sidRegex.MatchString(line) && !starRegex.MatchString(line) && !fixBlockRegex.MatchString(line) {
			sct2.RawLines = append(sct2.RawLines, line)
			continue
		}

		if fixBlockRegex.MatchString(line) {
			inMaps = false
			inFixes = true
			sct2.RawLines = append(sct2.RawLines, line)
			continue
		} else if sidRegex.MatchString(line) || starRegex.MatchString(line) {
			if isInSlice {
				sct2.Maps = append(sct2.Maps, workingMap)
				isInSlice = false
			}
			if sidRegex.MatchString(line) {
				isSID = true
			} else {
				isSID = false
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
			sct2.RawLines = append(sct2.RawLines, line)
			if fixRegex.MatchString(line) {
				res := fixRegex.FindAllStringSubmatch(line, -1)
				pt, err := ConvertSct2Point(res[0][2], res[0][3], res[0][4], res[0][5])
				if err != nil {
					return nil, err
				}
				sct2.Fixes[res[0][1]] = pt
			}
		}

		if inMaps || MapOnly {
			if mapRegex.MatchString(line) {
				if isInSlice {
					sct2.Maps = append(sct2.Maps, workingMap)
					isInSlice = false
				}

				res := mapRegex.FindAllStringSubmatch(line, -1)
				workingMap = Sct2Map{
					Name:        res[0][1],
					RawNameLine: line,
					IsSID:       isSID,
				}
				isInSlice = true
			} else if isInSlice && lineFixRegex.MatchString(line) {
				res := lineFixRegex.FindAllStringSubmatch(line, -1)
				l := Sct2Line{
					Start:      sct2.Fixes[res[0][1]],
					End:        sct2.Fixes[res[0][3]],
					LineEnding: res[0][5],
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
					Start:      point1,
					End:        point2,
					LineEnding: res[0][9],
				})

			}
		}
	}

	if isInSlice {
		sct2.Maps = append(sct2.Maps, workingMap)
		isInSlice = false
	}

	return sct2, nil
}

func (s *Sct2) ToSct2() ([]string, error) {
	lines := make([]string, 0)

	lines = append(lines, s.RawLines...)

	if !MapOnly {
		lines = append(lines, "[SID]")
		lines = append(lines, s.getMaps(true)...)
		lines = append(lines, "[STAR]")
	}
	lines = append(lines, s.getMaps(false)...)

	return lines, nil
}

func (s *Sct2) getMaps(isSid bool) []string {
	lines := make([]string, 0)
	for _, mapData := range s.Maps {
		if mapData.IsSID == isSid {
			lines = append(lines, mapData.RawNameLine)
			for _, line := range mapData.Lines {
				if !line.Remove {
					start := ConvertToSct2(line.Start.Lat, line.Start.Lon)
					end := ConvertToSct2(line.End.Lat, line.End.Lon)
					lines = append(lines, fmt.Sprintf("\t%s %s %s %s %s", start[0], start[1], end[0], end[1], line.LineEnding))
				}
			}
		}
	}

	return lines
}
