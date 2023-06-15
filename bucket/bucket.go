package bucket

import (
	"context"

	"github.com/dgraph-io/badger/v4"
	"github.com/naivary/objst/logger"
	"github.com/naivary/objst/object"
)

// Bucket is the actual object storage
// containing all objects in a flat hierachy.
type Bucket struct {
	store  *badger.DB
	logger *logger.Logger
}

func New(opts badger.Options) (*Bucket, error) {
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	b := &Bucket{
		store:  db,
		logger: logger.New(context.Background()),
	}
	return b, nil
}

func (b Bucket) Create(obj *object.Object) error {
	return b.store.Update(func(txn *badger.Txn) error {
		if err := obj.IsValid(); err != nil {
			return err
		}
		data, err := obj.Marshal()
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(obj.ID), data)
		return txn.SetEntry(e)
	})
}

func (b Bucket) Get(id string) (*object.Object, error) {
	var obj object.Object
	err := b.store.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		data := make([]byte, item.ValueSize())
		if _, err := item.ValueCopy(data); err != nil {
			return err
		}
		return obj.Unmarshal(data)
	})
	return &obj, err
}

func (b Bucket) Delete(id string) error {
	return b.store.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
}
