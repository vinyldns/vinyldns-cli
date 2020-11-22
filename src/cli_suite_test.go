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
	cleanUp = func(deleteZones bool) {
		var (
			zones []vinyldns.Zone
			err   error
		)

		// poll until zones created by the tests are completely created
		for {
			if !deleteZones {
				break
			}

			zones, err = vinylClient.Zones()
			Expect(err).NotTo(HaveOccurred())

			if len(zones) != 0 {
				break
			}
		}

		for _, z := range zones {
			_, err = vinylClient.ZoneDelete(z.ID)
			Expect(err).NotTo(HaveOccurred())
		}

		// poll until all zones are deleted
		for {
			zones, err = vinylClient.Zones()
			Expect(err).NotTo(HaveOccurred())

			if len(zones) == 0 {
				break
			}
		}

		// There's a window of time following zone deletion in which
		// VinylDNS continues to believe the group is a zone admin.
		// We sleep for 3 seconds to allow VinylDNS to get itself straight.
		time.Sleep(3 * time.Second)

		var groups []vinyldns.Group
		groups, err = vinylClient.Groups()
		Expect(err).NotTo(HaveOccurred())

		for _, g := range groups {
			_, err = vinylClient.GroupDelete(g.ID)
			Expect(err).NotTo(HaveOccurred())
		}

		// poll until all groups are deleted
		for {
			groups, err := vinylClient.Groups()
			Expect(err).NotTo(HaveOccurred())

			if len(groups) == 0 {
				break
			}
		}
	}
	deleteAllGroupsAndZones = func() {
		cleanUp(true)
	}
	deleteAllGroups = func() {
		cleanUp(false)
	}
)

func TestVinylDNSCLI(t *testing.T) {
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

	// ensure there are no pre-existing groups
	gs, err := vinylClient.Groups()
	if err != nil {
		t.Error(err)
	}

	for _, g := range gs {
		_, err := vinylClient.GroupDelete(g.ID)
		if err != nil {
			t.Error(err)
		}
	}

	// ensure there are no pre-existing zones
	zs, err := vinylClient.Zones()
	if err != nil {
		t.Error(err)
	}

	for _, z := range zs {
		_, err := vinylClient.ZoneDelete(z.ID)
		if err != nil {
			t.Error(err)
		}
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "vinyldns CLI integration test suite")
}
