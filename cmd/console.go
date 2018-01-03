// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

type signinToken struct {
	Token string `json:"SigninToken"`
}

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Logs into and opens console in default browser using aws cli profile",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		openConsole("awscreds", profile, service)
	},
}

var profile string
var service string

func init() {

	rootCmd.AddCommand(consoleCmd)

	consoleCmd.Flags().StringVarP(&profile, "profile", "p", "Default", "AWS CLI profile name")
	consoleCmd.Flags().StringVarP(&service, "service", "s", "", "AWS Service to connect to")

}

func openConsole(name string, profile string, service string) error {

	if service == "" {
		service = "console"
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profile,
	}))

	stsClient := sts.New(sess)

	var duration int64 = 43200 // 12 hours
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": "*",
				"Resource": "*"
			}
		]
	}`

	input := sts.GetFederationTokenInput{Name: &name, DurationSeconds: &duration, Policy: &policy}

	token, err := stsClient.GetFederationToken(&input)

	if err != nil {
		log.Fatal(err)
	}

	sessionString := "{" +
		"\"sessionId\":\"" + *token.Credentials.AccessKeyId + "\"," +
		"\"sessionKey\":\"" + *token.Credentials.SecretAccessKey + "\"," +
		"\"sessionToken\":\"" + *token.Credentials.SessionToken + "\"" +
		"}"

	federationURL, err := url.Parse("https://signin.aws.amazon.com/federation")

	if err != nil {
		panic(err)
	}

	federationParams := url.Values{}
	federationParams.Add("Action", "getSigninToken")
	federationParams.Add("Session", sessionString)
	federationURL.RawQuery = federationParams.Encode()

	resp, err := http.Get(federationURL.String())

	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var t signinToken

	err = json.Unmarshal(data, &t)

	loginURL, err := url.Parse("https://signin.aws.amazon.com/federation")

	if err != nil {
		log.Fatal(err)
	}

	parameters := url.Values{}
	parameters.Add("Action", "login")
	parameters.Add("Destination", "https://console.aws.amazon.com/"+service+"/home")
	parameters.Add("SigninToken", t.Token)
	loginURL.RawQuery = parameters.Encode()

	err = exec.Command("open", loginURL.String()).Start()

	if err != nil {
		panic(err)
	}

	return nil
}
