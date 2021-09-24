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

import (
	"fmt"
	"strings"

	"hawton.dev/hygieia/pkg/utils"
)

func ValidateConfig(cfg *Config) error {
	if !strings.EqualFold(strings.ToLower(cfg.Filter.Type), "radius") && !strings.EqualFold(strings.ToLower(cfg.Filter.Type), "polygon") {
		return fmt.Errorf("invalid filter type %s, expected radius or polygon", cfg.Filter.Type)
	}

	if !strings.EqualFold(strings.ToLower(cfg.Filter.Direction), "inside") && !strings.EqualFold(strings.ToLower(cfg.Filter.Direction), "outside") {
		return fmt.Errorf("invalid filter direction %s, expected inside or outside", cfg.Filter.Direction)
	}

	if !strings.EqualFold(strings.ToLower(cfg.Filter.Type), "radius") {
		if cfg.Radius.Radius <= 0 {
			return fmt.Errorf("invalid filter radius %f, expected > 0", cfg.Radius.Radius)
		}

		if !utils.StringInSlice(cfg.Radius.Unit, []string{"km", "mi", "sm", "nm"}) {
			return fmt.Errorf("invalid filter radius unit %s, expected km, mi, sm or nm", cfg.Radius.Unit)
		}
	}

	return nil
}
