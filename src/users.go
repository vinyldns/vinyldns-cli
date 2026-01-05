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
	"fmt"
	"strings"

	"github.com/urfave/cli"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func user(c *cli.Context) error {
	identifier := c.String("user-id")
	if identifier == "" {
		identifier = c.String("user-name")
	}
	if identifier == "" {
		return fmt.Errorf("one of the flags must be provided: '--user-id', '--user-name'")
	}

	client := client(c)
	u, err := client.User(identifier)
	if err != nil {
		return err
	}

	return printUserInfo(c, u)
}

func userLock(c *cli.Context) error {
	identifier, err := getOption(c, "user-id")
	if err != nil {
		return err
	}

	client := client(c)
	u, err := client.UserLock(identifier)
	if err != nil {
		return err
	}

	return printUserInfo(c, u)
}

func userUnlock(c *cli.Context) error {
	identifier, err := getOption(c, "user-id")
	if err != nil {
		return err
	}

	client := client(c)
	u, err := client.UserUnlock(identifier)
	if err != nil {
		return err
	}

	return printUserInfo(c, u)
}

func printUserInfo(c *cli.Context, u vinyldns.UserInfo) error {
	if c.GlobalString(outputFlag) == "json" {
		return printJSON(u)
	}

	data := [][]string{
		{"ID", u.ID},
		{"UserName", u.UserName},
		{"LockStatus", u.LockStatus},
		{"GroupIDs", joinUserGroupIDs(u)},
	}

	printBasicTable(data)
	return nil
}

func joinUserGroupIDs(u vinyldns.UserInfo) string {
	ids := make([]string, 0, len(u.GroupID))
	for _, g := range u.GroupID {
		if g != "" {
			ids = append(ids, g)
		}
	}

	return strings.Join(ids, ", ")
}
