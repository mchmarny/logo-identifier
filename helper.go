package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	me "github.com/mchmarny/gcputil/metric"
)

const (
	idPrefix       = "id-"
	imageNameSufix = "/photo.jpg"
)

var (
	metricClient *me.Client
)

func initAuth(ctx context.Context) {
	c, err := me.NewClient(ctx)
	if err != nil {
		logger.Fatalf("Error while initializing metrics client %v", err)
	}
	metricClient = c
}

func serverSizeResizePlusPic(picURL string, size int) string {

	logger.Printf("URL: %s", picURL)
	if !strings.HasSuffix(picURL, imageNameSufix) {
		logger.Printf("Not valid profile picture format")
		return picURL
	}

	sizedImageName := fmt.Sprintf("/s%d%s", size, imageNameSufix)
	logger.Printf("Sized image name: %s", sizedImageName)

	return strings.Replace(picURL, imageNameSufix, sizedImageName, -1)

}

func makeDailySessionID(email string) string {
	emailID := makeID(email)
	datePrefix := time.Now().UTC().Format("2006-01-02")
	return fmt.Sprintf("s-%s-%s", datePrefix, emailID)
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
		logger.Fatalf("Error while getting id: %v\n", err)
	}
	return idPrefix + id.String()
}
