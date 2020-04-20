// Copyright © 2017 Alex Beltran <alex.e.beltran@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var interfacesCmd = &cobra.Command{
	Use:   "interfaces",
	Short: "lists all available network interfaces",
	Long:  `lists all available network interfaces`,

	Run: interfaces,
}

func interfaces(cmd *cobra.Command, args []string) {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{}
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)

    ifaces, err := net.Interfaces()
    if err != nil {
        log.Fatalf("localAddresses:", err)
        return
    }
    for _, i := range ifaces {
    	log.Infof(i.Name)
        addrs, err := i.Addrs()
        if err != nil {
	        log.Fatalf("localAddresses:", err)
            continue
        }
        for _, a := range addrs {
            switch v := a.(type) {
            case *net.IPAddr:
                log.Infof("%v : %s (%s)\n", i.Name, v, v.IP.DefaultMask())
            }

        }
    }
}

func init() {
	RootCmd.AddCommand(interfacesCmd)
}
