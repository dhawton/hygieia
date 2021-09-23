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

type Sct2 struct {
	RawLines []string
	Maps     []Sct2Map
	Fixes    map[string]Sct2Point
}

type Sct2Map struct {
	Name        string
	RawNameLine string
	Lines       []Sct2Line
	IsSID       bool
}

type Sct2Line struct {
	Start  Sct2Point
	End    Sct2Point
	Remove bool
}

type Sct2Point struct {
	Lat float64
	Lon float64
}
