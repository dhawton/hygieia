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

package clean

import (
	"encoding/json"
	"os"

	"hawton.dev/hygieia/pkg/sct2parse"
	"hawton.dev/log4g"
)

var log = log4g.Category("clean")

func Start(input string, output string, cfg Config) error {
	sct2, err := sct2parse.Parse(input)
	if err != nil {
		log.Error("Error parsing sct2: %s", err.Error())
		return err
	}

	json, err := json.MarshalIndent(sct2, "", "  ")
	if err != nil {
		return err
	}

	os.WriteFile(output, json, 0644)

	return nil
}
