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

	"github.com/crackcomm/go-clitable"
	"github.com/urfave/cli"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func zoneBackendIDs(c *cli.Context) error {
	client := client(c)
	ids, err := client.ZoneBackendIDs()
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(ids)
	}

	data := [][]string{}
	for _, id := range ids {
		data = append(data, []string{id})
	}

	if len(data) != 0 {
		printTableWithHeaders([]string{"BackendID"}, data)
	} else {
		fmt.Print("No backend IDs found")
	}

	return nil
}

func zoneChangesFailure(c *cli.Context) error {
	filter, err := listFilterFromContext(c)
	if err != nil {
		return err
	}

	client := client(c)
	resp, err := client.ZoneChangesFailure(filter)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(resp)
	}

	for _, change := range resp.FailedZoneChanges {
		clitable.PrintHorizontal(map[string]interface{}{
			"Zone":       change.Zone.Name,
			"ZoneID":     change.Zone.ID,
			"UserID":     change.UserID,
			"ChangeType": change.ChangeType,
			"Status":     change.Status,
			"Created":    change.Created,
			"ID":         change.ID,
		})
	}

	return nil
}

func zonesDeletedChanges(c *cli.Context) error {
	filter, err := deletedZonesFilterFromContext(c)
	if err != nil {
		return err
	}

	client := client(c)
	resp, err := client.ZonesDeleted(filter)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(resp)
	}

	for _, info := range resp.ZonesDeletedInfo {
		clitable.PrintHorizontal(map[string]interface{}{
			"Zone":           info.ZoneChange.Zone.Name,
			"ZoneID":         info.ZoneChange.Zone.ID,
			"UserID":         info.ZoneChange.UserID,
			"ChangeType":     info.ZoneChange.ChangeType,
			"Status":         info.ZoneChange.Status,
			"Created":        info.ZoneChange.Created,
			"ID":             info.ZoneChange.ID,
			"AdminGroupName": info.AdminGroupName,
			"UserName":       info.UserName,
			"AccessLevel":    info.AccessLevel,
		})
	}

	return nil
}

func zoneACLRuleCreate(c *cli.Context) error {
	client := client(c)
	zoneID, err := getZoneID(client, c.String("zone-id"), c.String("zone-name"))
	if err != nil {
		return err
	}

	data := []byte(c.String("json"))
	rule := &vinyldns.ACLRule{}
	if err := json.Unmarshal(data, &rule); err != nil {
		return err
	}

	updated, err := client.ZoneACLRuleCreate(zoneID, rule)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(updated)
	}

	zoneName := updated.Zone.Name
	if zoneName == "" {
		zoneName = zoneID
	}
	fmt.Printf("Updated zone %s\n", zoneName)
	return nil
}

func zoneACLRuleDelete(c *cli.Context) error {
	client := client(c)
	zoneID, err := getZoneID(client, c.String("zone-id"), c.String("zone-name"))
	if err != nil {
		return err
	}

	data := []byte(c.String("json"))
	rule := &vinyldns.ACLRule{}
	if err := json.Unmarshal(data, &rule); err != nil {
		return err
	}

	updated, err := client.ZoneACLRuleDelete(zoneID, rule)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(updated)
	}

	zoneName := updated.Zone.Name
	if zoneName == "" {
		zoneName = zoneID
	}
	fmt.Printf("Updated zone %s\n", zoneName)
	return nil
}

func listFilterFromContext(c *cli.Context) (vinyldns.ListFilter, error) {
	filter := vinyldns.ListFilter{
		NameFilter: c.String("name-filter"),
		StartFrom:  c.String("start-from"),
	}

	maxItems, err := parseIntFlag(c, "max-items")
	if err != nil {
		return filter, err
	}
	filter.MaxItems = maxItems
	return filter, nil
}

func deletedZonesFilterFromContext(c *cli.Context) (vinyldns.DeletedZonesFilter, error) {
	filter := vinyldns.DeletedZonesFilter{
		NameFilter: c.String("name-filter"),
		StartFrom:  c.String("start-from"),
	}

	maxItems, err := parseIntFlag(c, "max-items")
	if err != nil {
		return filter, err
	}
	filter.MaxItems = maxItems

	ignoreAccess, err := parseBoolFlag(c, "ignore-access")
	if err != nil {
		return filter, err
	}
	filter.IgnoreAccess = ignoreAccess
	return filter, nil
}
