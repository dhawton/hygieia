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

package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"hawton.dev/hygieia/cmd/clean"
	"hawton.dev/hygieia/cmd/dat2sct"
	"hawton.dev/hygieia/internal/utils"
	"hawton.dev/log4g"
)

func main() {
	app := cli.App{
		Name:  "Hygieia",
		Usage: "Clean your sct2 maps of unneeded information",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose output",
				Aliases: []string{"v"},
			},
		},
		Commands: []*cli.Command{
			clean.Command(),
			dat2sct.Command(),
			{
				Name:  "version",
				Usage: "Get version information",
				Action: func(c *cli.Context) error {
					// Can just return since version is part of the global run.
					return nil
				},
			},
		},
		Before: func(c *cli.Context) error {
			utils.GlobalRun(c)
			return nil
		},
		After: func(c *cli.Context) error {
			if c.Bool("verbose") {
				log4g.Category("main").Debug("Setting debug")
				log4g.SetLogLevel(log4g.DEBUG)
			}
			return nil
		},
	}

	app.Run(os.Args)
}
