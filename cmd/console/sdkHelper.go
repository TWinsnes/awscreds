package console

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// SdkHelper describes an object that helps
type SdkHelper interface {
	GetFederationToken(profile string, name string, duration int64) (AwsCredentials, error)
	GetCallerIdentity(sts *sts.STS) (string, error)
}

// DefaultSdkHelper is the default sdk helper implementation
type DefaultSdkHelper struct{}

// GetFederationToken s
func (DefaultSdkHelper) GetFederationToken(profile string, name string, duration int64) (AwsCredentials, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profile,
	}))

	sess.Config.Credentials = credentials.NewSharedCredentials("", profile)

	stsClient := sts.New(sess)

	localPolicy := policy // can't reference a const, wil have to copy this

	input := sts.GetFederationTokenInput{Name: &name, DurationSeconds: &duration, Policy: &localPolicy}

	token, err := stsClient.GetFederationToken(&input)

	credentials := AwsCredentials{
		AccessKeyID:     *token.Credentials.AccessKeyId,
		SecretAccessKey: *token.Credentials.SecretAccessKey,
		SessionToken:    *token.Credentials.SessionToken,
	}

	return credentials, err
}

func (DefaultSdkHelper) GetCallerIdentity(stsClient *sts.STS) (string, error) {
	var callerIdentityInput *sts.GetCallerIdentityInput

	output, err := stsClient.GetCallerIdentity(callerIdentityInput)

	if err != nil {
		return "", err
	}

	return *output.Arn, err
}
