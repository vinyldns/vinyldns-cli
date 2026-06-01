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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/crackcomm/go-clitable"
	"github.com/urfave/cli"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func recordSetUpdate(c *cli.Context) error {
	data := []byte(c.String("json"))
	rs := &vinyldns.RecordSet{}
	if err := json.Unmarshal(data, &rs); err != nil {
		return err
	}
	if rs.ZoneID == "" || rs.ID == "" {
		return fmt.Errorf("record set JSON must include zoneId and id")
	}

	client := client(c)
	updated, err := client.RecordSetUpdate(rs)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(updated)
	}

	fmt.Printf("Updated record set %s\n", rs.ID)
	return nil
}

func recordSetCount(c *cli.Context) error {
	client := client(c)
	zoneID, err := getZoneID(client, c.String("zone-id"), c.String("zone-name"))
	if err != nil {
		return err
	}

	count, err := client.RecordSetCount(zoneID)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(count)
	}

	data := [][]string{
		{"Count", strconv.Itoa(count.Count)},
	}
	printBasicTable(data)
	return nil
}

func recordSetHistory(c *cli.Context) error {
	client := client(c)
	zoneID, err := getZoneID(client, c.String("zone-id"), c.String("zone-name"))
	if err != nil {
		return err
	}

	fqdn, err := getOption(c, "record-set-fqdn")
	if err != nil {
		return err
	}
	recordType, err := getOption(c, "record-set-type")
	if err != nil {
		return err
	}

	filter := vinyldns.RecordSetChangeHistoryFilter{
		ZoneID:     zoneID,
		FQDN:       fqdn,
		RecordType: recordType,
		StartFrom:  c.String("start-from"),
	}
	maxItems, err := parseIntFlag(c, "max-items")
	if err != nil {
		return err
	}
	filter.MaxItems = maxItems

	history, err := client.RecordSetChangeHistory(filter)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(history)
	}

	for _, change := range history.RecordSetChanges {
		clitable.PrintHorizontal(map[string]interface{}{
			"Zone":          change.Zone.Name,
			"RecordSetName": change.RecordSet.Name,
			"RecordSetID":   change.RecordSet.ID,
			"UserID":        change.UserID,
			"ChangeType":    change.ChangeType,
			"Status":        change.Status,
			"Created":       change.Created,
			"ID":            change.ID,
		})
	}

	return nil
}

func recordSetChangesFailure(c *cli.Context) error {
	client := client(c)
	zoneID, err := getZoneID(client, c.String("zone-id"), c.String("zone-name"))
	if err != nil {
		return err
	}

	filter := vinyldns.ListFilter{
		StartFrom: c.String("start-from"),
	}
	maxItems, err := parseIntFlag(c, "max-items")
	if err != nil {
		return err
	}
	filter.MaxItems = maxItems

	failures, err := client.RecordSetChangesFailure(zoneID, filter)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(failures)
	}

	for _, change := range failures.FailedRecordSetChanges {
		clitable.PrintHorizontal(map[string]interface{}{
			"Zone":          change.Zone.Name,
			"RecordSetName": change.RecordSet.Name,
			"RecordSetID":   change.RecordSet.ID,
			"UserID":        change.UserID,
			"ChangeType":    change.ChangeType,
			"Status":        change.Status,
			"Created":       change.Created,
			"ID":            change.ID,
		})
	}

	return nil
}
