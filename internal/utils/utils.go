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

package utils

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/urfave/cli/v2"
	"hawton.dev/hygieia/internal/version"
	"hawton.dev/log4g"
)

func GlobalRun(c *cli.Context) {
	intro := figure.NewFigure("Hygieia", "", false).Slicify()
	for i := 0; i < len(intro); i++ {
		fmt.Println(intro[i])
	}
	fmt.Println("Hygieia " + version.FriendlyVersion())
	fmt.Println("")
	fmt.Println("Hygieia Copyright (C) 2021 Daniel A. Hawton <daniel@hawton.com>, Raaj Patel")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY.")
	fmt.Println("This is free software, and you are welcome to redistribute it")
	fmt.Println("under certain conditions; view license at https://www.gnu.org/licenses/gpl-3.0.en.html.")
	fmt.Println("")

	verbose := c.Bool("verbose")
	if verbose {
		log4g.SetLogLevel(log4g.DEBUG)
	}
}
