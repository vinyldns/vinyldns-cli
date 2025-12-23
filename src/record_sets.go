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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/crackcomm/go-clitable"
	"github.com/urfave/cli"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func recordSetChanges(c *cli.Context) error {
	client := client(c)
	zh, err := client.RecordSetChanges(c.String("zone-id"), vinyldns.ListFilterRecordSetChanges{})
	if err != nil {
		return err
	}
	rsc := zh.RecordSetChanges

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(rsc)
	}

	for _, c := range rsc {
		clitable.PrintHorizontal(map[string]interface{}{
			"Zone":          c.Zone.Name,
			"RecordSetName": c.RecordSet.Name,
			"RecordSetID":   c.RecordSet.ID,
			"UserID":        c.UserID,
			"ChangeType":    c.ChangeType,
			"Status":        c.Status,
			"Created":       c.Created,
			"ID":            c.ID,
		})
	}

	return nil
}

func recordSetChange(c *cli.Context) error {
	client := client(c)
	rsc, err := client.RecordSetChange(c.String("zone-id"), c.String("record-set-id"), c.String("change-id"))
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(rsc)
	}

	clitable.PrintHorizontal(map[string]interface{}{
		"Zone":          rsc.Zone.Name,
		"RecordSetName": rsc.RecordSet.Name,
		"RecordSetID":   rsc.RecordSet.ID,
		"UserID":        rsc.UserID,
		"ChangeType":    rsc.ChangeType,
		"Status":        rsc.Status,
		"Created":       rsc.Created,
		"ID":            rsc.ID,
	})

	return nil
}

func recordSets(c *cli.Context) error {
	client := client(c)
	rs, err := client.RecordSetsListAll(c.String("zone-id"), vinyldns.ListFilter{})
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(rs)
	}

	s := []map[string]interface{}{}
	for _, r := range rs {
		m := map[string]interface{}{}
		m["Name"] = r.Name
		m["ID"] = r.ID
		m["Type"] = r.Type
		m["Status"] = r.Status
		s = append(s, m)
	}

	if len(s) != 0 {
		clitable.PrintTable([]string{"Name", "ID", "Type", "Status"}, s)
	} else {
		fmt.Printf("No record sets found")
	}

	return nil
}

func searchRecordSets(c *cli.Context) error {
	client := client(c)
	filterOptions := vinyldns.GlobalListFilter{}
	recordNameFilter, err := getOption(c, "record-name-filter")
	if err != nil {
		return err
	}
	filterOptions.RecordNameFilter = recordNameFilter
	startFrom := c.String("start-from")
	if startFrom != "" {
		filterOptions.StartFrom = startFrom
	}
	maxItemsString := c.String("max-items")
	if maxItemsString != "" {
		maxItems, err := strconv.Atoi(maxItemsString)
		if err != nil {
			return err
		}
		filterOptions.MaxItems = maxItems
	}
	recordTypeFilter := c.String("record-type-filter")
	if recordTypeFilter != "" {
		filterOptions.RecordTypeFilter = recordTypeFilter
	}
	recordOwnerGroup := c.String("record-owner-group")
	if recordOwnerGroup != "" {
		filterOptions.RecordOwnerGroupFilter = recordOwnerGroup
	}
	nameSortString := c.String("name-sort")
	nameSort := vinyldns.ASC
	if nameSortString == "DESC" {
		nameSort = vinyldns.DESC
	}
	filterOptions.NameSort = nameSort

	var rs []vinyldns.RecordSet
	if filterOptions.MaxItems <= 0 {
		rs, err = client.RecordSetsGlobalListAll(filterOptions)
		if err != nil {
			return err
		}
	} else {
		rs, _, err = client.RecordSetsGlobal(filterOptions)
		if err != nil {
			return err
		}
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(rs)
	}

	s := []map[string]interface{}{}
	for _, r := range rs {
		m := map[string]interface{}{}
		m["Name"] = r.Name
		m["ID"] = r.ID
		m["Type"] = r.Type
		m["Status"] = r.Status
		s = append(s, m)
	}

	if len(s) != 0 {
		clitable.PrintTable([]string{"Name", "ID", "Type", "Status"}, s)
	} else {
		fmt.Printf("No record sets found")
	}

	return nil
}

func recordSet(c *cli.Context) error {
	client := client(c)
	rs, err := client.RecordSet(c.String("zone-id"), c.String("record-set-id"))
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(rs)
	}

	clitable.PrintHorizontal(map[string]interface{}{
		"Zone":    rs.ZoneID,
		"Name":    rs.Name,
		"Account": rs.Account,
		"ID":      rs.ID,
		"Type":    rs.Type,
		"Records": getRecord(rs.Records),
		"Created": rs.Created,
		"Status":  rs.Status,
		"Updated": rs.Updated,
		"TTL":     rs.TTL,
	})

	return nil
}

func recordSetCreate(c *cli.Context) error {
	client := client(c)
	zoneID, err := getZoneID(client, c.String("zone-id"), c.String("zone-name"))
	if err != nil {
		return err
	}

	name, err := getOption(c, "record-set-name")
	if err != nil {
		return err
	}

	rtype, err := getOption(c, "record-set-type")
	if err != nil {
		return err
	}
	t := typeSwitch(rtype)
	if len(t) == 0 {
		return fmt.Errorf("unknown --record-set-type %s", rtype)
	}

	rdataS, err := getOption(c, "record-set-data")
	if err != nil {
		return err
	}

	rdata := strings.Split(rdataS, ",")

	var records []vinyldns.Record

	if t == "CNAME" {
		records = []vinyldns.Record{
			{
				CName: rdata[0],
			},
		}
	} else if t == "MX" {
		i, err := strconv.Atoi(rdata[0])

		if err != nil {
			return err
		}

		records = []vinyldns.Record{
			{
				Preference: i,
				Exchange:   rdata[1],
			},
		}
	} else if t == "PTR" {
		records = []vinyldns.Record{
			{
				PTRDName: rdata[0],
			},
		}
	} else if t == "TXT" {
		records = []vinyldns.Record{
			{
				Text: rdataS,
			},
		}
	} else {
		records = []vinyldns.Record{
			{
				Address: rdata[0],
			},
		}
	}

	rs := &vinyldns.RecordSet{
		ZoneID:  zoneID,
		Name:    c.String("record-set-name"),
		Type:    t,
		TTL:     c.Int("record-set-ttl"),
		Records: records,
	}

	rsc, err := client.RecordSetCreate(rs)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(rsc)
	}

	fmt.Printf("Created record set %s\n", name)
	return nil
}

func recordSetDelete(c *cli.Context) error {
	id := c.String("record-set-id")
	if len(id) == 0 {
		return errors.New("--record-set-id is required")
	}

	client := client(c)
	d, err := client.RecordSetDelete(c.String("zone-id"), id)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(d)
	}

	fmt.Printf("Deleted record set %s\n", id)
	return nil
}

func getRecord(recs []vinyldns.Record) string {
	records := []string{}

	for _, r := range recs {
		records = getRecordValue(records, r.Address, "Address")
		records = getRecordValue(records, r.Algorithm, "Algorithm")
		records = getRecordValue(records, r.CName, "CNAME")
		records = getRecordValue(records, r.Exchange, "Exchange")
		records = getRecordValue(records, r.Expire, "Expire")
		records = getRecordValue(records, r.Fingerprint, "Fingerprint")
		records = getRecordValue(records, r.MName, "MNAME")
		records = getRecordValue(records, r.Minimum, "Minimum")
		records = getRecordValue(records, r.NSDName, "NSDNAME")
		records = getRecordValue(records, r.Port, "Port")
		records = getRecordValue(records, r.Preference, "Preference")
		records = getRecordValue(records, r.Priority, "Priority")
		records = getRecordValue(records, r.PTRDName, "PTRDNAME")
		records = getRecordValue(records, r.Refresh, "Refresh")
		records = getRecordValue(records, r.Retry, "Retry")
		records = getRecordValue(records, r.RName, "RNAME")
		records = getRecordValue(records, r.Serial, "Serial")
		records = getRecordValue(records, r.Target, "Target")
		records = getRecordValue(records, r.Text, "Text")
		records = getRecordValue(records, r.Type, "Type")
		records = getRecordValue(records, r.Weight, "Weight")
	}

	return strings.Join(records, "\n")
}
