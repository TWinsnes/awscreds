package console

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// SdkHelper describes an object that helps
type SdkHelper interface {
	GetStsClient(profile string) (stsiface.STSAPI, error)
}

// DefaultSdkHelper is the default sdk helper implementation
type DefaultSdkHelper struct{}

func (DefaultSdkHelper) GetStsClient(profile string) (stsiface.STSAPI, error) {

	var sessionOptions session.Options
	if profile != "" {
		sessionOptions = session.Options{
			Profile: profile,
		}
	} else {
		sessionOptions = session.Options{}
	}

	sess := session.Must(session.NewSessionWithOptions(sessionOptions))

	// if profile is not "" asume profile is defined and override credentials in case env vars have been set
	if profile != "" {
		sess.Config.Credentials = credentials.NewSharedCredentials("", profile)
	}

	stsClient := sts.New(sess)

	return stsClient, nil
}
