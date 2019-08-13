package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testEmail = "Test@Domain.com"
)

func TestServerSizeResizePlusPic(t *testing.T) {
	testImage := "https://lh6.googleusercontent.com/-65SFt9rUmD0/AAAAAAAAAAI/AAAAAAAB698/8pIgz0b5NG8/photo.jpg"
	expectedImg := "https://lh6.googleusercontent.com/-65SFt9rUmD0/AAAAAAAAAAI/AAAAAAAB698/8pIgz0b5NG8/s100/photo.jpg"

	sizedImg := serverSizeResizePlusPic(testImage, 100)
	assert.NotNil(t, sizedImg)
	assert.Equalf(t, expectedImg, sizedImg, "Images sizes don't equal %s != %s", expectedImg, sizedImg)
}

func TestServerSizeResizeInvalidPic(t *testing.T) {
	testImage := "https://test.domain.com/someotherpic.jpg"
	sizedImg := serverSizeResizePlusPic(testImage, 100)

	assert.NotNil(t, sizedImg)
	assert.Equalf(t, testImage, sizedImg, "Images don't equal %s != %s", testImage, sizedImg)
}

func TestSessionID(t *testing.T) {
	id1 := makeDailySessionID(testEmail)
	assert.NotNil(t, id1)
}

func TestID(t *testing.T) {

	id1 := makeID(testEmail)
	assert.NotNil(t, id1)

	testEmail2 := strings.ToLower(testEmail)

	id2 := makeID(testEmail2)
	assert.NotNil(t, id2)
	assert.Equalf(t, id1, id2, "Sessions IDs don't equal %s != %s", id1, id2)

	email, err := parseEmail(id2)
	assert.Nil(t, err)
	assert.NotNil(t, email)
	assert.Equalf(t, email, testEmail2, "Email addresses don't equal %s != %s", email, testEmail2)

}
