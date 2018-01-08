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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

type signinToken struct {
	Token string `json:"SigninToken"`
}

type awsCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
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
var sessionDuration string
var printKeys bool

func init() {

	rootCmd.AddCommand(consoleCmd)

	consoleCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS CLI profile name")
	consoleCmd.Flags().StringVarP(&service, "service", "s", "", "AWS Service to connect to")
	consoleCmd.Flags().StringVarP(&sessionDuration, "session-duration", "t", "12h", "Length of session duration (suffix with s/m/h)")
	consoleCmd.Flags().BoolVar(&printKeys, "printkeys", false, "Set this to print federated keys to console")

}

func parseSessionDuration(sessionDuration string) (sessionSeconds int64) {
	// Try to parse duration string as-is
	sessionSeconds, err := strconv.ParseInt(sessionDuration, 10, 64)
	if err != nil {
		// If duration string fails to parse, assume there is a time suffix
		durationPrefix, err := strconv.ParseInt(sessionDuration[0:len(sessionDuration)-1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		durationSuffix := sessionDuration[len(sessionDuration)-1:]

		switch durationSuffix {
		case "h":
			sessionSeconds = int64(durationPrefix * 60 * 60)
		case "m":
			sessionSeconds = int64(durationPrefix * 60)
		case "s":
			sessionSeconds = int64(durationPrefix)
		default:
			log.Fatalf("Session duration suffix \"%s\" is not valid", durationSuffix)
		}
	}
	return sessionSeconds
}

func openConsole(name string, profile string, service string) error {

	if service == "" {
		service = "console"
	}

	envSessionToken := os.Getenv("AWS_SESSION_TOKEN")

	var credentials awsCredentials
	if envSessionToken == "" {

		var sessionOptions session.Options
		if profile == "" {
			sessionOptions = session.Options{}
		} else {
			sessionOptions = session.Options{
				Profile: profile,
			}
		}
		sess := session.Must(session.NewSessionWithOptions(sessionOptions))

		duration := parseSessionDuration(sessionDuration)
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

		stsClient := sts.New(sess)
		input := sts.GetFederationTokenInput{Name: &name, DurationSeconds: &duration, Policy: &policy}
		tokenResponse, err := stsClient.GetFederationToken(&input)
		if err != nil {
			log.Fatal(err)
		}
		credentials = awsCredentials{
			AccessKeyID:     *tokenResponse.Credentials.AccessKeyId,
			SecretAccessKey: *tokenResponse.Credentials.SecretAccessKey,
			SessionToken:    *tokenResponse.Credentials.SessionToken,
		}
	} else {
		credentials = awsCredentials{
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			SessionToken:    envSessionToken,
		}
	}
	if credentials.AccessKeyID == "" || credentials.SecretAccessKey == "" {
		log.Fatal(
			"\"AWS_ACCESS_KEY_ID\" and \"AWS_SECRET_ACCESS_KEY\" environment " +
				"variables must be set when using \"AWS_SESSION_TOKEN\"",
		)
	}

	sessionString := "{" +
		"\"sessionId\":\"" + credentials.AccessKeyID + "\"," +
		"\"sessionKey\":\"" + credentials.SecretAccessKey + "\"," +
		"\"sessionToken\":\"" + credentials.SessionToken + "\"" +
		"}"

	if printKeys {
		fmt.Printf("Session ID:    %s \n", credentials.AccessKeyID)
		fmt.Printf("Session Key:   %s \n", credentials.SecretAccessKey)
		fmt.Printf("Session Token: %s \n", credentials.SessionToken)
	}

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
