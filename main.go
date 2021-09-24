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
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/urfave/cli/v2"
	"hawton.dev/hygieia/cmd/clean"
	"hawton.dev/hygieia/cmd/dat2sct"
	internalConfig "hawton.dev/hygieia/internal/config"
	"hawton.dev/hygieia/pkg/config"
	"hawton.dev/log4g"
)

var log = log4g.Category("main")

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
					&cli.BoolFlag{
						Name:    "maponly",
						Usage:   "Input is not a full sct2, just a map",
						Aliases: []string{"m"},
						Value:   false,
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() != 2 {
						return cli.Exit("Missing required arguments", 1)
					}

					globalrun(c)

					input := c.Args().Get(0)
					output := c.Args().Get(1)
					cfg := c.String("config")

					if _, err := os.Stat(cfg); os.IsNotExist(err) {
						return cli.Exit("Config "+cfg+" file does not exist", 1)
					}

					if _, err := os.Stat(input); os.IsNotExist(err) {
						return cli.Exit("Input file does not exist", 1)
					}

					yml := internalConfig.Config{}
					err := config.LoadConfigYaml(cfg, &yml)
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}

					if err := internalConfig.ValidateConfig(&yml); err != nil {
						fmt.Printf("Error processing config: %s", err.Error())
						return cli.Exit("Config file is invalid", 1)
					}

					if c.Bool("maponly") {
						yml.MapOnly = true
					}

					return clean.Start(input, output, yml)
				},
			},
			{
				Name:      "dat2sct",
				Usage:     "Convert FAA .dat files to sct2",
				ArgsUsage: "[input file] [output file]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Name of map",
						Value:   "HYGIEIA_CONVERTED",
					},
					&cli.BoolFlag{
						Name:    "maponly",
						Usage:   "Output as just a map",
						Aliases: []string{"m"},
						Value:   false,
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() != 2 {
						return cli.Exit("Missing required arguments", 1)
					}

					globalrun(c)

					input := c.Args().Get(0)
					output := c.Args().Get(1)

					if _, err := os.Stat(input); os.IsNotExist(err) {
						return cli.Exit("Input file does not exist", 1)
					}

					return dat2sct.Start(input, output, c.String("name"), c.Bool("maponly"))
				},
			},
		},
	}

	app.Run(os.Args)
}

func globalrun(c *cli.Context) {
	intro := figure.NewFigure("Hygieia", "", false).Slicify()
	for i := 0; i < len(intro); i++ {
		log.Info(intro[i])
	}
	log.Info("Thanks for using Hygieia")
	log.Info("")
	log.Info("Hygieia Copyright (C) 2021 Daniel A. Hawton <daniel@hawton.com>, Raaj Patel")
	log.Info("This program comes with ABSOLUTELY NO WARRANTY.")
	log.Info("This is free software, and you are welcome to redistribute it")
	log.Info("under certain conditions; view license at https://www.gnu.org/licenses/gpl-3.0.en.html.")
	log.Info("")

	verbose := c.Bool("verbose")
	if verbose {
		log4g.SetLogLevel(log4g.DEBUG)
	}
}
