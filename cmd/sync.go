// Copyright © 2019 VinkDong <dong@wenqi.us>
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
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"github.com/vinkdong/gox/log"
	"github.com/vinkdong/image-sync/pkg/docker"
	"gopkg.in/yaml.v2"
	"fmt"
	"time"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync two docker registry",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			config string
			err    error
		)

		if config, err = cmd.Flags().GetString("config"); err != nil {
			log.Error(err)
			os.Exit(128)
		}
		if config == "" {
			log.Error("must specify a config file")
		}

		data, err := ioutil.ReadFile(config)
		if err != nil {
			log.Error(err)
			os.Exit(128)
		}
		sync := &docker.NamedSync{}
		yaml.Unmarshal(data, sync)
		if err := sync.Sync.Do(); err != nil {
			log.Error(err)
			os.Exit(128)
		}

		daemon, err := cmd.Flags().GetBool("daemon")
		if err != nil {
			log.Error(err)
			os.Exit(128)
		}
		if daemon {
			for {
				select {
				case <-time.Tick(time.Minute * 3):
					if err := sync.Sync.Do(); err != nil {
						log.Error(err)
						os.Exit(128)
					}
				}
			}
		}
		fmt.Println("sync succeed")
	},
}

func init() {
	syncCmd.Flags().StringP("config", "c", "", "specify a config file")
	syncCmd.Flags().BoolP("daemon", "d", false, "run as daemon")
	rootCmd.AddCommand(syncCmd)
}
