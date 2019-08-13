package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidAccessToken(t *testing.T) {

	url := "https://storage.googleapis.com/kdemo-logos/gmail.png"

	ctx := context.Background()
	desc, err := getLogoFromURL(ctx, url)

	assert.Nil(t, err)
	assert.NotEmpty(t, desc)
}
