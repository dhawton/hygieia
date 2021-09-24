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

package geo

import "math"

func ConvertNMToKM(nm float64) float64 {
	return nm * 1.852
}

func ConvertSMToKM(sm float64) float64 {
	return sm * 1.15078
}

func ConvertKMToNM(km float64) float64 {
	return km / 1.852
}

func ConvertKMToSM(km float64) float64 {
	return km / 1.15078
}

func ConvertDegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func CalcGreatCircleDistance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	radius := 6371.0
	lat1 = ConvertDegreesToRadians(lat1)
	lon1 = ConvertDegreesToRadians(lon1)
	lat2 = ConvertDegreesToRadians(lat2)
	lon2 = ConvertDegreesToRadians(lon2)

	d_lat := lat2 - lat1
	d_lon := lon2 - lon1
	h := math.Sin(d_lat/2)*math.Sin(d_lat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(d_lon/2)*math.Sin(d_lon/2)

	return 2 * radius * math.Asin(math.Sqrt(h))
}
