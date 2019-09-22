package mongo

import (
	"context"
	"log"

	"github.com/fidellr/jastip/backend/rover"

	"github.com/fidellr/jastip/backend/rover/models"
	"github.com/fidellr/jastip/backend/rover/repository"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	contentCollectionName = "screen_contents"
)

type contentMongoRepository struct {
	Session *mgo.Session
	DBName  string
}

type contentRequirement func(*contentMongoRepository)

func ContentSession(session *mgo.Session) contentRequirement {
	return func(c *contentMongoRepository) {
		c.Session = session
	}
}

func ContentDBName(dbName string) contentRequirement {
	return func(c *contentMongoRepository) {
		c.DBName = dbName
	}
}

func NewContentMongo(reqs ...contentRequirement) repository.ContentRepository {
	repo := new(contentMongoRepository)
	for _, req := range reqs {
		req(repo)
	}

	return repo
}

func (u *contentMongoRepository) CreateScreenContent(ctx context.Context, m *models.Content) (err error) {
	session := u.Session.Clone()
	defer session.Close()

	err = session.DB(u.DBName).C(contentCollectionName).Insert(m)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (u *contentMongoRepository) FetchContent(ctx context.Context, filter *rover.Filter) ([]*models.Content, string, error) {
	session := u.Session.Clone()
	defer session.Close()

	query := make(bson.M)
	if filter.Cursor != "" {
		createdAt, err := rover.DecodeTime(filter.Cursor)
		if err == nil {
			query["cursor"] = bson.M{"$lt": createdAt}
		} else {
			log.Fatalf("Failed to pass cursor %s : %s \n", filter.Cursor, err.Error())
		}
	}

	if filter.RoleName != "" {
		query["role"] = bson.M{"role_name": filter.RoleName}
	}

	var m []*models.Content
	err := session.DB(u.DBName).C(contentCollectionName).Find(query).Limit(filter.Num).Sort("-created_at").All(&m)
	if err != nil {
		log.Fatalf("Failed to fetch screen content : %s", err.Error())
		return nil, "", err
	}

	if len(m) == 0 {
		return make([]*models.Content, 0), "", err
	}

	lastContents := m[len(m)-1]
	nextCursor := rover.EncodeTime(lastContents.CreatedAt)
	return m, nextCursor, nil
}

func (u *contentMongoRepository) GetContentByScreen(ctx context.Context, screenName string) (*models.Content, error) {
	session := u.Session.Clone()
	defer session.Close()

	var m *models.Content

	if err := session.DB(u.DBName).C(contentCollectionName).Find(bson.M{"screen": screenName}).One(&m); err != nil {
		log.Fatalf("Failed to get screen content : %s", err.Error())
		return nil, err
	}

	return m, nil
}

func (u *contentMongoRepository) UpdateByContentID(ctxt context.Context, contentID string, m *models.Content) (err error) {
	session := u.Session.Clone()
	defer session.Close()

	idB := bson.ObjectIdHex(contentID)
	err = session.DB(u.DBName).C(contentCollectionName).Update(bson.M{"_id": idB}, m)
	if err != nil {
		log.Fatalf("Failed to update content by screen : %s", err.Error())
		return err
	}

	return nil
}
