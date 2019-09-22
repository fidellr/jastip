package image

import (
	"context"
	"log"

	"github.com/fidellr/jastip/backend/plateu/models"
	"github.com/globocom/gothumbor"
)

func thumborizeImage(ctx context.Context, imageURL string, m *models.Image) (thumborizedURL string, err error) {
	secretKey := "jastip_way_2019"
	thumborOpts := gothumbor.ThumborOptions{Width: m.Width, Height: m.Height}
	thumborizedURL, err = gothumbor.GetCryptedThumborPath(secretKey, imageURL, thumborOpts)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return "", err
	}

	return thumborizedURL, nil
}
