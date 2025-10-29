/*
Copyright 2018 Comcast Cable Communications Management, LLC
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

	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func zoneByName(c *vinyldns.Client, name string) (vinyldns.Zone, error) {
	var z vinyldns.Zone
	zone, err := c.ZoneByName(name)
	if err != nil {
		return z, err
	}

	return zone, nil
}

func getZoneID(c *vinyldns.Client, id, name string) (string, error) {
	if id != "" {
		return id, nil
	}

	z, err := zoneByName(c, name)
	if err != nil {
		return "", err
	}

	return z.ID, nil
}

func getZone(c *vinyldns.Client, name, id string) (vinyldns.Zone, error) {
	if name != "" {
		return zoneByName(c, name)
	}

	return c.Zone(id)
}

func getZoneDetails(c *vinyldns.Client, id string) (vinyldns.ZoneDetails, error) {

	return c.ZoneDetails(id)
}

func validateConnection(connection string, c *vinyldns.ZoneConnection) (bool, error) {
	// if all are empty, we assume the user does not want to declare a connection
	if c.Key == "" && c.KeyName == "" && c.Name == "" && c.PrimaryServer == "" {
		return false, nil
	}

	// if any but not all are empty, we have a problem
	if c.Key == "" || c.KeyName == "" || c.Name == "" || c.PrimaryServer == "" {
		return false, fmt.Errorf("%s connection requires '--%s-connection-key-name', '--%s-connection-key', and '--%s-connection-primary-server'", connection, connection, connection, connection)
	}

	return true, nil
}
