/*
Copyright 2018 Comcast Cable Communications Management, LLC
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
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func client(c *cli.Context) *vinyldns.Client {
	validateEnv(c)
	return &vinyldns.Client{
		AccessKey:  c.GlobalString(accessKeyFlag),
		SecretKey:  c.GlobalString(secretKeyFlag),
		Host:       c.GlobalString(hostFlag),
		HTTPClient: &http.Client{},
	}
}

func getOption(c *cli.Context, name string) (string, error) {
	val := c.String(name)
	var err error
	if len(val) == 0 {
		err = fmt.Errorf("--%s is required", name)
	}
	return val, err
}

func validateEnv(c *cli.Context) {
	h := c.GlobalString(hostFlag)
	ak := c.GlobalString(accessKeyFlag)
	sk := c.GlobalString(secretKeyFlag)
	missing := []string{}

	if h == "" {
		missing = append(missing, h)
		fmt.Printf("\nPlease pass '--%s' or set 'VINYLDNS_HOST'\n", hostFlag)
	}
	if ak == "" {
		missing = append(missing, h)
		fmt.Printf("\nPlease pass '--%s' or set 'VINYLDNS_ACCESS_KEY'\n", accessKeyFlag)
	}
	if sk == "" {
		missing = append(missing, h)
		fmt.Printf("\nPlease pass '--%s' or set 'VINYLDNS_SECRET_KEY'\n", secretKeyFlag)
	}

	if len(missing) > 0 {
		os.Exit(1)
	}
}

func printBasicTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(data)
	table.SetRowLine(true)
	table.Render()
}

func printTableWithHeaders(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(data)
	table.SetRowLine(true)
	table.Render()
}

func printJSON(i interface{}) error {
	j, err := json.Marshal(i)
	if err != nil {
		return err
	}

	fmt.Println(string(j))

	return nil
}
