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

package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"hawton.dev/hygieia/internal/clean"
	"hawton.dev/hygieia/pkg/config"
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
			{
				Name:      "clean",
				Usage:     "Clean your sct2 maps of unneeded information",
				ArgsUsage: "[input file] [output file]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "Path to config file",
						Value:   "config.yaml",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() != 2 {
						return cli.Exit("Missing required arguments", 1)
					}

					verbose := c.Bool("verbose")
					if verbose {
						log4g.SetLogLevel(log4g.DEBUG)
					}

					input := c.Args().Get(0)
					output := c.Args().Get(1)
					cfg := c.String("config")

					if _, err := os.Stat(cfg); os.IsNotExist(err) {
						return cli.Exit("Config "+cfg+" file does not exist", 1)
					}

					if _, err := os.Stat(input); os.IsNotExist(err) {
						return cli.Exit("Input file does not exist", 1)
					}

					yml := clean.Config{}
					err := config.LoadConfigYaml(cfg, &yml)
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}

					return clean.Start(input, output, yml)
				},
			},
		},
	}

	app.Run(os.Args)
}
