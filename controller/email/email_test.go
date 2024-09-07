package email

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailFormat(t *testing.T) {
	v := "bobowu@keiper.com"
	require.True(t, validateEmailForamt(v))
	v = "bobo.wu@keiper.com"
	require.True(t, validateEmailForamt(v))
	v = "asd4a45ads4"
	require.False(t, validateEmailForamt(v))
}
