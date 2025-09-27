package body

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBody(t *testing.T) {
	body := NewBody()
	data := []byte("hello world @!@\n")
	body.ContentLength = len(data)
	n, done, err := body.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, body)
	assert.Equal(t, "hello world @!@\n", body.Body)
	assert.Equal(t, 16, n)
	assert.True(t, done)
}
