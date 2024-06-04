package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	type scenario struct {
		name string
		arg  string
		test func(*testing.T, error)
	}
	scenarios := []scenario{
		{
			"Invalid characters",
			"test ",
			func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			"More than 50 characters provided",
			"F9e99Sy3gYZDcyDsTThdTBPcZ57hAGcPHTcNiS2hlJhADSgKdKo",
			func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			"Valid name provided",
			"F9e99Sy3g_-DcyDsTThdTBPcZ57hAGcPHTcNiS2hlJhADSgKdK",
			func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			s.test(t, ValidateName(s.arg))
		})
	}
}
