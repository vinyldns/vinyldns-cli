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

var _ = Describe("its commands for working with groups", func() {
	var (
		session    *gexec.Session
		err        error
		args       []string
		groupsArgs []string
	)

	JustBeforeEach(func() {
		args = append(baseArgs, groupsArgs...)
		cmd := exec.Command(exe, args...)
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	JustAfterEach(func() {
		session.Terminate()
	})

	Describe("its 'groups' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				groupsArgs = []string{
					"groups",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("List all vinyldns groups"))
			})
		})

		Context("when no groups exist", func() {
			Context("when not passed an --output", func() {
				BeforeEach(func() {
					groupsArgs = []string{
						"groups",
					}
				})

				It("prints the correct data", func() {
					Eventually(session.Out, 5).Should(gbytes.Say("No groups found"))
				})
			})

			Context("when passed --output=json", func() {
				BeforeEach(func() {
					groupsArgs = []string{
						"--output=json",
						"groups",
					}
				})

				It("prints the correct data", func() {
					Eventually(session.Out, 5).Should(gbytes.Say(`\[\]`))
				})
			})
		})

		Context("when groups exist", func() {
			var (
				group *vinyldns.Group
				name  string = "test-group"
			)

			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup(name))
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				err = deleteAllGroups()
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when it's not passed the --output=json option", func() {
				BeforeEach(func() {
					groupsArgs = []string{
						"groups",
					}
				})

				It("prints groups details", func() {
					output := fmt.Sprintf(`+------------+--------------------------------------+
|    NAME    |                  ID                  |
+------------+--------------------------------------+
| %s | %s |
+------------+--------------------------------------+`, group.Name, group.ID)

					Eventually(func() string {
						return string(session.Out.Contents())
					}).Should(ContainSubstring(output))
				})
			})

			Context("when it's passed the --output=json option", func() {
				BeforeEach(func() {
					groupsArgs = []string{
						"--output=json",
						"groups",
					}
				})

				It("prints the groups JSON", func() {
					Eventually(func() string {
						return string(session.Out.Contents())
					}).Should(ContainSubstring(fmt.Sprintf(`"name":"%s","email":"email@email.com","description":"description","status":"Active"`, name)))
				})
			})
		})
	})

	Describe("its 'group-create' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				groupsArgs = []string{
					"group-create",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Create a vinyldns group"))
			})
		})

		Context("when it's passed group JSON", func() {
			Context("when it's not passed '--output=json'", func() {
				var (
					name string = "group-test-create"
				)

				BeforeEach(func() {
					g := makeGroup(name)
					j, err := json.Marshal(g)
					Expect(err).NotTo(HaveOccurred())

					groupsArgs = []string{
						"group-create",
						"--json",
						string(j),
					}
				})

				AfterEach(func() {
					err = deleteAllGroups()
					Expect(err).NotTo(HaveOccurred())
				})

				It("creates the group and prints a helpful message", func() {
					Eventually(session.Out, 5).Should(gbytes.Say(fmt.Sprintf("Created group %s", name)))
				})
			})
		})
	})

	Describe("its 'group' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				groupsArgs = []string{
					"group",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Retrieve details for vinyldns group"))
			})
		})

		Context("when it's passed the name of a group that exists", func() {
			var (
				group *vinyldns.Group
				name  string = "group-test"
			)

			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(&vinyldns.Group{
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
				})

			})

			AfterEach(func() {
				err = deleteAllGroups()
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when it's not passed --output=json", func() {
				BeforeEach(func() {
					groupsArgs = []string{
						"group",
						fmt.Sprintf("--name=%s", name),
					}
				})

				It("prints the group details", func() {
					output := fmt.Sprintf(`+-------------+--------------------------------------+
| Name        | %s                           |
+-------------+--------------------------------------+
| ID          | %s |
+-------------+--------------------------------------+
| Email       | email@email.com                      |
+-------------+--------------------------------------+
| Description | description                          |
+-------------+--------------------------------------+
| Status      | Active                               |
+-------------+--------------------------------------+
| Members     | ok                                   |
+-------------+--------------------------------------+
| Admins      | ok                                   |
+-------------+--------------------------------------+`, group.Name, group.ID)

					Eventually(func() string {
						return string(session.Out.Contents())
					}).Should(ContainSubstring(output))
				})
			})

			Context("when it's passed --output=json", func() {
				BeforeEach(func() {
					groupsArgs = []string{
						"--output=json",
						"group",
						fmt.Sprintf("--name=%s", name),
					}
				})

				It("prints the group details JSON", func() {
					Eventually(func() string {
						return string(session.Out.Contents())
					}).Should(ContainSubstring(fmt.Sprintf(`"name":"%s","email":"email@email.com","description":"description","status":"Active","created":`, name)))
				})
			})
		})
	})

	Describe("its 'group-update' command", func() {
		Context("when it's passed '--help'", func() {
			BeforeEach(func() {
				groupsArgs = []string{
					"group-update",
					"--help",
				}
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Update a vinyldns group"))
			})
		})

		Context("when it's passed the JSON for a valid, existing group", func() {
			var (
				group       *vinyldns.Group
				err         error
				updatedDesc string = "updated-description"
				name        string = "group-update-test"
			)

			BeforeEach(func() {
				g := makeGroup(name)
				group, err = vinylClient.GroupCreate(g)
				Expect(err).NotTo(HaveOccurred())

				groupUpdated := group
				groupUpdated.Description = updatedDesc
				j, err := json.Marshal(groupUpdated)
				Expect(err).NotTo(HaveOccurred())

				fmt.Println(string(j))

				groupsArgs = []string{
					"group-update",
					"--json",
					string(j),
				}
			})

			AfterEach(func() {
				err = deleteAllGroups()
				Expect(err).NotTo(HaveOccurred())
			})

			It("prints a message", func() {
				Eventually(session.Out, 5).Should(gbytes.Say(fmt.Sprintf("Updated group %s", name)))
			})

			It("updates the group", func() {
				time.Sleep(3 * time.Second)

				g, err := vinylClient.Group(group.ID)
				fmt.Println(g)

				Expect(err).NotTo(HaveOccurred())
				Expect(g.Description).To(Equal(updatedDesc))
			})
		})
	})
})
