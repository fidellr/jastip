package image

import (
	"context"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fidellr/jastip_way/backend/plateu/models"
)

func crawler(ctx context.Context, imageUrl string) (*models.Image, error) {
	done := make(chan bool, 2)
	var doc *goquery.Document
	var err error
	go func() {
		doc, err = goquery.NewDocument(imageUrl)
		if err != nil {
			log.Fatalf("%s", err.Error())
			<-done
			return
		}
	}()

	if <-done {
		go func() {
			log.Println("Start downloading...")
			imageUrl, err = GetImageURL(doc)
			if err != nil {
				<-done
				return
			}

			<-done
			return
		}()
	}

	var m *models.Image
	width, height := GetImageDimension(imageUrl)
	// thumborizedURL, err := thumborizeImage(ctx, imageUrl, m)
	if err != nil {
		log.Fatalf("Failed to thumborized image url : %s", err.Error())
		return nil, err
	}

	m = &models.Image{
		CreatedAt: time.Now(),
		// URL:       thumborizedURL,
		Height: height,
		Width:  width,
	}

	return m, nil
}

func QuerySelector(s object, selector string) *goquery.Selection {
	return s.Find(selector).First()
}
