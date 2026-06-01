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

	"github.com/urfave/cli"
)

func recordSetOwnershipRequest(c *cli.Context) error {
	return recordSetOwnershipTransfer(c, "request")
}

func recordSetOwnershipApprove(c *cli.Context) error {
	return recordSetOwnershipTransfer(c, "approve")
}

func recordSetOwnershipReject(c *cli.Context) error {
	return recordSetOwnershipTransfer(c, "reject")
}

func recordSetOwnershipCancel(c *cli.Context) error {
	return recordSetOwnershipTransfer(c, "cancel")
}

func recordSetOwnershipTransfer(c *cli.Context, action string) error {
	client := client(c)
	zoneID, err := getZoneID(client, c.String("zone-id"), c.String("zone-name"))
	if err != nil {
		return err
	}

	recordSetID, err := getOption(c, "record-set-id")
	if err != nil {
		return err
	}

	requestedOwnerGroupID, err := getOption(c, "requested-owner-group-id")
	if err != nil {
		return err
	}

	rs, err := client.RecordSet(zoneID, recordSetID)
	if err != nil {
		return err
	}

	if rs.ZoneID == "" {
		rs.ZoneID = zoneID
	}
	if rs.ID == "" {
		rs.ID = recordSetID
	}

	switch action {
	case "request":
		_, err = client.RecordSetOwnershipTransferRequest(&rs, requestedOwnerGroupID)
	case "approve":
		_, err = client.RecordSetOwnershipTransferApprove(&rs, requestedOwnerGroupID)
	case "reject":
		_, err = client.RecordSetOwnershipTransferReject(&rs, requestedOwnerGroupID)
	case "cancel":
		_, err = client.RecordSetOwnershipTransferCancel(&rs, requestedOwnerGroupID)
	default:
		return fmt.Errorf("unknown ownership transfer action %q", action)
	}

	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON("OK")
	}

	fmt.Printf("Updated record set %s\n", recordSetID)
	return nil
}
