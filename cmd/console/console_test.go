package console

import (
	"testing"

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

	c := Console{
		SessionDuration: "3h",
	}

	output, err := c.parseSessionDuration()

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expected, output)
}

func TestParseSessionDurationMinutes(t *testing.T) {

	t.Parallel()
	var expected int64 = 900

	c := Console{
		SessionDuration: "15m",
	}

	output, err := c.parseSessionDuration()

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expected, output)
}

func TestParseSessionDurationSeconds(t *testing.T) {

	t.Parallel()
	var expected int64 = 1000

	c := Console{
		SessionDuration: "1000s",
	}

	output, err := c.parseSessionDuration()

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expected, output)
}

func TestParseSessionDurationDefault(t *testing.T) {

	t.Parallel()
	var expected int64 = 1200

	c := Console{
		SessionDuration: "1200",
	}

	output, err := c.parseSessionDuration()

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expected, output)
}

func TestParseSessionDurationBadSuffix(t *testing.T) {

	t.Parallel()
	var expected int64

	c := Console{
		SessionDuration: "1200x",
	}

	output, err := c.parseSessionDuration()

	assert.Error(t, err, "Expected an error")
	assert.Equal(t, expected, output)
}
