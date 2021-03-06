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

package config

type Config struct {
	Filter  Filter  `yaml:"filter"`
	Points  []Point `yaml:"points"`
	Radius  Radius  `yaml:"radius"`
	MapOnly bool    `yaml:"map_only"`
}

type Filter struct {
	Type      string `yaml:"type"`
	Direction string `yaml:"direction"`
}

type Radius struct {
	Center   Point   `yaml:"center"`
	Radius   float64 `yaml:"radius"`
	Unit     string  `yaml:"unit"`
	KMRadius float64 `yaml:"-"`
}

type Distance struct {
	Radius float64 `yaml:"radius"`
	Unit   string  `yaml:"unit"`
}

type Point struct {
	Lat float64 `yaml:"lat"`
	Lon float64 `yaml:"lon"`
}
