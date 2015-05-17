package mongostorage

import (
	"io"
	"errors"
	"blob-server/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStorage struct {
	session *mgo.Session
	gridFs *mgo.GridFS
}

type MongoIdError struct {
	err error
}

func (err MongoIdError) NotFound() bool {
	return err.err == mgo.ErrNotFound
}

func (err MongoIdError) IllFormed() bool {
	return err.err == ErrInvalidHexId
}

func (err MongoIdError) Error() string {
	return err.err.Error()
}

var (
	ErrInvalidHexId = errors.New("invalid hex id")
)

func Start(url string, prefix string) (*MongoStorage, error) {
	session, err := mgo.Dial(url)
	if (err != nil) {
		return nil, err
	}
	mongoStorage := &MongoStorage{
		session: session,
		gridFs: session.DB("").GridFS(prefix),
	}
	return mongoStorage, nil
}

func (mongo MongoStorage) Stop() {
	mongo.session.Close()
}

func (mongo MongoStorage) Post(reader io.Reader) (string, error) {
	file, err := mongo.gridFs.Create("")
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}
	return file.Id().(bson.ObjectId).Hex(), nil
}

func doIdRequest(id string, f func() error) storage.IdRequestError {
	if (!bson.IsObjectIdHex(id)) {
		return MongoIdError{ErrInvalidHexId}
	}
	err := f()
	if (err == nil) {
		return nil
	} else {
		return MongoIdError{err}
	}
}

func (mongo MongoStorage) Get(id string, writer io.Writer) storage.IdRequestError {
	return doIdRequest(id, func() error {
		file, err := mongo.gridFs.OpenId(bson.ObjectIdHex(id))
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
}

func (mongo MongoStorage) GetList() ([]string, error) {
	var ids []bson.ObjectId
	err := mongo.gridFs.Find(nil).Distinct("_id", &ids)
	if err != nil {
		return nil, err
	}
	files := make([]string, len(ids))
	for i, id := range ids {
		files[i] = id.Hex()
	}
	return files, nil
}

func (mongo MongoStorage) Delete(id string) storage.IdRequestError {
	return doIdRequest(id, func() error {
		return mongo.gridFs.RemoveId(bson.ObjectIdHex(id))
	})
}
