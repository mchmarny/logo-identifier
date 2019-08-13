package main

import (
	"context"

	vision "cloud.google.com/go/vision/apiv1"
)

func getLogoFromURL(ctx context.Context, url string) (desc string, err error) {

	client, e := vision.NewImageAnnotatorClient(ctx)
	if e != nil {
		logger.Printf("Error while creating image client: %v", e)
		return "", e
	}

	image := vision.NewImageFromURI(url)
	annotations, e := client.DetectLogos(ctx, image, nil, 10)
	if e != nil {
		logger.Printf("Error while decoding logos for %s: %v", url, e)
		return "", e
	}

	if len(annotations) == 0 {
		logger.Printf("No logo desc found for: %s", url)
		return "", nil
	}

	d := annotations[0].Description
	logger.Printf("Found logos: %s", d)
	return d, nil
}
