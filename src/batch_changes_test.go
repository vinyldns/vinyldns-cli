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

var _ = Describe("its commands for working with batch changes", func() {
	var (
		session *gexec.Session
		err     error
		args    []string
		bcArgs  []string
	)

	JustBeforeEach(func() {
		args = append(baseArgs, bcArgs...)
		cmd := exec.Command(exe, args...)
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	JustAfterEach(func() {
		session.Terminate()
	})

	Describe("its 'batch-changes' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				bcArgs = []string{
					"batch-changes",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("List all batch changes"))
			})
		})
	})

	Describe("its 'batch-change-create' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				bcArgs = []string{
					"batch-change-create",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Create a batch change"))
			})
		})

		Context("when it's passed JSON", func() {
			var (
				zone   *vinyldns.ZoneUpdateResponse
				group  *vinyldns.Group
				zName  string = "vinyldns."
				rsName string = fmt.Sprintf("batch-change.%s", zName)
			)

			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup("test-group"))
				Expect(err).NotTo(HaveOccurred())

				zone, err = vinylClient.ZoneCreate(makeZone(zName, group.ID))
				Expect(err).NotTo(HaveOccurred())

				// poll until the new zone exists
				for {
					exists, err := vinylClient.ZoneExists(zone.Zone.ID)
					Expect(err).NotTo(HaveOccurred())

					if exists {
						break
					}
				}

				jsonData := fmt.Sprintf(`{
					"comments": "request on behalf of someone",
					"changes": [{
						"inputName": "%s",
						"changeType": "Add",
						"type": "A",
						"ttl": 7200,
						"record": {
							"address": "1.1.1.1"
						}
					}]
				}`, rsName)

				bcArgs = []string{
					"batch-change-create",
					fmt.Sprintf("--json=%s", jsonData),
				}
			})

			AfterEach(func() {
				// sleep until the record set is completely created
				// TODO: this can be improved
				time.Sleep(5 * time.Second)

				err = deleteRecordInZone(zone.Zone.ID, rsName)
				Expect(err).NotTo(HaveOccurred())

				err = deleteAllGroupsAndZones()
				Expect(err).NotTo(HaveOccurred())
			})

			It("prints the correct message", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("request on behalf of someone"))
			})
		})
	})
})
