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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TWinsnes/awscreds/config"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Writes defaults to config file",
	Long:  `Writes all the default configuration values to the .awsconfig file in your home directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		err := getUserConfigValues(conf)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = conf.SaveConfig(cfgFile)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Wrote default configuration values to file")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func getUserConfigValues(conf config.Handler) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Secret Backend [%s]", conf.SecretBackend())
	secretBackend, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	// remove trailing
	secretBackend = strings.TrimSpace(secretBackend)
	if secretBackend != "" {
		err := conf.SetSecretBackend(secretBackend)

		if err != nil {
			return err
		}
	}

	return nil
}
