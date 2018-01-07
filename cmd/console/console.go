//
package console

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type signinToken struct {
	Token string `json:"SigninToken"`
}

// Browser describes an object that interacts with a browser
type Browser interface {
	Open(url string) error
}

// DefaultBrowser represents system default browser
type DefaultBrowser struct{}

// Console wrapper for the console command
type Console struct {
	Profile   string
	Service   string
	PrintKeys bool
}

// OpenConsole opens the console using
func (c *Console) OpenConsole(name string, browser Browser) error {

	if c.Service == "" {
		c.Service = "console"
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: c.Profile,
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

	if c.PrintKeys {
		fmt.Printf("Session ID: %s \n", *token.Credentials.AccessKeyId)
		fmt.Printf("Session Key: %s \n", *token.Credentials.SecretAccessKey)
		fmt.Printf("Session Token: %s \n", *token.Credentials.SessionToken)
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

	loginURL, err := c.getLoginURL(c.Service, t.Token)

	if err != nil {
		return err
	}

	err = browser.Open(loginURL)

	return err
}

func (c *Console) getLoginURL(service string, token string) (string, error) {
	urlStruct, err := url.Parse("https://signin.aws.amazon.com/federation")

	if err != nil {
		return "", err
	}

	parameters := url.Values{}
	parameters.Add("Action", "login")
	parameters.Add("Destination", "https://console.aws.amazon.com/"+service+"/home")
	parameters.Add("SigninToken", token)
	urlStruct.RawQuery = parameters.Encode()

	loginURL := urlStruct.String()
	return loginURL, err
}

// Open Opens url in default browser
func (DefaultBrowser) Open(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
