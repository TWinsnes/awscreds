package console

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sts"
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

	sdkHelper := mockSdkHelper{}

	c := Console{}

	name, err := c.getAwsUsername(nil, sdkHelper)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "UserName", name)
}

func TestOpenConsole(t *testing.T) {
	c := Console{
		SessionDuration: "12h",
	}
	sdkHelper := mockSdkHelper{}
	browser := mockBrowser{}

	err := c.OpenConsole(browser, sdkHelper)

	assert.NoError(t, err, "expected no error")
}

type mockSdkHelper struct{}

func (mockSdkHelper) GetCallerIdentity(stsClient *sts.STS) (string, error) {
	return "sdfsd/UserName", nil
}

func (mockSdkHelper) GetFederationToken(profile string, name string, duration int64) (AwsCredentials, error) {
	var creds AwsCredentials
	return creds, nil
}

type mockBrowser struct{}

func (mockBrowser) Open(url string) error {
	return nil
}
