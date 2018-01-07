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
