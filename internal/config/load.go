package config

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func Load(to any, lookAt []string) error {
	var notfound int
	var content []byte
	var err error
	for i := 0; i < len(lookAt); i++ {
		content, err = os.ReadFile(lookAt[i])
		if err == nil {
			break
		}
		if os.IsNotExist(err) || os.IsPermission(err) {
			notfound++
		}
	}
	if notfound == len(lookAt) {
		return fmt.Errorf("no config found or accessible at %v", lookAt)
	}
	if len(content) == 0 {
		return fmt.Errorf("no config content")
	}

	err = toml.Unmarshal(content, to)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	return nil
}
