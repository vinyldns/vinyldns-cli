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

var _ = Describe("its commands for working with record sets", func() {
	var (
		session        *gexec.Session
		err            error
		args           []string
		recordSetsArgs []string
	)

	JustBeforeEach(func() {
		args = append(baseArgs, recordSetsArgs...)
		cmd := exec.Command(exe, args...)
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	JustAfterEach(func() {
		session.Terminate()
	})

	Describe("its 'record-sets' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				recordSetsArgs = []string{
					"record-sets",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("List all record sets associated with a zone"))
			})
		})
	})

	Describe("its 'search-record-sets' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				recordSetsArgs = []string{
					"search-record-sets",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("List all record sets matching given record name filter"))
			})
		})

		Context("when the search returns no results", func() {
			BeforeEach(func() {
				recordSetsArgs = []string{
					"search-record-sets",
					"--record-name-filter=so*",
					"--record-type-filter=CNAME",
					"--record-type-filter=mx",
					"--max-items=50",
					"--name-sort=DESC",
				}
			})

			It("prints a message", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("No record sets found"))
			})
		})

		Context("when the search returns results", func() {
			var (
				err   error
				group *vinyldns.Group
				zone  *vinyldns.ZoneUpdateResponse
				rs    *vinyldns.RecordSetUpdateResponse
			)

			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup("record-sets-group"))
				Expect(err).NotTo(HaveOccurred())

				zone, err = vinylClient.ZoneCreate(makeZone("vinyldns.", group.ID))
				Expect(err).NotTo(HaveOccurred())

				// poll until zone creation is complete
				for {
					exists, err := vinylClient.ZoneExists(zone.Zone.ID)
					Expect(err).NotTo(HaveOccurred())
					if exists {
						break
					}
				}

				rs, err = vinylClient.RecordSetCreate(&vinyldns.RecordSet{
					ZoneID: zone.Zone.ID,
					Name:   "name",
					Type:   "A",
					TTL:    200,
					Records: []vinyldns.Record{{
						Address: "127.0.0.1",
					}},
				})
				Expect(err).NotTo(HaveOccurred())

				// poll until the record set creation is complete
				for {
					_, err := vinylClient.RecordSet(zone.Zone.ID, rs.RecordSet.ID)
					if err == nil {
						break
					}

					_, ok := err.(*vinyldns.Error)
					Expect(ok).To(BeTrue())
				}

				recordSetsArgs = []string{
					"search-record-sets",
					"--record-name-filter=*name*",
				}
			})

			AfterEach(func() {
				_, err := vinylClient.RecordSetDelete(zone.Zone.ID, rs.RecordSet.ID)
				Expect(err).NotTo(HaveOccurred())
				deleteAllGroupsAndZones()
			})

			It("prints the search results", func() {
				output := fmt.Sprintf(`|-------------------------------------------------------------|
| Name | ID                                   | Type | Status |
|-------------------------------------------------------------|
| name | %s | A    | Active |
|-------------------------------------------------------------|`, rs.RecordSet.ID)

				Eventually(func() string {
					return string(session.Out.Contents())
				}).Should(ContainSubstring(output))
			})
		})
	})

	Describe("its 'record-set-create' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				recordSetsArgs = []string{
					"record-set-create",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("add a record set in a zone"))
			})
		})

		Context("when tasked in creating a record", func() {
			var (
				err    error
				group  *vinyldns.Group
				zone   *vinyldns.ZoneUpdateResponse
				rsName string
				zName  string = "vinyldns."
			)

			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup("record-sets-group"))
				Expect(err).NotTo(HaveOccurred())

				zone, err = vinylClient.ZoneCreate(makeZone(zName, group.ID))
				Expect(err).NotTo(HaveOccurred())

				// poll until zone creation is complete
				for {
					exists, err := vinylClient.ZoneExists(zone.Zone.ID)
					Expect(err).NotTo(HaveOccurred())
					if exists {
						break
					}
				}

			})

			AfterEach(func() {
				rss, err := vinylClient.RecordSets(zone.Zone.ID)
				Expect(err).NotTo(HaveOccurred())

				for _, rs := range rss {
					if rs.Name == rsName {
						_, err := vinylClient.RecordSetDelete(zone.Zone.ID, rs.ID)
						Expect(err).NotTo(HaveOccurred())
					}
				}

				deleteAllGroupsAndZones()
			})

			Context("when it's tasked in creating a CNAME", func() {
				BeforeEach(func() {
					rsName = "some-cname"

					recordSetsArgs = []string{
						"record-set-create",
						fmt.Sprintf("--zone-name=%s", zName),
						fmt.Sprintf("--record-set-name=%s", rsName),
						"--record-set-type=CNAME",
						"--record-set-ttl=123",
						"--record-set-data=test.com",
					}
				})

				It("prints a useful message", func() {
					Eventually(session.Out, 5).Should(gbytes.Say(fmt.Sprintf("Created record set %s", rsName)))
				})

				It("creates the record set", func() {
					found := false

					// sleep for 3 seconds until creation is complete
					// TODO: this could be improved
					time.Sleep(3 * time.Second)

					rss, err := vinylClient.RecordSets(zone.Zone.ID)
					Expect(err).NotTo(HaveOccurred())

					for _, rs := range rss {
						if rs.Name == rsName {
							found = true
							break
						}
					}

					Expect(found).To(BeTrue())
				})
			})

			Context("when it's tasked in creating an MX record", func() {
				BeforeEach(func() {
					rsName = "some-mx"

					recordSetsArgs = []string{
						"record-set-create",
						fmt.Sprintf("--zone-name=%s", zName),
						fmt.Sprintf("--record-set-name=%s", rsName),
						"--record-set-type=mx",
						"--record-set-ttl=123",
						"--record-set-data=3,test.com",
					}
				})

				It("prints a useful message", func() {
					Eventually(session.Out, 5).Should(gbytes.Say(fmt.Sprintf("Created record set %s", rsName)))
				})

				It("creates the record set", func() {
					found := false

					// sleep for 3 seconds until creation is complete
					// TODO: this could be improved
					time.Sleep(3 * time.Second)

					rss, err := vinylClient.RecordSets(zone.Zone.ID)
					Expect(err).NotTo(HaveOccurred())

					for _, rs := range rss {
						if rs.Name == rsName {
							found = true
							break
						}
					}

					Expect(found).To(BeTrue())
				})
			})

			Context("when it's tasked in creating an TXT record", func() {
				BeforeEach(func() {
					rsName = "some-txt"

					recordSetsArgs = []string{
						"record-set-create",
						fmt.Sprintf("--zone-name=%s", zName),
						fmt.Sprintf("--record-set-name=%s", rsName),
						"--record-set-type=TXT",
						"--record-set-ttl=123",
						"--record-set-data=test TXT",
					}
				})

				It("prints a useful message", func() {
					Eventually(session.Out, 5).Should(gbytes.Say(fmt.Sprintf("Created record set %s", rsName)))
				})

				It("creates the record set", func() {
					found := false

					// sleep for 3 seconds until creation is complete
					// TODO: this could be improved
					time.Sleep(3 * time.Second)

					rss, err := vinylClient.RecordSets(zone.Zone.ID)
					Expect(err).NotTo(HaveOccurred())

					for _, rs := range rss {
						if rs.Name == rsName {
							found = true
							break
						}
					}

					Expect(found).To(BeTrue())
				})
			})
		})
	})
})
