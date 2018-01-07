// Copyright © 2018 Thomas Winsnes <twinsnes@live.com>
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
	"github.com/TWinsnes/awscreds/cmd/console"
	"github.com/spf13/cobra"
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Logs into and opens console in default browser using aws cli profile",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		browser := console.DefaultBrowser{}
		consoleOptions.OpenConsole("awscreds", browser)
	},
}

var consoleOptions = &console.Console{}

func init() {

	rootCmd.AddCommand(consoleCmd)

	consoleCmd.Flags().StringVarP(&consoleOptions.Profile, "profile", "p", "Default", "AWS CLI profile name")
	consoleCmd.Flags().StringVarP(&consoleOptions.Service, "service", "s", "", "AWS Service to connect to")
	consoleCmd.Flags().BoolVar(&consoleOptions.PrintKeys, "printkeys", false, "Set this to print federated keys to console")

}
