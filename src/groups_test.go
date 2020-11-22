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
		makeGroup  = func(name string) *vinyldns.Group {
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
	)

	JustBeforeEach(func() {
		args = append(baseArgs, groupsArgs...)
		cmd := exec.Command(exe, args...)
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
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

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
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

				It("does not error", func() {
					Expect(err).NotTo(HaveOccurred())
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

				It("does not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("prints the correct data", func() {
					Eventually(session.Out, 5).Should(gbytes.Say(`\[\]`))
				})
			})
		})

		Context("when groups exist", func() {
			var (
				group *vinyldns.Group
				name  string = "ok-groups-test"
			)

			BeforeEach(func() {
				group, err = vinylClient.GroupCreate(makeGroup(name))
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				_, err = vinylClient.GroupDelete(group.ID)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when it's not passed the --output=json option", func() {
				BeforeEach(func() {
					groupsArgs = []string{
						"groups",
					}
				})

				It("does not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("prints groups details", func() {
					output := fmt.Sprintf(`+----------------+--------------------------------------+
|      NAME      |                  ID                  |
+----------------+--------------------------------------+
| ok-groups-test | %s |
+----------------+--------------------------------------+`, group.ID)

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

				It("does not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("prints the groups JSON", func() {
					Eventually(func() string {
						return string(session.Out.Contents())
					}).Should(ContainSubstring(`"name":"ok-groups-test","email":"email@email.com","description":"description","status":"Active"`))
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

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Create a vinyldns group"))
			})
		})

		Context("when it's passed group JSON", func() {
			Context("when it's not passed '--output=json'", func() {
				var (
					name string = "ok-group-create-test"
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
					groups, err := vinylClient.Groups()

					for _, g := range groups {
						if g.Name == name {
							_, err = vinylClient.GroupDelete(g.ID)
							Expect(err).NotTo(HaveOccurred())
						}
					}
				})

				It("does not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("creates the group and prints a helpful message", func() {
					Eventually(session.Out, 5).Should(gbytes.Say("Created group ok-group-create-test"))
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

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("prints a useful description", func() {
				Eventually(session.Out, 5).Should(gbytes.Say("Retrieve details for vinyldns group"))
			})
		})

		Context("when it's passed the name of a group that exists", func() {
			var (
				group *vinyldns.Group
				name  string = "ok-group-test"
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
				_, err = vinylClient.GroupDelete(group.ID)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when it's not passed --output=json", func() {
				BeforeEach(func() {
					groupsArgs = []string{
						"group",
						fmt.Sprintf("--name=%s", name),
					}
				})

				It("does not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("prints the group details", func() {
					output := fmt.Sprintf(`+-------------+--------------------------------------+
| Name        | ok-group-test                        |
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
+-------------+--------------------------------------+`, group.ID)

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

				It("does not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("prints the group details JSON", func() {
					Eventually(func() string {
						return string(session.Out.Contents())
					}).Should(ContainSubstring(`"name":"ok-group-test","email":"email@email.com","description":"description","status":"Active","created":`))
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

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
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
				name        string = "ok-group-update-test"
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
				_, err := vinylClient.GroupDelete(group.ID)
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
