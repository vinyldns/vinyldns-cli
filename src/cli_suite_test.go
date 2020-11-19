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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

var (
	exe         string
	baseArgs    []string
	vinylClient *vinyldns.Client
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
		panic(err)
	}
	for _, g := range gs {
		vinylClient.GroupDelete(g.ID)
	}

	// ensure there are no pre-existing zones
	zs, err := vinylClient.Zones()
	if err != nil {
		panic(err)
	}
	for _, z := range zs {
		vinylClient.ZoneDelete(z.ID)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "vinyldns CLI integration test suite")
}
