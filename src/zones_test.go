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
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

var _ = Describe("its commands for working with zones", func() {
	var (
		session   *gexec.Session
		err       error
		args      []string
		zonesArgs []string
		name      string = "vinyldns."
		group     *vinyldns.Group
		makeGroup = func() *vinyldns.Group {
			return &vinyldns.Group{
				Name:        "zones-test-group",
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
			var zones []vinyldns.Zone

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

	JustBeforeEach(func() {
		args = append(baseArgs, zonesArgs...)
		cmd := exec.Command(exe, args...)
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	JustAfterEach(func() {
		session.Terminate().Wait()
	})

	Describe("its 'zones' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				zonesArgs = []string{
					"zones",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("List all vinyldns zones"))
			})
		})

		Context("when no zones exist", func() {
			Context("when not passed an --output", func() {
				BeforeEach(func() {
					zonesArgs = []string{
						"zones",
					}
				})

				It("prints the correct data", func() {
					Eventually(session.Out, 5).Should(gbytes.Say("No zones found"))
				})
			})

			Context("when passed an --output=json", func() {
				BeforeEach(func() {
					zonesArgs = []string{
						"--output=json",
						"zones",
					}
				})

				It("prints the correct data", func() {
					Eventually(session.Out, 5).Should(gbytes.Say(`\[\]`))
				})
			})
		})

		Context("when zones exist", func() {
			var (
				zone *vinyldns.ZoneUpdateResponse
				name string = "vinyldns."
			)

			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup())
				Expect(err).NotTo(HaveOccurred())

				zone, err = vinylClient.ZoneCreate(makeZone(name, group.ID))
				Expect(err).NotTo(HaveOccurred())

				// wait to be sure the zone is fully created
				// TODO: this can be improved
				time.Sleep(3 * time.Second)
			})

			AfterEach(func() {
				_, err = vinylClient.ZoneDelete(zone.Zone.ID)
				Expect(err).NotTo(HaveOccurred())

				for {
					exists, err := vinylClient.ZoneExists(zone.Zone.ID)
					Expect(err).NotTo(HaveOccurred())

					if !exists {
						break
					}
				}

				_, err = vinylClient.GroupDelete(group.ID)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when it's not passed the --output=json option", func() {
				BeforeEach(func() {
					zonesArgs = []string{
						"zones",
					}
				})

				It("prints zone details", func() {
					output := fmt.Sprintf(`+-----------+--------------------------------------+
|   NAME    |                  ID                  |
+-----------+--------------------------------------+
| vinyldns. | %s |
+-----------+--------------------------------------+`, zone.Zone.ID)

					Eventually(func() string {
						return string(session.Out.Contents())
					}).Should(ContainSubstring(output))
				})
			})
		})
	})

	Describe("its 'zone' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				zonesArgs = []string{
					"zone",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("view zone details"))
			})
		})

		Context("when the zone exists", func() {
			var (
				zone *vinyldns.ZoneUpdateResponse
			)

			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup())
				Expect(err).NotTo(HaveOccurred())

				zone, err = vinylClient.ZoneCreate(makeZone(name, group.ID))
				Expect(err).NotTo(HaveOccurred())

				// wait to be sure the zone is fully created
				time.Sleep(3 * time.Second)
			})

			AfterEach(func() {
				_, err = vinylClient.ZoneDelete(zone.Zone.ID)
				Expect(err).NotTo(HaveOccurred())

				for {
					exists, err := vinylClient.ZoneExists(zone.Zone.ID)
					Expect(err).NotTo(HaveOccurred())

					if !exists {
						break
					}
				}

				_, err = vinylClient.GroupDelete(group.ID)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("it's passed a '--zone-name'", func() {
				BeforeEach(func() {
					zonesArgs = []string{
						"zone",
						fmt.Sprintf("--zone-name=%s", name),
					}
				})

				It("prints the zone's details", func() {
					output := fmt.Sprintf(`+--------+--------------------------------------+
| Name   | %s                            |
+--------+--------------------------------------+
| ID     | %s |
+--------+--------------------------------------+
| Status | Active                               |
+--------+--------------------------------------+`, name, zone.Zone.ID)

					Eventually(func() string {
						return string(session.Out.Contents())
					}).Should(ContainSubstring(output))

				})
			})

			Context("it's passed a '--zone-id'", func() {
				BeforeEach(func() {
					zonesArgs = []string{
						"zone",
						fmt.Sprintf("--zone-id=%s", zone.Zone.ID),
					}
				})

				It("prints the zone's details", func() {
					output := fmt.Sprintf(`+--------+--------------------------------------+
| Name   | %s                            |
+--------+--------------------------------------+
| ID     | %s |
+--------+--------------------------------------+
| Status | Active                               |
+--------+--------------------------------------+`, name, zone.Zone.ID)

					Eventually(func() string {
						return string(session.Out.Contents())
					}).Should(ContainSubstring(output))
				})
			})
		})
	})

	Describe("its 'zone-create' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				zonesArgs = []string{
					"zone-create",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Create a zone"))
			})
		})

		Context("when it's not passed connection details", func() {
			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup())
				Expect(err).NotTo(HaveOccurred())

				zonesArgs = []string{
					"zone-create",
					fmt.Sprintf("--name=%s", name),
					"--email=admin@test.com",
					fmt.Sprintf("--admin-group-name=%s", group.Name),
				}
			})

			AfterEach(func() {
				deleteAllGroupsAndZones()
			})

			It("prints a message reporting that the zone has been created", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Created zone vinyldns."))
			})
		})

		Context("when it's passed valid connection details", func() {
			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup())
				Expect(err).NotTo(HaveOccurred())

				zonesArgs = []string{
					"zone-create",
					fmt.Sprintf("--name=%s", name),
					"--email=admin@test.com",
					fmt.Sprintf("--admin-group-name=%s", group.Name),
					"--zone-connection-key-name=vinyldns.",
					"--zone-connection-key=nzisn+4G2ldMn0q1CV3vsg==",
					"--zone-connection-primary-server=vinyldns-bind9",
					fmt.Sprintf("--transfer-connection-key-name=%s", name),
					"--transfer-connection-key=nzisn+4G2ldMn0q1CV3vsg==",
					"--transfer-connection-primary-server=vinyldns-bind9",
				}
			})

			AfterEach(func() {
				deleteAllGroupsAndZones()
			})

			It("prints a message reporting that the zone has been created", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Created zone vinyldns."))
			})
		})

		Context("when it's passed invalid connection details", func() {
			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup())
				Expect(err).NotTo(HaveOccurred())

				zonesArgs = []string{
					"zone-create",
					fmt.Sprintf("--name=%s", name),
					"--email=admin@test.com",
					fmt.Sprintf("--admin-group-name=%s", group.Name),
					"--zone-connection-key=nzisn+4G2ldMn0q1CV3vsg==",
					"--zone-connection-primary-server=vinyldns-bind9",
				}
			})

			AfterEach(func() {
				deleteAllGroups()
			})

			It("prints an explanatory message to stderr", func() {
				Eventually(session.Err, 5).Should(gbytes.Say("zone connection requires '--zone-connection-key-name', '--zone-connection-key', and '--zone-connection-primary-server'"))
			})

			It("exits 1", func() {
				Eventually(session, 3).Should(gexec.Exit(1))
			})
		})

		Context("when it's passed invalid transfer connection details", func() {
			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup())
				Expect(err).NotTo(HaveOccurred())

				zonesArgs = []string{
					"zone-create",
					fmt.Sprintf("--name=%s", name),
					"--email=admin@test.com",
					fmt.Sprintf("--admin-group-name=%s", group.Name),
					"--transfer-connection-key=nzisn+4G2ldMn0q1CV3vsg==",
					"--transfer-connection-primary-server=vinyldns-bind9",
				}
			})

			AfterEach(func() {
				deleteAllGroups()
			})

			It("prints an explanatory message to stderr", func() {
				Eventually(session.Err, 5).Should(gbytes.Say("transfer connection requires '--transfer-connection-key-name', '--transfer-connection-key', and '--transfer-connection-primary-server'"))
			})

			It("exits 1", func() {
				Eventually(session, 3).Should(gexec.Exit(1))
			})
		})
	})

	Describe("its 'zone-update' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				zonesArgs = []string{
					"zone-update",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("update zone details"))
			})
		})

		Context("when it's passed a JSON string", func() {
			var (
				zone     *vinyldns.ZoneUpdateResponse
				newEmail string = "updated@email.com"
			)

			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup())
				Expect(err).NotTo(HaveOccurred())

				zone, err = vinylClient.ZoneCreate(makeZone(name, group.ID))
				Expect(err).NotTo(HaveOccurred())

				// poll until the new zone exists
				for {
					exists, err := vinylClient.ZoneExists(zone.Zone.ID)
					Expect(err).NotTo(HaveOccurred())

					if exists {
						break
					}
				}

				zone.Zone.Email = newEmail
				j, err := json.Marshal(zone.Zone)
				Expect(err).NotTo(HaveOccurred())

				zonesArgs = []string{
					"zone-update",
					"--json",
					string(j),
				}
			})

			AfterEach(func() {
				deleteAllGroupsAndZones()
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Updated zone vinyldns."))
			})

			It("updates the zone", func() {
				z, err := vinylClient.Zone(zone.Zone.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(z.Email).NotTo(Equal(newEmail))
			})
		})
	})
})
