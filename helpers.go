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
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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
