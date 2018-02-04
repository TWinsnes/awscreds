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
	"fmt"
	"os"

	"github.com/TWinsnes/awscreds/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var cfgFile string
var conf config.Handler

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awscreds",
	Short: "Helps make working with aws credentials from terminal easier.",
	Long:  `Command line utility to help work with aws credentials.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.awscreds)")

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		// flag not set, use default

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cfgFile = home + "/.awscreds"

	}

	conf = config.LoadConfig(cfgFile)

	fmt.Println(conf.SecretBackend())
}
