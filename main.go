package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: "vibrato-tom",
	}))

	stsClient := sts.New(sess)

	var duration int64 = 43200 // 12 hours
	name := "temptoken"
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

	url := "https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=" + sessionString

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	stringData := string(data[:])

	log.Print(stringData)
}
