package mongo

import (
	"context"
	"log"

	"github.com/fidellr/jastip/backend/uranus/repository"

	"github.com/fidellr/jastip/backend/uranus"
	"github.com/fidellr/jastip/backend/uranus/models"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	userAccountCollectionName = "user_account"
)

type userMongoRepository struct {
	Session *mgo.Session
	DBName  string
}

type userRequirement func(*userMongoRepository)

func UserSession(session *mgo.Session) userRequirement {
	return func(c *userMongoRepository) {
		c.Session = session
	}
}

func UserDBName(dbName string) userRequirement {
	return func(c *userMongoRepository) {
		c.DBName = dbName
	}
}

func NewUserMongo(reqs ...userRequirement) repository.UserAccountRepository {
	repo := new(userMongoRepository)
	for _, req := range reqs {
		req(repo)
	}

	return repo
}

func (u *userMongoRepository) CreateUserAccount(ctx context.Context, m *models.UserAccount) error {
	session := u.Session.Clone()
	defer session.Close()

	err := session.DB(u.DBName).C(userAccountCollectionName).Insert(m)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (u *userMongoRepository) Fetch(ctx context.Context, filter *uranus.Filter) ([]*models.UserAccount, string, error) {
	session := u.Session.Clone()
	defer session.Close()

	var m []*models.UserAccount
	query := make(bson.M)
	if filter.Cursor != "" {
		createdAt, err := uranus.DecodeTime(filter.Cursor)
		if err == nil {
			query["created_at"] = bson.M{"$lt": createdAt}
		} else {
			log.Fatalf("Failed to pass cursor %s : %s \n", filter.Cursor, err.Error())
		}
	}

	if filter.RoleName != "" {
		query["role"] = bson.M{"role_name": filter.RoleName}
	}

	err := session.DB(u.DBName).C(userAccountCollectionName).Find(query).Limit(filter.Num).Sort("-created_at").All(&m)
	if err != nil {
		log.Println(err.Error())
		return make([]*models.UserAccount, 0), "", err
	}

	if len(m) == 0 {
		return make([]*models.UserAccount, 0), "", err
	}

	lastUsers := m[len(m)-1]
	nextCursors := uranus.EncodeTime(lastUsers.CreatedAt)
	return m, nextCursors, nil
}

func (u *userMongoRepository) GetUserByID(ctx context.Context, uuid string) (*models.UserAccount, error) {
	session := u.Session.Clone()
	defer session.Close()

	var m *models.UserAccount
	uuidB := bson.ObjectIdHex(uuid)
	err := session.DB(u.DBName).C(userAccountCollectionName).Find(bson.M{"_id": uuidB}).One(&m)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return m, nil
}

func (u *userMongoRepository) SuspendAccount(ctx context.Context, uuid string) (bool, error) {
	session := u.Session.Clone()
	defer session.Close()

	var m *models.UserAccount
	uuidB := bson.ObjectIdHex(uuid)
	c := session.DB(u.DBName).C(userAccountCollectionName)
	err := c.Find(bson.M{"_id": uuidB}).One(&m)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	m.IsBanned = true
	err = c.Update(bson.M{"_id": uuidB}, m)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	return true, nil
}

func (u *userMongoRepository) RemoveAccount(ctx context.Context, uuid string) (bool, error) {
	session := u.Session.Clone()
	defer session.Close()

	uuidB := bson.ObjectIdHex(uuid)
	err := session.DB(u.DBName).C(userAccountCollectionName).Remove(bson.M{"_id": uuidB})
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	return true, nil
}

func (u *userMongoRepository) UpdateUserByID(ctx context.Context, uuid string, m *models.UserAccount) error {
	session := u.Session.Clone()
	defer session.Close()

	uuidB := bson.ObjectIdHex(uuid)
	err := session.DB(u.DBName).C(userAccountCollectionName).Update(bson.M{"_id": uuidB}, m)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
