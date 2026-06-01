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

func groupChange(c *cli.Context) error {
	client := client(c)
	change, err := client.GroupChange(c.String("group-change-id"))
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(change)
	}

	data := [][]string{
		{"ID", change.ID},
		{"UserID", change.UserID},
		{"UserName", change.UserName},
		{"ChangeType", change.ChangeType},
		{"Created", change.Created},
		{"Message", change.GroupChangeMessage},
	}
	printBasicTable(data)

	fmt.Println("\nNew Group...")
	printGroup(change.NewGroup)
	fmt.Println("\nOld Group...")
	printGroup(change.OldGroup)

	return nil
}

func groupValidDomains(c *cli.Context) error {
	client := client(c)
	domains, err := client.GroupValidDomains()
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(domains)
	}

	data := [][]string{}
	for _, domain := range domains {
		data = append(data, []string{domain})
	}

	if len(data) != 0 {
		printTableWithHeaders([]string{"Domain"}, data)
	} else {
		fmt.Print("No domains found")
	}

	return nil
}
