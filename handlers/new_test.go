package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValidUrl(t *testing.T) {
	var token = "abcdef0123"

	u, err := getValidUrl("http://valid.url.co.uk", token)
	assert.Nil(t, err)
	assert.Equal(t, u, "http://valid.url.co.uk/abcdef0123")

	// Make sure trailing slash is not causing issues
	u, err = getValidUrl("http://valid.url.co.uk/", token)
	assert.Nil(t, err)
	assert.Equal(t, u, "http://valid.url.co.uk/abcdef0123")

	// URL formats that are not accepted
	_, err = getValidUrl("////", token)
	assert.NotNil(t, err)

	_, err = getValidUrl("http://invalid.url.co.uk/sample-path", token)
	assert.NotNil(t, err)

	_, err = getValidUrl("http://invalid.url.co.uk/?some=query", token)
	assert.NotNil(t, err)

	_, err = getValidUrl("http://invalid.url.co.uk/#a-fragment", token)
	assert.NotNil(t, err)
}
