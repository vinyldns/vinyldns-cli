/*
Copyright 2026 Comcast Cable Communications Management, LLC
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli"
)

func status(c *cli.Context) error {
	client := client(c)
	s, err := client.Status()
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(s)
	}

	data := [][]string{
		{"ProcessingDisabled", fmt.Sprintf("%t", s.ProcessingDisabled)},
		{"Color", s.Color},
		{"KeyName", s.KeyName},
		{"Version", s.Version},
	}

	printBasicTable(data)
	return nil
}

func statusUpdate(c *cli.Context) error {
	processingDisabledRaw, err := getOption(c, "processing-disabled")
	if err != nil {
		return err
	}

	processingDisabled, err := strconv.ParseBool(processingDisabledRaw)
	if err != nil {
		return err
	}

	client := client(c)
	s, err := client.StatusUpdate(processingDisabled)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(s)
	}

	data := [][]string{
		{"ProcessingDisabled", fmt.Sprintf("%t", s.ProcessingDisabled)},
		{"Color", s.Color},
		{"KeyName", s.KeyName},
		{"Version", s.Version},
	}

	printBasicTable(data)
	return nil
}
