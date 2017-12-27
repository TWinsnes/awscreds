package console

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type signinToken struct {
	Token string `json:"SigninToken"`
}

// OpenConsole Opens console window logged into defined profile
func OpenConsole(name string, profile string) error {

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

	log.Print(t.Token)

	loginURL, err := url.Parse("https://signin.aws.amazon.com/federation")

	if err != nil {
		log.Fatal(err)
	}

	parameters := url.Values{}
	parameters.Add("Action", "login")
	parameters.Add("Destination", "https://console.aws.amazon.com/console/home")
	parameters.Add("SigninToken", t.Token)
	loginURL.RawQuery = parameters.Encode()

	// /federation?Action=login&Destination=https://console.aws.amazon.com/console/home&SigninToken=" + t.Token

	log.Print(loginURL.String())

	err = exec.Command("open", loginURL.String()).Start()

	if err != nil {
		panic(err)
	}

	return nil
}
