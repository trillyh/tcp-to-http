package headers 

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test invalid letters
	headers = NewHeaders()
	data = []byte (" H©st: localhost:42069\r\n\r\n ")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test append to headers with the same field-name
	headers = NewHeaders()
	data = []byte ("set-person: tung\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, "tung", headers.Get("set-person"))
	assert.Equal(t, 18, n)
	assert.False(t, done)

	data = []byte ("set-person: trilly\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, 20, n)
	assert.False(t, done)

	// this one capitalized
	data = []byte ("Set-Person: trillion\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, "tung, trilly, trillion", headers.Get("set-person"))
	assert.Equal(t, 22, n)
	assert.False(t, done)
}
