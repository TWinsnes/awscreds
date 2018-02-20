package console

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/stretchr/testify/assert"
)

func TestGetLoginUrl(t *testing.T) {

	t.Parallel()
	expected := "https://signin.aws.amazon.com/federation?Action=login&Destination=https%3A%2F%2Fconsole.aws.amazon.com%2Fs3%2Fhome&SigninToken=xxx"

	c := Console{}

	output, err := c.getLoginURL("s3", "xxx")

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expected, output)

}

func TestParseSessionDurationHours(t *testing.T) {

	t.Parallel()
	var expected int64 = 10800

	c := Console{}

	output, err := c.parseSessionDuration("3h")

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expected, output)
}

func TestParseSessionDurationMinutes(t *testing.T) {

	t.Parallel()
	var expected int64 = 900

	c := Console{}

	output, err := c.parseSessionDuration("15m")

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expected, output)
}

func TestParseSessionDurationSeconds(t *testing.T) {

	t.Parallel()
	var expected int64 = 1000

	c := Console{}

	output, err := c.parseSessionDuration("1000s")

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expected, output)
}

func TestParseSessionDurationDefault(t *testing.T) {

	t.Parallel()
	var expected int64 = 1200

	c := Console{}

	output, err := c.parseSessionDuration("1200")

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expected, output)
}

func TestParseSessionDurationBadSuffix(t *testing.T) {

	t.Parallel()
	var expected int64

	c := Console{
		SessionDuration: "1200x",
	}

	output, err := c.parseSessionDuration(c.SessionDuration)

	assert.Error(t, err, "Expected an error")
	assert.Equal(t, expected, output)
}

func TestGetAwsUsername(t *testing.T) {

	stsClient := mockStsApi{}

	stsClient.CallerIdentity = "somethingrandom/UserName"

	c := Console{}

	name, err := c.getAwsUsername(stsClient)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "UserName", name)
}

func TestGetAwsUsernameTooLong(t *testing.T) {
	stsClient := mockStsApi{}

	stsClient.CallerIdentity = "somethingrandom/somethinglongerthan32charactersmaybe"

	c := Console{}

	name, err := c.getAwsUsername(stsClient)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "somethinglongerthan32charactersm", name)
}

type mockSdkHelper struct{}

type mockStsApi struct {
	CallerIdentity string
}

func (mockSdkHelper) GetCallerIdentity(stsClient stsiface.STSAPI) (string, error) {
	return "sdfsd/UserName", nil
}

func (mockSdkHelper) GetFederationToken(profile string, name string, duration int64) (AwsCredentials, error) {
	var creds AwsCredentials
	return creds, nil
}

func (mockSdkHelper) GetStsClient(profile string) (stsiface.STSAPI, error) {
	client := mockStsApi{}
	return client, nil
}

type mockBrowser struct{}

func (mockBrowser) Open(url string) error {
	return nil
}

func (m mockStsApi) AssumeRole(*sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	return nil, nil
}
func (m mockStsApi) AssumeRoleWithContext(aws.Context, *sts.AssumeRoleInput, ...request.Option) (*sts.AssumeRoleOutput, error) {
	return nil, nil
}
func (m mockStsApi) AssumeRoleRequest(*sts.AssumeRoleInput) (*request.Request, *sts.AssumeRoleOutput) {
	return nil, nil
}

func (m mockStsApi) AssumeRoleWithSAML(*sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error) {
	return nil, nil
}
func (m mockStsApi) AssumeRoleWithSAMLWithContext(aws.Context, *sts.AssumeRoleWithSAMLInput, ...request.Option) (*sts.AssumeRoleWithSAMLOutput, error) {
	return nil, nil
}
func (m mockStsApi) AssumeRoleWithSAMLRequest(*sts.AssumeRoleWithSAMLInput) (*request.Request, *sts.AssumeRoleWithSAMLOutput) {
	return nil, nil
}

func (m mockStsApi) AssumeRoleWithWebIdentity(*sts.AssumeRoleWithWebIdentityInput) (*sts.AssumeRoleWithWebIdentityOutput, error) {
	return nil, nil
}
func (m mockStsApi) AssumeRoleWithWebIdentityWithContext(aws.Context, *sts.AssumeRoleWithWebIdentityInput, ...request.Option) (*sts.AssumeRoleWithWebIdentityOutput, error) {
	return nil, nil
}
func (m mockStsApi) AssumeRoleWithWebIdentityRequest(*sts.AssumeRoleWithWebIdentityInput) (*request.Request, *sts.AssumeRoleWithWebIdentityOutput) {
	return nil, nil
}

func (m mockStsApi) DecodeAuthorizationMessage(*sts.DecodeAuthorizationMessageInput) (*sts.DecodeAuthorizationMessageOutput, error) {
	return nil, nil
}
func (m mockStsApi) DecodeAuthorizationMessageWithContext(aws.Context, *sts.DecodeAuthorizationMessageInput, ...request.Option) (*sts.DecodeAuthorizationMessageOutput, error) {
	return nil, nil
}
func (m mockStsApi) DecodeAuthorizationMessageRequest(*sts.DecodeAuthorizationMessageInput) (*request.Request, *sts.DecodeAuthorizationMessageOutput) {
	return nil, nil
}

func (m mockStsApi) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	output := sts.GetCallerIdentityOutput{}
	value := m.CallerIdentity
	output.Arn = &value
	return &output, nil
}
func (m mockStsApi) GetCallerIdentityWithContext(aws.Context, *sts.GetCallerIdentityInput, ...request.Option) (*sts.GetCallerIdentityOutput, error) {
	return nil, nil
}
func (m mockStsApi) GetCallerIdentityRequest(*sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput) {
	return nil, nil
}

func (m mockStsApi) GetFederationToken(*sts.GetFederationTokenInput) (*sts.GetFederationTokenOutput, error) {
	return nil, nil
}
func (m mockStsApi) GetFederationTokenWithContext(aws.Context, *sts.GetFederationTokenInput, ...request.Option) (*sts.GetFederationTokenOutput, error) {
	return nil, nil
}
func (m mockStsApi) GetFederationTokenRequest(*sts.GetFederationTokenInput) (*request.Request, *sts.GetFederationTokenOutput) {
	return nil, nil
}

func (m mockStsApi) GetSessionToken(*sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
	return nil, nil
}
func (m mockStsApi) GetSessionTokenWithContext(aws.Context, *sts.GetSessionTokenInput, ...request.Option) (*sts.GetSessionTokenOutput, error) {
	return nil, nil
}
func (m mockStsApi) GetSessionTokenRequest(*sts.GetSessionTokenInput) (*request.Request, *sts.GetSessionTokenOutput) {
	return nil, nil
}
