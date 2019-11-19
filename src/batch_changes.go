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
	"strings"
	"strconv"

	clitable "github.com/crackcomm/go-clitable"

	"github.com/urfave/cli"

	"github.com/vinyldns/go-vinyldns/vinyldns"
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
		changeElem = append(changeElem, `"TTL" - `+strconv.Itoa(r.TTL))

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

func changeList(chs []vinyldns.RecordChange) string {
	changes := []string{}

  for _, r := range chs {
    recordData, _ := json.Marshal(&(r.Record))

    changes = append(changes,
      `"ChangeType" - `+r.ChangeType,
      `"InputName" - `+r.InputName,
      `"Type" - `+r.Type,
      `"TTL" - `+strconv.Itoa(r.TTL),
      `"Record" - `+string(recordData),
      `"Status" - `+r.Status)
  }

  return strings.Join(changes, "\n")
}


func batchChangeCreate(c *cli.Context) error {
	data := []byte(c.String("json"))
	batchChange := &vinyldns.BatchRecordChange{}
	if err := json.Unmarshal(data, &batchChange); err != nil {
		return err
	}
	client := client(c)
	bc, err := client.BatchRecordChangeCreate(batchChange)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(bc)
	}

  formattedData := [][]string{
    {"ID", bc.ID},
    {"UserName", bc.UserName},
    {"UserID", bc.UserID},
    {"Status", bc.Status},
    {"Comments", bc.Comments},
    {"Changes", changeList(bc.Changes)},
    {"CreatedTimestamp", bc.CreatedTimestamp},
    {"OwnerGroupID", bc.OwnerGroupID},
    {"ApprovalStatus", bc.ApprovalStatus},
    {"ReviewerID", bc.ReviewerID},
    {"ReviewerUserName", bc.ReviewerUserName},
    {"ReviewerTimestamp", bc.ReviewTimestamp},
    {"ReviewComment", bc.ReviewComment},
    {"ScheduledTime", bc.ScheduledTime},
    {"CancelledTimestamp", bc.CancelledTimestamp},
  }

  printBasicTable(formattedData)

	return nil
}
