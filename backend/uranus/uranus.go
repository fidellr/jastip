package uranus

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/fidellr/jastip/backend/uranus/models"
)

var (
	timeFormat = "2006-01-02T15:04:05.999Z07:00"
)

type UserAccountUsecase interface {
	CreateUserAccount(ctx context.Context, userAccountM *models.UserAccount) error
	Fetch(ctx context.Context, filter *Filter) ([]*models.UserAccount, string, error)
	GetUserByID(ctx context.Context, uuid string) (*models.UserAccount, error)
	SuspendAccount(ctx context.Context, uuid string) (bool, error)
	RemoveAccount(ctx context.Context, uuid string) (bool, error)
	UpdateUserByID(ctx context.Context, uuid string, userAccountM *models.UserAccount) error
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
