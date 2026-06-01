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

	"github.com/urfave/cli"
)

func ping(c *cli.Context) error {
	client := client(c)
	resp, err := client.Ping()
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(resp)
	}

	fmt.Print(resp)
	return nil
}

func health(c *cli.Context) error {
	client := client(c)
	if err := client.Health(); err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON("OK")
	}

	fmt.Print("OK")
	return nil
}

func color(c *cli.Context) error {
	client := client(c)
	resp, err := client.Color()
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(resp)
	}

	fmt.Print(resp)
	return nil
}

func metricsPrometheus(c *cli.Context) error {
	client := client(c)
	resp, err := client.MetricsPrometheus(c.StringSlice("name"))
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(resp)
	}

	fmt.Print(resp)
	return nil
}
