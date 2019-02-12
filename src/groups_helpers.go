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
	"strings"

	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func getGroup(c *vinyldns.Client, name, id string) (*vinyldns.Group, error) {
	if name != "" {
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
