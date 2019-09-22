package mongo

import (
	"context"
	"log"

	"github.com/fidellr/jastip/backend/plateu"

	"github.com/fidellr/jastip/backend/plateu/models"
	"github.com/fidellr/jastip/backend/plateu/repository"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	imageCollectionName = "images"
)

type imageMongoRepository struct {
	Session *mgo.Session
	DBName  string
}

type imageRequirement func(*imageMongoRepository)

func ImageSession(session *mgo.Session) imageRequirement {
	return func(c *imageMongoRepository) {
		c.Session = session
	}
}

func ImageDBName(dbName string) imageRequirement {
	return func(c *imageMongoRepository) {
		c.DBName = dbName
	}
}

func NewImageMongo(reqs ...imageRequirement) repository.ImageRepository {
	repo := new(imageMongoRepository)
	for _, req := range reqs {
		req(repo)
	}

	return repo
}

func (r *imageMongoRepository) StoreImage(ctx context.Context, m *models.Image) (err error) {
	session := r.Session.Clone()
	defer session.Close()

	err = session.DB(r.DBName).C(imageCollectionName).Insert(m)
	if err != nil {
		log.Printf("Failed to Insert image : %s", err.Error())
		return err
	}

	return nil
}

func (r *imageMongoRepository) FetchImages(ctx context.Context, filter *plateu.Filter) ([]*models.Image, string, error) {
	session := r.Session.Clone()
	defer session.Close()

	query := make(bson.M)
	if filter.Cursor != "" {
		createdAt, err := plateu.DecodeTime(filter.Cursor)
		if err == nil {
			query["cursor"] = bson.M{"$lt": createdAt}
		} else {
			log.Printf("Failed to pass cursor %s : %s \n", filter.Cursor, err.Error())
		}
	}

	if filter.RoleName != "" {
		query["role"] = bson.M{"role_name": filter.RoleName}
	}

	var m []*models.Image
	err := session.DB(r.DBName).C(imageCollectionName).Find(query).All(&m)
	if err != nil {
		log.Printf("Failed to fetch screen content : %s \n", err.Error())
		return nil, "", err
	}

	if len(m) == 0 {
		return make([]*models.Image, 0), "", err
	}

	lastImages := m[len(m)-1]
	nextCursor := plateu.EncodeTime(lastImages.CreatedAt)
	return m, nextCursor, nil
}

func (r *imageMongoRepository) GetImageByID(ctx context.Context, imageID string) (*models.Image, error) {
	session := r.Session.Clone()
	defer session.Close()

	var m *models.Image
	imageIDb := bson.ObjectIdHex(imageID)
	if err := session.DB(r.DBName).C(imageCollectionName).Find(bson.M{"_id": imageIDb}).One(&m); err != nil {
		log.Printf("Failed to get image : %s", err.Error())
		return nil, err
	}

	return m, nil
}

func (r *imageMongoRepository) UpdateImageByID(ctx context.Context, imageID string, image *models.Image) (err error) {
	session := r.Session.Clone()
	defer session.Close()

	imageIDb := bson.ObjectIdHex(imageID)
	err = session.DB(r.DBName).C(imageCollectionName).Update(bson.M{"_id": imageIDb}, image)
	if err != nil {
		log.Printf("Failed to update image : %s", err.Error())
		return err
	}

	return nil
}

func (r *imageMongoRepository) RemoveImageByID(ctx context.Context, imageID string) (err error) {
	session := r.Session.Clone()
	defer session.Close()

	imageIDb := bson.ObjectIdHex(imageID)
	if err = session.DB(r.DBName).C(imageCollectionName).Remove(bson.M{"_id": imageIDb}); err != nil {
		log.Printf("Failed to remoove image : %s", err.Error())
		return err
	}

	return nil
}
