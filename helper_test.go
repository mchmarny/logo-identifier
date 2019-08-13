package main

import (
	"log"
	"strings"
	"testing"
)


func TestServerSizeResizePlusPic(t *testing.T) {

	testImage := "https://lh6.googleusercontent.com/-65SFt9rUmD0/AAAAAAAAAAI/AAAAAAAB698/8pIgz0b5NG8/photo.jpg"
	expectedImg := "https://lh6.googleusercontent.com/-65SFt9rUmD0/AAAAAAAAAAI/AAAAAAAB698/8pIgz0b5NG8/s100/photo.jpg"

	sizedImg := serverSizeResizePlusPic(testImage, 100)
	log.Printf("Sized image: %s", sizedImg)

	if sizedImg != expectedImg {
		t.Errorf("Failed to resize valid Google Plus image")
	}

}

func TestServerSizeResizeInvalidPic(t *testing.T) {

	testImage := "https://test.domain.com/someotherpic.jpg"

	sizedImg := serverSizeResizePlusPic(testImage, 100)
	log.Printf("Sized image: %s", sizedImg)

	if testImage != sizedImg {
		t.Errorf("Resize invalid Google Plus image")
	}

}

func TestID(t *testing.T) {

	testEmail := "Test@Chmarny.com"

	id1 := makeID(testEmail)
	log.Printf("ID1: %s", id1)

	testEmail2 := strings.ToLower(testEmail)

	id2 := makeID(testEmail2)
	log.Printf("ID2: %s", id2)

	if id1 != id2 {
		t.Errorf("Failed to generate case insensitive ID")
		return
	}

	email, err := parseEmail(id2)
	if err != nil {
		t.Errorf("Error parsing email: %v", err)
	}

	log.Printf("Email: %s", email)

	if email != testEmail2 {
		t.Errorf("Failed to generate case insensitive ID")
	}

}
