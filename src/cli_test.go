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
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("the --help flag", func() {
	var (
		session *gexec.Session
		err     error
	)

	BeforeEach(func() {
		args := []string{
			"--help",
		}

		cmd := exec.Command(exe, args...)
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	})

	AfterEach(func() {
		if session != nil {
			session.Terminate()
		}
	})

	It("does not error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	It("prints NAME info", func() {
		Eventually(session.Out, 5).Should(gbytes.Say("vinyldns - A CLI to the vinyldns DNS-as-a-service API"))
	})

	It("prints USAGE info", func() {
		Eventually(session.Out, 5).Should(gbytes.Say("USAGE"))
	})

	It("prints COMMANDS info", func() {
		Eventually(session.Out, 5).Should(gbytes.Say("COMMANDS"))
	})

	It("prints GLOBAL OPTIONS info", func() {
		Eventually(session.Out, 5).Should(gbytes.Say("GLOBAL OPTIONS"))
	})
})
