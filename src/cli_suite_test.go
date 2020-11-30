// +build integration

/*
Copyright 2020 Comcast Cable Communications Management, LLC
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
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

var (
	exe         string
	baseArgs    []string
	vinylClient *vinyldns.Client
	makeGroup   = func(name string) *vinyldns.Group {
		return &vinyldns.Group{
			Name:        name,
			Description: "description",
			Email:       "email@email.com",
			Admins: []vinyldns.User{{
				UserName: "ok",
				ID:       "ok",
			}},
			Members: []vinyldns.User{{
				UserName: "ok",
				ID:       "ok",
			}},
		}
	}
	makeZone = func(name, adminGroupID string) *vinyldns.Zone {
		return &vinyldns.Zone{
			Name:         name,
			Email:        "email@email.com",
			AdminGroupID: adminGroupID,
		}
	}
	deleteGroupsAndZones = func(waitForZoneCreation bool) error {
		var (
			zones []vinyldns.Zone
			err   error
		)

		// Poll until zones created by the tests are completely created.
		// TODO: This could perhaps be improved, as it may lead to an infinite
		// loop if ever a zone creation is _expected_, but is unsuccessful.
		for {
			if !waitForZoneCreation {
				break
			}

			zones, err = vinylClient.Zones()
			if err != nil {
				return err
			}

			if len(zones) != 0 {
				break
			}
		}

		zones, err = vinylClient.Zones()
		if err != nil {
			return err
		}

		for _, z := range zones {
			_, err = vinylClient.ZoneDelete(z.ID)
			if err != nil {
				return err
			}
		}

		// poll until all zones are deleted
		for {
			zones, err = vinylClient.Zones()
			if err != nil {
				return err
			}

			if len(zones) == 0 {
				break
			}
		}

		// There's a window of time following zone deletion in which
		// VinylDNS continues to believe the group is a zone admin.
		// We sleep to allow VinylDNS to get itself straight.
		time.Sleep(6 * time.Second)

		var groups []vinyldns.Group
		groups, err = vinylClient.Groups()
		if err != nil {
			return err
		}

		for _, g := range groups {
			_, err = vinylClient.GroupDelete(g.ID)
			if err != nil {
				return err
			}
		}

		// poll until all groups are deleted
		for {
			groups, err := vinylClient.Groups()
			if err != nil {
				return err
			}

			if len(groups) == 0 {
				break
			}
		}

		return nil
	}
	deleteAllGroupsAndZones = func() error {
		return deleteGroupsAndZones(true)
	}
	deleteAllGroups = func() error {
		return deleteGroupsAndZones(false)
	}
	deleteRecordInZone = func(zoneID, rsName string) error {
		rss, err := vinylClient.RecordSets(zoneID)
		if err != nil {
			return err
		}

		for _, rs := range rss {
			if rs.Name == rsName {
				_, err := vinylClient.RecordSetDelete(zoneID, rs.ID)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
	createGroupAndZone = func(groupName, zoneName string) (*vinyldns.Group, *vinyldns.ZoneUpdateResponse, error) {
		group, err := vinylClient.GroupCreate(makeGroup(groupName))
		if err != nil {
			return group, nil, err
		}

		zone, err := vinylClient.ZoneCreate(makeZone(zoneName, group.ID))
		if err != nil {
			return group, zone, err
		}

		// poll until zone creation is complete
		for {
			exists, err := vinylClient.ZoneExists(zone.Zone.ID)
			if err != nil {
				return group, zone, err
			}

			if exists {
				break
			}
		}

		return group, zone, err
	}
)

func TestVinylDNSCLI(t *testing.T) {
	// The tests assume a local VinylDNS API is available on port 9000.
	host := "http://localhost:9000"
	accessKey := "okAccessKey"
	secretKey := "okSecretKey"

	exe = "../bin/vinyldns"
	baseArgs = []string{
		fmt.Sprintf("--host=%s", host),
		fmt.Sprintf("--access-key=%s", accessKey),
		fmt.Sprintf("--secret-key=%s", secretKey),
	}

	vinylClient = vinyldns.NewClient(vinyldns.ClientConfiguration{
		Host:      host,
		AccessKey: accessKey,
		SecretKey: secretKey,
	})

	err := deleteGroupsAndZones(false)
	if err != nil {
		t.Error(err)
	}

	config.DefaultReporterConfig.SlowSpecThreshold = 30
	config.GinkgoConfig.ParallelTotal = 1

	RegisterFailHandler(Fail)
	RunSpecs(t, "VinylDNS CLI integration test suite")
}
