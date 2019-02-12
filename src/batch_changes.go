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

	clitable "github.com/crackcomm/go-clitable"
	"github.com/urfave/cli"
)

func batchChanges(c *cli.Context) error {
	client := client(c)
	rc, err := client.BatchRecordChanges()
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(rc)
	}

	changes := []map[string]interface{}{}
	for _, r := range rc {
		m := map[string]interface{}{}
		m["ID"] = r.ID
		m["CreatedTimestamp"] = r.CreatedTimestamp
		m["Comments"] = r.Comments

		changes = append(changes, m)
	}

	if len(changes) != 0 {
		clitable.PrintTable([]string{"ID", "CreatedTimestamp", "Comments"}, changes)
	} else {
		fmt.Println("No batch changes found")
	}

	return nil
}

func batchChange(c *cli.Context) error {
	client := client(c)
	rc, err := client.BatchRecordChange(c.String("batch-change-id"))
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(rc)
	}

	change := [][]string{}
	for _, r := range rc.Changes {
		changeElem := []string{}

		changeElem = append(changeElem, `"ChangeType" - `+r.ChangeType)
		changeElem = append(changeElem, `"InputName" - `+r.InputName)
		changeElem = append(changeElem, `"Type" - `+r.Type)
		changeElem = append(changeElem, `"TTL" - `+string(r.TTL))

		recordData, err := json.Marshal(&(r.Record))
		if err != nil {
			return err
		}
		changeElem = append(changeElem, `"Record" - `+string(recordData))
		changeElem = append(changeElem, `"Status" - `+r.Status)

		change = append(change, changeElem)
	}

	if len(change) != 0 {
		printBasicTable(change)
	} else {
		fmt.Println("No batch change found with id: " + c.String("batch-change-id"))
	}

	return nil
}
