package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

const (
	idPrefix = "id-"
	imageNameSufix = "/photo.jpg"
)

func serverSizeResizePlusPic(picURL string, size int) string {

	log.Printf("URL: %s", picURL)
	if !strings.HasSuffix(picURL, imageNameSufix) {
		log.Printf("Not valid profile picture format")
		return picURL
	}

	sizedImageName := fmt.Sprintf("/s%d%s", size, imageNameSufix)
	log.Printf("Sized image name: %s", sizedImageName)

	return strings.Replace(picURL, imageNameSufix, sizedImageName, -1)

}

func makeID(email string) string {

	// normalize
	s := strings.TrimSpace(email)
	s = strings.ToLower(s)

	// encode
	encoded := base64.StdEncoding.EncodeToString([]byte(s))

	return idPrefix + encoded

}

func parseEmail(id string) (email string, err error) {

	// check format
	if !strings.HasPrefix(id, idPrefix) {
		return "", fmt.Errorf("Invalid ID format: %s", id)
	}

	// trim
	id2 := strings.TrimPrefix(id, idPrefix)

	// decode
	decoded, err := base64.StdEncoding.DecodeString(id2)
	if err != nil {
		return "", err
	}

	return string(decoded), nil

}

func makeUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Error while getting id: %v\n", err)
	}
	return idPrefix + id.String()
}
