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
	"net/http"
	"os"
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
	client := client(c)
	g, err := client.Group(c.String("id"))
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
		cname := r.CName
		address := r.Address

		if cname != "" {
			records = append(records, "CNAME: "+cname)
		}
		if address != "" {
			records = append(records, "address: "+address)
		}
	}

	return strings.Join(records, "\n")
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
