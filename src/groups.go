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

	"github.com/urfave/cli"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func groups(c *cli.Context) error {
	client := client(c)
	groups, err := client.Groups()
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(groups)
	}

	return printGroupsTable(groups)
}

func printGroupsTable(groups []vinyldns.Group) error {
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

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(g)
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

func groupCreate(c *cli.Context) error {
	data := []byte(c.String("json"))
	group := &vinyldns.Group{}
	if err := json.Unmarshal(data, &group); err != nil {
		return err
	}
	client := client(c)
	create, err := client.GroupCreate(group)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(create)
	}

	fmt.Printf("Created group %s\n", group.Name)

	return nil
}

func groupDelete(c *cli.Context) error {
	id := c.String("group-id")
	client := client(c)
	deleted, err := client.GroupDelete(id)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(deleted)
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

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(admins)
	}

	printUsers(admins)

	return nil
}

func groupMembers(c *cli.Context) error {
	client := client(c)
	members, err := client.GroupMembers(c.String("group-id"))
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(members)
	}

	printUsers(members)

	return nil
}

func groupActivity(c *cli.Context) error {
	client := client(c)
	activity, err := client.GroupActivity(c.String("group-id"))
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(activity)
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
