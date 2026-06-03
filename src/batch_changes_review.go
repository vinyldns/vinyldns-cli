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
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func batchChangeApprove(c *cli.Context) error {
	return batchChangeReview(c, "approve")
}

func batchChangeReject(c *cli.Context) error {
	return batchChangeReview(c, "reject")
}

func batchChangeCancel(c *cli.Context) error {
	return batchChangeReview(c, "cancel")
}

func batchChangeReview(c *cli.Context, action string) error {
	changeID, err := getOption(c, "batch-change-id")
	if err != nil {
		return err
	}

	var review *vinyldns.BatchChangeReview
	if comment := c.String("review-comment"); comment != "" {
		review = &vinyldns.BatchChangeReview{ReviewComment: comment}
	}

	client := client(c)
	var resp *vinyldns.BatchRecordChange
	switch action {
	case "approve":
		resp, err = client.BatchRecordChangeApprove(changeID, review)
	case "reject":
		resp, err = client.BatchRecordChangeReject(changeID, review)
	case "cancel":
		resp, err = client.BatchRecordChangeCancel(changeID, review)
	default:
		return fmt.Errorf("unknown review action %q", action)
	}
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(resp)
	}

	data := [][]string{
		{"ID", resp.ID},
		{"Status", resp.Status},
		{"ApprovalStatus", resp.ApprovalStatus},
		{"ReviewComment", resp.ReviewComment},
		{"ReviewerUserName", resp.ReviewerUserName},
	}

	printBasicTable(data)
	return nil
}
