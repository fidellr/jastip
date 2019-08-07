package plateu

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/fidellr/jastip_way/backend/plateu/models"
	"github.com/globalsign/mgo/bson"
)

var (
	timeFormat = "2006-01-02T15:04:05.999Z07:00"
)

type ImageUsecase interface {
	StoreImage(ctx context.Context, m *models.Image) error
	FetchImages(ctx context.Context, filter *Filter) ([]*models.Image, string, error)
	GetImageByID(ctx context.Context, imageID string) (*models.Image, error)
	UpdateImageByID(ctx context.Context, imageID string, m *models.Image) error
}

type Filter struct {
	Num      int
	Cursor   string
	RoleName string
}

type ErrValidation struct {
	ErrorVal error
}

func (e *ErrValidation) Error() string {
	return e.ErrorVal.Error()
}

func CreateCursor(cursorData bson.D) (string, error) {
	data, err := bson.Marshal(cursorData)
	return base64.RawURLEncoding.EncodeToString(data), err
}

func ParseCursor(c string) (cursorData bson.D, err error) {
	var data []byte
	if data, err = base64.RawURLEncoding.DecodeString(c); err != nil {
		return
	}

	err = bson.Unmarshal(data, &cursorData)
	return
}

func EncodeTime(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}

func DecodeTime(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}
