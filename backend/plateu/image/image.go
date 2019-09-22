package image

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/fidellr/jastip/backend/plateu/models"
	"github.com/pkg/errors"

	"github.com/fidellr/jastip/backend/plateu"
	"github.com/fidellr/jastip/backend/plateu/repository"
)

type service struct {
	repository     repository.ImageRepository
	validator      plateu.Validate
	contextTimeout time.Duration
}

type requirement func(*service)
type object interface {
	Find(string) *goquery.Selection
}

type FBProfileData struct {
	ProfileURL string
	ImageURL   string
}

func Repository(repository repository.ImageRepository) requirement {
	return func(s *service) {
		s.repository = repository
	}
}

func Timeout(timeout time.Duration) requirement {
	return func(s *service) {
		s.contextTimeout = timeout
	}
}

func Validator(validator plateu.Validate) requirement {
	return func(s *service) {
		s.validator = validator
	}
}

func NewService(reqs ...requirement) plateu.ImageUsecase {
	s := new(service)
	for _, option := range reqs {
		option(s)
	}

	return s
}

func GetImageURL(doc *goquery.Document) (string, error) {
	// s := QuerySelector(doc, ".1kf.img")
	s := QuerySelector(doc, ".XjzKX > div span img")
	log.Println(s.Html())
	if s.Length() == 0 {
		s = QuerySelector(doc, "img.scaledImageFitWidth")
		log.Fatalf("s is len(%d)", s.Length())
	}

	url, ok := s.Attr("src")
	if !ok {
		return "", errors.New("cannot find image url")
	}

	return url, nil
}

func GuessImageContentType(file *os.File) (format string, err error) {
	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

// func ParseProfile(ctx context.Context, html, profileURL string) (*FBProfileData, error) {
// 	fb := FBProfileData{ProfileURL: profileURL}
// 	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
// 	if err != nil {
// 		log.Fatalf("Failed read html : %s", err.Error())
// 		return nil, err
// 	}

// 	fb.ImageURL, err = GetImageURL(doc)
// 	if err != nil {
// 		log.Fatalf("Failed to get image url : %s", err.Error())
// 		return nil, err
// 	}
// 	image := new(models.Image)
// 	thumborizeURL, err := thumborizeImage(ctx, fb.ImageURL, image)
// 	if err != nil {
// 		log.Fatalf("Failed to thumborized image url : %s", err.Error())
// 		return nil, err
// 	}

// 	fb.ProfileURL = profileURL
// 	fb.ImageURL = thumborizeURL

// 	return &fb, nil
// }

// func Parse(ctx context.Context, url string) (*FBProfileData, error) {
// 	doc, err := goquery.NewDocument(url)
// 	if err != nil {
// 		log.Fatalf("Failed to documentized fb profile : %s", err.Error())
// 		return nil, err
// 	}

// 	s := QuerySelector(doc, "._11kf.img")
// 	cmt, err := s.Html()
// 	if err != nil {
// 		log.Fatalf("Failed to read html : %s", err.Error())
// 		return nil, err
// 	}

// 	return ParseProfile(ctx, cmt[5:len(cmt)-4], url)
// }

func (s *service) StoreImage(ctx context.Context, m *models.Image) (err error) {
	if ctx == nil {
		err = plateu.ErrContextNil
		return err
	}

	if err = s.validator.ValidateStruct(m); err != nil {
		err = errors.Wrap(err, "error validating image")
		return err
	}

	m.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err = s.repository.StoreImage(ctx, m)
	if err != nil {
		err = errors.Wrap(err, "error storing image")
		return err
	}

	return nil
}

func (s *service) FetchImages(ctx context.Context, filter *plateu.Filter) ([]*models.Image, string, error) {
	if ctx == nil {
		err := plateu.ErrContextNil
		return nil, "", err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if filter.Num == 0 {
		filter.Num = int(3)
	}

	images, nextCursor, err := s.repository.FetchImages(ctx, filter)
	if err != nil {
		return nil, "", err
	}

	return images, nextCursor, nil
}

func (s *service) GetImageByID(ctx context.Context, imageID string) (*models.Image, error) {
	if ctx == nil {
		err := plateu.ErrContextNil
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	image, err := s.repository.GetImageByID(ctx, imageID)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *service) UpdateImageByID(ctx context.Context, imageID string, m *models.Image) (err error) {
	if ctx == nil {
		err = plateu.ErrContextNil
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	done := make(chan bool)
	go func() {
		if err != nil {
			log.Printf("Failed to parse docs : %s", err.Error())
			done <- true
			return
		}

		done <- true
		return
	}()

	if <-done {
		err = s.repository.UpdateImageByID(ctx, imageID, m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) RemoveImageByID(ctx context.Context, imageID string) (err error) {
	if ctx == nil {
		err = plateu.ErrContextNil
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err = s.repository.RemoveImageByID(ctx, imageID); err != nil {
		return err
	}

	return nil
}
