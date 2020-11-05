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
	)

	JustBeforeEach(func() {
		args = append(baseArgs, zonesArgs...)
		cmd := exec.Command(exe, args...)
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	})

	JustAfterEach(func() {
		if session != nil {
			session.Terminate()
		}
	})

	Describe("its 'zones' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				zonesArgs = []string{
					"zones",
					"--help",
				}
			})

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
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

				It("does not error", func() {
					Expect(err).NotTo(HaveOccurred())
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

				It("does not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("prints the correct data", func() {
					Eventually(session.Out, 5).Should(gbytes.Say(`\[\]`))
				})
			})
		})

		Context("when zones exist", func() {
			var (
				zone  *vinyldns.ZoneUpdateResponse
				group *vinyldns.Group
				name  string = "vinyldns."
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

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("view zone details"))
			})
		})

		Context("when the zone exists", func() {
			var (
				zone  *vinyldns.ZoneUpdateResponse
				group *vinyldns.Group
				name  string = "vinyldns."
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

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Create a zone"))
			})
		})
	})
})
