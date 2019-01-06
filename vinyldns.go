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
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/vinyldns/go-vinyldns/vinyldns"

	clitable "github.com/crackcomm/go-clitable"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

// passed in via Makefile
var version string

const hostFlag = "host"
const accessKeyFlag = "access-key"
const secretKeyFlag = "secret-key"

func main() {
	app := cli.NewApp()
	app.Name = "vinyldns"
	app.Version = version
	app.Usage = "A CLI to the vinyldns DNS-as-a-service API"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   hostFlag,
			Usage:  "vinyldns API Hostname",
			EnvVar: "VINYLDNS_HOST",
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s, ak", accessKeyFlag),
			Usage:  "vinyldns access key",
			EnvVar: "VINYLDNS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s, sk", secretKeyFlag),
			Usage:  "vinyldns secret key",
			EnvVar: "VINYLDNS_SECRET_KEY",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:        "groups",
			Usage:       "groups",
			Description: "List all vinyldns groups",
			Action:      groups,
		},
		{
			Name:        "group",
			Usage:       "group --group-id <groupID>",
			Description: "Retrieve details for vinyldns group",
			Action:      group,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "group-id",
					Usage: "The group ID",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "The group name (in alternative to group-id)",
				},
			},
		},
		{
			Name:        "group-create",
			Usage:       "group-create --json <groupJSON>",
			Description: "Create a vinyldns group",
			Action:      groupCreate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "The vinyldns JSON representing the group",
				},
			},
		},
		{
			Name:        "group-delete",
			Usage:       "group-delete --group-id <groupID>",
			Description: "Delete the targeted vinyldns group",
			Action:      groupDelete,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "group-id",
					Usage: "The group ID",
				},
			},
		},
		{
			Name:        "group-admins",
			Usage:       "group-admins --group-id <groupID>",
			Description: "Retrieve details for vinyldns group admins",
			Action:      groupAdmins,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "group-id",
					Usage: "The group ID",
				},
			},
		},
		{
			Name:        "group-members",
			Usage:       "group-members --group-id <groupID>",
			Description: "Retrieve details for vinyldns group members",
			Action:      groupMembers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "group-id",
					Usage: "The group ID",
				},
			},
		},
		{
			Name:        "group-activity",
			Usage:       "group-activity --group-id <groupID>",
			Description: "Retrieve change activity details for vinyldns group activity",
			Action:      groupActivity,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "group-id",
					Usage: "The group ID",
				},
			},
		},
		{
			Name:        "zones",
			Usage:       "zones",
			Description: "List all vinyldns zones",
			Action:      zones,
		},
		{
			Name:        "zone",
			Usage:       "zone --zone-id <zoneID>",
			Description: "view zone details",
			Action:      zone,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
			},
		},
		{
			Name:        "zone-create",
			Usage:       "zone-create --name <name> --email <email> --admin-group-id <adminGroupID> --transfer-connection-name <transferConnectionName> --transfer-connection-key <transferConnectionKey> --transfer-connection-key-name <transferConnectionKeyName> --transfer-connection-primary-server <transferConnectionPrimaryServer> --zone-connection-name <zoneConnectionName> --zone-connection-key <zoneConnectionKey> --zone-connection-key-name <zoneConnectionKeyName> --zone-connection-primary-server <zoneConnectionPrimaryServer>",
			Description: "Create a zone",
			Action:      zoneCreate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The zone name",
				},
				cli.StringFlag{
					Name:  "email",
					Usage: "The zone email",
				},
				cli.StringFlag{
					Name:  "admin-group-id",
					Usage: "The zone admin group ID",
				},
				cli.StringFlag{
					Name:  "admin-group-name",
					Usage: "The zone admin group name (an alternative to admin-group-id)",
				},
				cli.StringFlag{
					Name:  "transfer-connection-key-name",
					Usage: "The zone transfer connection key name",
				},
				cli.StringFlag{
					Name:  "transfer-connection-key",
					Usage: "The zone transfer connection key",
				},
				cli.StringFlag{
					Name:  "transfer-connection-primary-server",
					Usage: "The zone transfer connection primary server",
				},
				cli.StringFlag{
					Name:  "zone-connection-key-name",
					Usage: "The zone connection key name",
				},
				cli.StringFlag{
					Name:  "zone-connection-key",
					Usage: "The zone connection key",
				},
				cli.StringFlag{
					Name:  "zone-connection-primary-server",
					Usage: "The zone zone connection primary server",
				},
			},
		},
		{
			Name:        "zone-delete",
			Usage:       "zone-delete --zone-id <zoneID>",
			Description: "Delete a zone",
			Action:      zoneDelete,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "The zone ID",
				},
			},
		},
		{
			Name:        "zone-connection",
			Usage:       "zone-connection --zone-id <zoneID>",
			Description: "view zone connection details",
			Action:      zoneConnection,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
			},
		},
		{
			Name:        "zone-changes",
			Usage:       "zone-changes --zone-changes <zoneID>",
			Description: "view zone change history details",
			Action:      zoneChanges,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
			},
		},
		{
			Name:        "record-set-changes",
			Usage:       "record-set-changes --zone-id <zoneID>",
			Description: "view record set change history details for a zone",
			Action:      recordSetChanges,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
			},
		},
		{
			Name:        "record-set",
			Usage:       "record-set --zone-id <zoneID> --record-set-id <recordSetID>",
			Description: "View record set details",
			Action:      recordSet,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
				cli.StringFlag{
					Name:  "record-set-id",
					Usage: "The record set ID",
				},
			},
		},
		{
			Name:        "record-set-change",
			Usage:       "record-set-change --zone-id <zoneID> --record-set-id <recordSetID> --change-id <changeID>",
			Description: "view record set change details for a zone",
			Action:      recordSetChange,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
				cli.StringFlag{
					Name:  "record-set-id",
					Usage: "The record set ID",
				},
				cli.StringFlag{
					Name:  "change-id",
					Usage: "The change ID",
				},
			},
		},
		{
			Name:        "record-set-create",
			Usage:       "record-set-create --zone-id <zoneID> --record-set-name <recordSetName> --record-set-type <type> --record-set-ttl <TTL> --record-set-data <rdata>",
			Description: "add a record set in a zone",
			Action:      recordSetCreate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
				cli.StringFlag{
					Name:  "record-set-name",
					Usage: "The record set name",
				},
				cli.StringFlag{
					Name:  "record-set-type",
					Usage: "The record set type",
				},
				cli.StringFlag{
					Name:  "record-set-ttl",
					Usage: "The record set TTL",
				},
				cli.StringFlag{
					Name:  "record-set-data",
					Usage: "The record set data",
				},
			},
		},
		{
			Name:        "record-set-delete",
			Usage:       "record-set-delete --zone-id <zoneID> --record-set-id <recordSetID>",
			Description: "delete record set in a zone",
			Action:      recordSetDelete,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
				cli.StringFlag{
					Name:  "record-set-id",
					Usage: "The record set ID",
				},
			},
		},
		{
			Name:        "record-sets",
			Usage:       "record-sets --zone-id <zoneID>",
			Description: "List all record sets associated with a zone",
			Action:      recordSets,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
			},
		},
		{
			Name:        "batch-changes",
			Usage:       "batch-changes",
			Description: "List all batch changes",
			Action:      batchChanges,
		},
		{
			Name:        "batch-change",
			Usage:       "batch-change --batch-change-id <batchChangeID>",
			Description: "view batch change details for a particular batch-id",
			Action:      batchChange,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "batch-change-id",
					Usage: "The batch change ID",
				},
			},
		},
	}
	app.RunAndExitOnError()
}

func groups(c *cli.Context) error {
	client := client(c)
	validateEnv(c)
	groups, err := client.Groups()
	if err != nil {
		return err
	}

	data := [][]string{}
	for _, g := range groups {
		data = append(data, []string{
			g.Name,
			g.ID,
		})
	}

	if len(data) != 0 {
		printTableWithHeaders([]string{"Name", "ID"}, data)
	} else {
		fmt.Printf("No groups found")
	}

	return nil
}

func group(c *cli.Context) error {
	id := c.String("group-id")
	name := c.String("name")
	g, err := getGroup(client(c), name, id)
	if err != nil {
		return err
	}

	data := [][]string{
		{"Name", g.Name},
		{"ID", g.ID},
		{"Email", g.Email},
		{"Description", g.Description},
		{"Status", g.Status},
		{"Members", userIDList(g.Members)},
		{"Admins", userIDList(g.Admins)},
	}

	printBasicTable(data)

	return nil
}

func getGroup(c *vinyldns.Client, name, id string) (*vinyldns.Group, error) {
	if name != "" {
		fmt.Println(fmt.Printf("here: %s", name))
		return groupByName(c, name)
	}

	return c.Group(id)
}

func groupByName(c *vinyldns.Client, name string) (*vinyldns.Group, error) {
	var g *vinyldns.Group
	groups, err := c.Groups()
	if err != nil {
		return g, err
	}

	for _, group := range groups {
		if group.Name == name {
			return &group, nil
		}
	}

	return g, fmt.Errorf("Group %s not found", name)
}

func groupCreate(c *cli.Context) error {
	data := []byte(c.String("json"))
	group := &vinyldns.Group{}
	if err := json.Unmarshal(data, &group); err != nil {
		return err
	}
	client := client(c)
	_, err := client.GroupCreate(group)
	if err != nil {
		return err
	}

	fmt.Printf("Created group %s\n", group.Name)

	return nil
}

func groupDelete(c *cli.Context) error {
	id := c.String("group-id")
	client := client(c)
	_, err := client.GroupDelete(id)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted group %s\n", id)

	return nil
}

func groupAdmins(c *cli.Context) error {
	client := client(c)
	admins, err := client.GroupAdmins(c.String("group-id"))
	if err != nil {
		return err
	}

	printUsers(admins)

	return nil
}

func groupMembers(c *cli.Context) error {
	client := client(c)
	admins, err := client.GroupMembers(c.String("group-id"))
	if err != nil {
		return err
	}

	printUsers(admins)

	return nil
}

func groupActivity(c *cli.Context) error {
	client := client(c)
	activity, err := client.GroupActivity(c.String("group-id"))
	if err != nil {
		return err
	}

	for _, c := range activity.Changes {
		fmt.Println("Changes...")

		data := [][]string{
			{"Created", c.Created},
			{"UserID", c.UserID},
			{"ChangeType", c.ChangeType},
		}

		printBasicTable(data)

		fmt.Println("\n\nNew Group...")
		printGroup(c.NewGroup)

		fmt.Println("\n\nOld Group...")
		printGroup(c.OldGroup)
		fmt.Println("\n=====")
	}

	return nil
}

func zones(c *cli.Context) error {
	client := client(c)
	validateEnv(c)
	zones, err := client.Zones()
	if err != nil {
		return err
	}

	data := [][]string{}
	for _, z := range zones {
		data = append(data, []string{
			z.Name,
			z.ID,
		})
	}

	if len(data) != 0 {
		printTableWithHeaders([]string{"Name", "ID"}, data)
	} else {
		fmt.Printf("No zones found")
	}

	return nil
}

func zone(c *cli.Context) error {
	client := client(c)
	z, err := client.Zone(c.String("zone-id"))
	if err != nil {
		return err
	}

	data := [][]string{
		{"Name", z.Name},
		{"ID", z.ID},
		{"Status", z.Status},
	}

	printBasicTable(data)

	return nil
}

func zoneDelete(c *cli.Context) error {
	id := c.String("zone-id")
	client := client(c)
	_, err := client.ZoneDelete(id)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted zone %s\n", id)

	return nil
}

func zoneCreate(c *cli.Context) error {
	client := client(c)
	id, err := getAdminGroupID(client, c.String("admin-group-id"), c.String("admin-group-name"))
	if err != nil {
		return err
	}
	connection := &vinyldns.ZoneConnection{
		Key:           c.String("zone-connection-key"),
		KeyName:       c.String("zone-connection-key-name"),
		Name:          c.String("zone-connection-key-name"),
		PrimaryServer: c.String("zone-connection-primary-server"),
	}
	tConnection := &vinyldns.ZoneConnection{
		Key:           c.String("transfer-connection-key"),
		KeyName:       c.String("transfer-connection-key-name"),
		Name:          c.String("transfer-connection-key-name"),
		PrimaryServer: c.String("transfer-connection-primary-server"),
	}
	z := &vinyldns.Zone{
		Name:         c.String("name"),
		Email:        c.String("email"),
		AdminGroupID: id,
	}

	zc, err := validateConnection("zone", connection)
	if err != nil {
		return err
	}
	if zc {
		z.Connection = connection
	}

	tc, err := validateConnection("transfer", tConnection)
	if err != nil {
		return err
	}
	if tc {
		z.TransferConnection = tConnection
	}

	created, err := client.ZoneCreate(z)
	if err != nil {
		return err
	}

	fmt.Printf("Created zone %s\n", created.Zone.Name)

	return nil
}

func getAdminGroupID(c *vinyldns.Client, id, name string) (string, error) {
	if id != "" {
		return id, nil
	}

	g, err := groupByName(c, name)
	if err != nil {
		return "", err
	}

	return g.ID, nil
}

func zoneConnection(c *cli.Context) error {
	client := client(c)
	id := c.String("zone-id")
	z, err := client.Zone(id)
	if err != nil {
		return err
	}
	con := z.Connection

	if con == nil {
		fmt.Printf("No zone connection found for zone %s\n", id)

		return nil
	}

	data := [][]string{
		{"Name", con.Name},
		{"KeyName", con.KeyName},
		{"Key", con.Key},
		{"PrimaryServer", con.PrimaryServer},
	}

	printBasicTable(data)

	return nil
}

func zoneChanges(c *cli.Context) error {
	client := client(c)
	zh, err := client.ZoneHistory(c.String("zone-id"))
	if err != nil {
		return err
	}
	cs := zh.ZoneChanges

	for _, c := range cs {
		clitable.PrintHorizontal(map[string]interface{}{
			"Zone":       c.Zone.Name,
			"ZoneID":     c.Zone.ID,
			"UserID":     c.UserID,
			"ChangeType": c.ChangeType,
			"Status":     c.Status,
			"Created":    c.Created,
			"ID":         c.ID,
		})
	}

	return nil
}

func recordSetChanges(c *cli.Context) error {
	client := client(c)
	zh, err := client.ZoneHistory(c.String("zone-id"))
	if err != nil {
		return err
	}
	rsc := zh.RecordSetChanges

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
	rs, err := client.RecordSets(c.String("zone-id"))
	if err != nil {
		return err
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

func batchChanges(c *cli.Context) error {
	client := client(c)
	rc, err := client.BatchRecordChanges()
	if err != nil {
		return err
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

func client(c *cli.Context) *vinyldns.Client {
	return &vinyldns.Client{
		AccessKey:  c.GlobalString(accessKeyFlag),
		SecretKey:  c.GlobalString(secretKeyFlag),
		Host:       c.GlobalString(hostFlag),
		HTTPClient: &http.Client{},
	}
}

func typeSwitch(t string) string {
	switch t {
	case "A", "a":
		return "A"
	case "AAAA", "aaaa":
		return "AAAA"
	case "CNAME", "cname":
		return "CNAME"
	}
	return ""
}

func getOption(c *cli.Context, name string) (string, error) {
	val := c.String(name)
	var err error
	if len(val) == 0 {
		err = fmt.Errorf("--%s is required", name)
	}
	return val, err
}

func recordSetCreate(c *cli.Context) error {
	client := client(c)
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
	rs := &vinyldns.RecordSet{
		ZoneID: c.String("zone-id"),
		Name:   c.String("record-set-name"),
		Type:   t,
		TTL:    c.Int("record-set-ttl"),
		Records: []vinyldns.Record{
			{
				Address: rdata[0],
			},
		},
	}

	_, err = client.RecordSetCreate(c.String("zone-id"), rs)
	if err != nil {
		return err
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
	_, err := client.RecordSetDelete(c.String("zone-id"), id)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted record set %s\n", id)
	return nil
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

func getRecordValue(records []string, recordValue interface{}, recordPrepend string) []string {
	var strVal string
	switch recordValue.(type) {
	case int:
		strVal = strconv.Itoa(recordValue.(int))
	case string:
		strVal = recordValue.(string)
	default:
		strVal = ""
	}
	if strVal != "" && strVal != "0" {
		records = append(records, recordPrepend+": "+strVal)
	}

	return records
}

func userIDList(mems []vinyldns.User) string {
	members := []string{}

	for _, m := range mems {
		members = append(members, m.ID)
	}

	return strings.Join(members, ", ")
}

func printUsers(users []vinyldns.User) {
	for _, u := range users {
		data := [][]string{
			{"UserName", u.UserName},
			{"Name", u.FirstName + " " + u.LastName},
			{"ID", u.ID},
			{"Email", u.Email},
			{"Created", u.Created},
		}

		printBasicTable(data)
	}
}

func printGroup(group vinyldns.Group) {
	data := [][]string{
		{"Name", group.Name},
		{"Status", group.Status},
		{"Created", group.Created},
		{"ID", group.ID},
	}

	printBasicTable(data)
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

func validateConnection(connection string, c *vinyldns.ZoneConnection) (bool, error) {
	// if all are empty, we assume the user does not want to declare a connection
	if c.Key == "" && c.KeyName == "" && c.Name == "" && c.PrimaryServer == "" {
		return false, nil
	}

	// if any but not all are empty, we have a problem
	if c.Key == "" || c.KeyName == "" || c.Name == "" || c.PrimaryServer == "" {
		return false, fmt.Errorf("%s connection requires '--%s-connection-key-name', '--%s-connection-key', and '--%s-connection-primary-server'", connection, connection, connection, connection)
	}

	return true, nil
}
