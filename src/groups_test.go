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
	"io/ioutil"
	"os/exec"

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
	})

	JustAfterEach(func() {
		if session != nil {
			session.Terminate()
		}
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

				groupsArgs = []string{
					"groups",
				}
			})

			AfterEach(func() {
				_, err = vinylClient.GroupDelete(group.ID)
				Expect(err).NotTo(HaveOccurred())
			})

			It("does not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the groups", func() {
				Eventually(session.Out, 5).Should(gbytes.Say(name))
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
			var (
				name string = "ok-group"
			)

			BeforeEach(func() {
				j, err := ioutil.ReadFile("../fixtures/group_create_json")
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
				Eventually(session.Out, 5).Should(gbytes.Say("Created group ok-group"))
			})
		})
	})
})
