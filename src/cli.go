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
	"os"
	"strings"

	"github.com/urfave/cli"
)

// passed in via Makefile
var version string

const hostFlag = "host"
const accessKeyFlag = "access-key"
const secretKeyFlag = "secret-key"
const outputFlag = "output"

func main() {
	app := cli.NewApp()
	app.Name = "vinyldns"
	app.Version = version
	app.Usage = "A CLI to the VinylDNS DNS-as-a-service API"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   hostFlag,
			Usage:  "VinylDNS API Hostname",
			EnvVar: "VINYLDNS_HOST",
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s, ak", accessKeyFlag),
			Usage:  "VinylDNS access key",
			EnvVar: "VINYLDNS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s, sk", secretKeyFlag),
			Usage:  "VinylDNS secret key",
			EnvVar: "VINYLDNS_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s, op", outputFlag),
			Usage:  "VinylDNS output format ('table' (default), 'json')",
			EnvVar: "VINYLDNS_FORMAT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:        "groups",
			Usage:       "groups",
			Description: "List all VinylDNS groups",
			Action:      groups,
		},
		{
			Name:        "group",
			Usage:       "group --group-id <groupID>",
			Description: "Retrieve details for VinylDNS group",
			Action: func(c *cli.Context) error {
				return requireAtLeast(c, group, "group-id", "name")
			},
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
			Description: "Create a VinylDNS group",
			Action:      groupCreate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "json",
					Usage:    "The VinylDNS JSON representing the group",
					Required: true,
				},
			},
		},
		{
			Name:        "group-update",
			Usage:       "group-update --json <groupJSON>",
			Description: "Update a VinylDNS group",
			Action:      groupUpdate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "json",
					Usage:    "The VinylDNS JSON representing the group",
					Required: true,
				},
			},
		},
		{
			Name:        "group-delete",
			Usage:       "group-delete --group-id <groupID>",
			Description: "Delete the targeted VinylDNS group",
			Action:      groupDelete,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "group-id",
					Usage:    "The group ID",
					Required: true,
				},
			},
		},
		{
			Name:        "group-admins",
			Usage:       "group-admins --group-id <groupID>",
			Description: "Retrieve details for VinylDNS group admins",
			Action:      groupAdmins,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "group-id",
					Usage:    "The group ID",
					Required: true,
				},
			},
		},
		{
			Name:        "group-members",
			Usage:       "group-members --group-id <groupID>",
			Description: "Retrieve details for VinylDNS group members",
			Action:      groupMembers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "group-id",
					Usage:    "The group ID",
					Required: true,
				},
			},
		},
		{
			Name:        "group-activity",
			Usage:       "group-activity --group-id <groupID>",
			Description: "Retrieve change activity details for VinylDNS group activity",
			Action:      groupActivity,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "group-id",
					Usage:    "The group ID",
					Required: true,
				},
			},
		},
		{
			Name:        "zones",
			Usage:       "zones",
			Description: "List all VinylDNS zones",
			Action:      zones,
		},
		{
			Name:        "zone",
			Usage:       "zone --zone-id <zoneID>",
			Description: "view zone details",
			Action: func(c *cli.Context) error {
				return requireAtLeast(c, zone, "zone-id", "zone-name")
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
				cli.StringFlag{
					Name:  "zone-name",
					Usage: "The zone name (an alternative to --zone-id)",
				},
			},
		},
		{
			Name:        "zone-create",
			Usage:       "zone-create --name <name> --email <email> --admin-group-id <adminGroupID> --transfer-connection-name <transferConnectionName> --transfer-connection-key <transferConnectionKey> --transfer-connection-key-name <transferConnectionKeyName> --transfer-connection-primary-server <transferConnectionPrimaryServer> --zone-connection-name <zoneConnectionName> --zone-connection-key <zoneConnectionKey> --zone-connection-key-name <zoneConnectionKeyName> --zone-connection-primary-server <zoneConnectionPrimaryServer>",
			Description: "Create a zone",
			Action: func(c *cli.Context) error {
				return requireAtLeast(c, zoneCreate, "admin-group-id", "admin-group-name")
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "name",
					Usage:    "The zone name",
					Required: true,
				},
				cli.StringFlag{
					Name:     "email",
					Usage:    "The zone email",
					Required: true,
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
			Name:        "zone-update",
			Usage:       "zone-update --json <zoneJSON>",
			Description: "update zone details",
			Action:      zoneUpdate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "json",
					Usage:    "The VinylDNS JSON representing the zone details",
					Required: true,
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
					Name:     "zone-id",
					Usage:    "The zone ID",
					Required: true,
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
					Name:     "zone-id",
					Usage:    "The zone ID",
					Required: true,
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
					Name:     "zone-id",
					Usage:    "The zone ID",
					Required: true,
				},
			},
		},
		{
			Name:        "zone-sync",
			Usage:       "zone-sync --zone-sync <zoneID>",
			Description: "starts zone sync process",
			Action: func(c *cli.Context) error {
				return requireAtLeast(c, zoneSync, "zone-id", "zone-name")
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
				cli.StringFlag{
					Name:  "zone-name",
					Usage: "The zone name (an alternative to --zone-id)",
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
					Name:     "zone-id",
					Usage:    "The zone ID",
					Required: true,
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
					Name:     "zone-id",
					Usage:    "The zone ID",
					Required: true,
				},
				cli.StringFlag{
					Name:     "record-set-id",
					Usage:    "The record set ID",
					Required: true,
				},
			},
		},
		{
			Name:        "record-set-change",
			Usage:       "record-set-change --zone-id <zoneID> --record-set-id <recordSetID> --change-id <changeID>",
			Description: "view record set change details for a zone",
			Action: func(c *cli.Context) error {
				return requireAtLeast(c, recordSetChange, "zone-id", "zone-name")
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "zone-id",
					Usage:    "The zone ID",
					Required: true,
				},
				cli.StringFlag{
					Name:     "record-set-id",
					Usage:    "The record set ID",
					Required: true,
				},
				cli.StringFlag{
					Name:     "change-id",
					Usage:    "The change ID",
					Required: true,
				},
			},
		},
		{
			Name:        "record-set-create",
			Usage:       "record-set-create --zone-id <zoneID> --record-set-name <recordSetName> --record-set-type <type> --record-set-ttl <TTL> --record-set-data <rdata>",
			Description: "add a record set in a zone",
			Action: func(c *cli.Context) error {
				return requireAtLeast(c, recordSetCreate, "zone-id", "zone-name")
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "zone-id",
					Usage: "The zone ID",
				},
				cli.StringFlag{
					Name:  "zone-name",
					Usage: "The zone name (an alternative to zone-id)",
				},
				cli.StringFlag{
					Name:     "record-set-name",
					Usage:    "The record set name",
					Required: true,
				},
				cli.StringFlag{
					Name:     "record-set-type",
					Usage:    "The record set type",
					Required: true,
				},
				cli.StringFlag{
					Name:     "record-set-ttl",
					Usage:    "The record set TTL",
					Required: true,
				},
				cli.StringFlag{
					Name:     "record-set-data",
					Usage:    "The record set data",
					Required: true,
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
					Name:     "zone-id",
					Usage:    "The zone ID",
					Required: true,
				},
				cli.StringFlag{
					Name:     "record-set-id",
					Usage:    "The record set ID",
					Required: true,
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
					Name:     "zone-id",
					Usage:    "The zone ID",
					Required: true,
				},
			},
		},
		{
			Name:        "search-record-sets",
			Usage:       "search-record-sets --record-name-filter <string>",
			Description: "List all record sets matching given record name filter",
			Action:      searchRecordSets,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "record-name-filter",
					Usage:    "Record name search string. At least two alpha-numeric characters are required.",
					Required: true,
				},
				cli.StringFlag{
					Name:     "start-from",
					Usage:    "The start key of the page.",
					Required: false,
				},
				cli.StringFlag{
					Name:     "max-items",
					Usage:    "The page limit.",
					Required: false,
				},
				cli.StringSliceFlag{
					Name:     "record-type-filter",
					Usage:    "Return record_sets whose type is present in the given list.",
					Required: false,
				},
				cli.StringFlag{
					Name:     "record-owner-group",
					Usage:    "Returns record_sets belonging to the given owner.",
					Required: false,
				},
				cli.StringFlag{
					Name:     "name-sort",
					Usage:    "Sort the results as per given order",
					Required: false,
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
					Name:     "batch-change-id",
					Usage:    "The batch change ID",
					Required: true,
				},
			},
		},
		{
			Name:        "batch-change-create",
			Usage:       "batch-change-create --json <batchChangeJSON>",
			Description: "Create a batch change",
			Action:      batchChangeCreate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "json",
					Usage:    "The VinylDNS JSON representing the batch change",
					Required: true,
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func requireAtLeast(c *cli.Context, action func(*cli.Context) error, flags ...string) error {
	haveAtLeastOne := false
	for _, flag := range flags {
		if c.String(flag) != "" {
			haveAtLeastOne = true
			break
		}
	}

	if !haveAtLeastOne {
		err := cli.ShowCommandHelp(c, c.Command.Name)
		if err != nil {
			return fmt.Errorf("error showing command help: %w", err)
		}
		prefixedFlags := make([]string, len(flags))
		for i := range flags {
			prefixedFlags[i] = fmt.Sprintf("'--%s'", flags[i])
		}
		return fmt.Errorf("one of the flags must be provided: %s", strings.Join(prefixedFlags, ", "))
	}
	return action(c)
}
