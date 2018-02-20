// Package console ...
package console

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

var policy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": "*",
			"Resource": "*"
		}
	]
}`

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
	Profile         string
	Service         string
	SessionDuration string
	PrintKeys       bool
	PrintURL        bool
}

// AwsCredentials acts as a credential storage structure across providers
type AwsCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// OpenConsole opens the console using
func (c *Console) OpenConsole(browser Browser, sdkHelper SdkHelper) error {

	var creds AwsCredentials
	var err error

	if c.Service == "" {
		c.Service = "console"
	}

	stsClient, err := sdkHelper.GetStsClient(c.Profile)

	if err != nil {
		return err
	}

	duration, err := c.parseSessionDuration(c.SessionDuration)

	if err != nil {
		return err
	}

	userName, err := c.getAwsUsername(stsClient)

	if err != nil {
		return err
	}

	// override the default sdk behaviour of using env vars over anything else
	// this follows the convention of the cli
	if c.Profile != "" {
		creds, err = c.getFederationToken(stsClient, userName, duration)

		if err != nil {
			return err
		}
	} else {

		// If profile has not been set explicitly we fall back to default SDK
		// behaviour. Since we can't use an sts token to get federated
		// credentials, we try to load those from environment variables.
		// If that fails, we leave it to the SDK to figure out credentials.
		envCredentials, envCredErr := c.getCredentialsFromEnvironment()

		if envCredErr == nil {
			creds = envCredentials
		} else {
			creds, err = c.getFederationToken(stsClient, userName, duration)

			if err != nil {
				return err
			}
		}
	}

	sessionString := "{" +
		"\"sessionId\":\"" + creds.AccessKeyID + "\"," +
		"\"sessionKey\":\"" + creds.SecretAccessKey + "\"," +
		"\"sessionToken\":\"" + creds.SessionToken + "\"" +
		"}"

	if c.PrintKeys {
		fmt.Printf("Session ID:    %s \n", creds.AccessKeyID)
		fmt.Printf("Session Key:   %s \n", creds.SecretAccessKey)
		fmt.Printf("Session Token: %s \n", creds.SessionToken)
	}

	federationURL, err := url.Parse("https://signin.aws.amazon.com/federation")

	if err != nil {
		return err
	}

	federationParams := url.Values{}
	federationParams.Add("Action", "getSigninToken")
	federationParams.Add("Session", sessionString)
	federationURL.RawQuery = federationParams.Encode()

	resp, err := http.Get(federationURL.String())

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()
	if err != nil {
		return err
	}

	var t signinToken

	err = json.Unmarshal(data, &t)

	loginURL, err := c.getLoginURL(c.Service, t.Token)

	if err != nil {
		return err
	}

	if c.PrintURL {
		fmt.Print(loginURL)
	} else {
		err = browser.Open(loginURL)
	}

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

func (c *Console) getAwsUsername(stsClient stsiface.STSAPI) (string, error) {

	var callerIdentityInput *sts.GetCallerIdentityInput

	output, err := stsClient.GetCallerIdentity(callerIdentityInput)

	if err != nil {
		return "", err
	}

	callerIdentity := *output.Arn

	if err != nil {
		return "", err
	}

	splitArn := strings.Split(callerIdentity, "/")
	username := splitArn[len(splitArn)-1]

	if len(username) > 32 {
		username = username[:32]
	}

	return username, nil
}

func (c *Console) getCredentialsFromEnvironment() (AwsCredentials, error) {
	credentials := AwsCredentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
	}

	if credentials.AccessKeyID == "" ||
		credentials.SecretAccessKey == "" ||
		credentials.SessionToken == "" {
		err := fmt.Errorf("\"AWS_ACCESS_KEY_ID\", \"AWS_SECRET_ACCESS_KEY\", " +
			"and \"AWS_SESSION_TOKEN\" environment variables must be set when" +
			" using environment variables for authentication.")

		return credentials, err
	}

	return credentials, nil
}

func (c *Console) parseSessionDuration(duration string) (int64, error) {
	// Try to parse duration string as-is
	sessionSeconds, err := strconv.ParseInt(duration, 10, 64)
	if err != nil {
		// If duration string fails to parse, assume there is a time suffix
		durationPrefix, err := strconv.ParseInt(duration[0:len(duration)-1], 10, 64)
		if err != nil {
			return 0, err
		}
		durationSuffix := duration[len(duration)-1:]

		switch durationSuffix {
		case "h":
			sessionSeconds = int64(durationPrefix * 60 * 60)
		case "m":
			sessionSeconds = int64(durationPrefix * 60)
		case "s":
			sessionSeconds = int64(durationPrefix)
		default:
			return 0, fmt.Errorf("Session duration suffix \"%s\" is not valid", durationSuffix)
		}
	}
	return sessionSeconds, nil
}

func (c *Console) getFederationToken(stsClient stsiface.STSAPI, name string, duration int64) (AwsCredentials, error) {

	input := sts.GetFederationTokenInput{Name: &name, DurationSeconds: &duration, Policy: &policy}

	token, err := stsClient.GetFederationToken(&input)

	credentials := AwsCredentials{
		AccessKeyID:     *token.Credentials.AccessKeyId,
		SecretAccessKey: *token.Credentials.SecretAccessKey,
		SessionToken:    *token.Credentials.SessionToken,
	}

	return credentials, err
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
