package elestio

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthSignIn(t *testing.T) {
	c, err := NewClient(os.Getenv("ELESTIO_INTEGRATION_EMAIL"), os.Getenv("ELESTIO_INTEGRATION_API_KEY"))
	require.NoError(t, err, "expected no error")
	require.NotNil(t, c, "expected non-nil client")
}
