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
	"math"
	"strconv"
	"strings"
)

func convertFloat(num string) (float64, error) {
	var f float64
	var e error
	if strings.Contains(num, ".") {
		f, e = strconv.ParseFloat(num, 64)
	} else {
		f, e = strconv.ParseFloat(num+".0", 64)
	}
	return f, e
}

func convertDMSToDD(dms string) (float64, error) {
	var d float64
	num := strings.Split(dms, ".")
	deg, err := convertFloat(num[0])
	if err != nil {
		return 0, err
	}

	d = deg

	min, err := convertFloat(num[1])
	if err != nil {
		return 0, err
	}
	d = d + min/60

	sec, err := convertFloat(fmt.Sprintf("%s.%s", num[2], num[3]))
	if err != nil {
		return 0, err
	}
	d = d + sec/3600

	return d, nil
}

func ConvertSct2Point(latType string, latNumber string, lonType string, lonNumber string) (Sct2Point, error) {
	var lat float64
	var lon float64

	lat, err := convertDMSToDD(latNumber)
	if err != nil {
		return Sct2Point{}, err
	}

	if latType == "S" {
		lat = -lat
	}

	lon, err = convertDMSToDD(lonNumber)
	if err != nil {
		return Sct2Point{}, err
	}

	if lonType == "W" {
		lon = -lon
	}

	return Sct2Point{lat, lon}, nil
}

func convertDDToDMS(dd float64) string {
	deg := int(math.Floor(dd))
	min := int(math.Floor((dd - float64(deg)) * 60))
	sec := ((dd-float64(deg))*60 - float64(min)) * 60

	return fmt.Sprintf("%03d.%d.%0.3f", deg, min, sec)
}

func ConvertToSct2(lat float64, lon float64) []string {
	latNegative := lat < 0
	lonNegative := lon < 0
	latAbs := math.Abs(lat)
	lonAbs := math.Abs(lon)
	latString := convertDDToDMS(latAbs)
	lonString := convertDDToDMS(lonAbs)

	if latNegative {
		latString = "S" + latString
	} else {
		latString = "N" + latString
	}

	if lonNegative {
		lonString = "W" + lonString
	} else {
		lonString = "E" + lonString
	}

	return []string{latString, lonString}
}
