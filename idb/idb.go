package idb

import (
	"encoding/json"
	"fmt"
	tcp_model "github.com/k8spacket/plugins/nodegraph/metrics/nodegraph/model"
	tls_model "github.com/k8spacket/plugins/tls-parser/metrics/model"
	"go.etcd.io/bbolt"
	"hash/fnv"
)

var BUCKET = []byte("bucket")

type DB[T tls_model.TLSDetails | tls_model.TLSConnection | tcp_model.ConnectionItem] struct {
	db *bbolt.DB
}

func StartDB[T tls_model.TLSDetails | tls_model.TLSConnection | tcp_model.ConnectionItem](dbname string) (*DB[T], error) {

	database, err := bbolt.Open(fmt.Sprintf("%s.db", dbname), 0600, nil)
	if err != nil {
		return nil, err
	}
	database.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket(BUCKET)
		if err != nil {
			return err
		}
		return nil
	})
	return &DB[T]{database}, nil
}

func (k *DB[T]) Close() error {
	return k.db.Close()
}

func (k *DB[T]) Read(key string) (T, error) {
	var value T
	return value, k.db.View(func(tx *bbolt.Tx) error {
		item := tx.Bucket(BUCKET).Get([]byte(key))
		err := json.Unmarshal(item, &value)
		if err != nil {
			return err
		}
		return nil
	})
}

func (k *DB[T]) ReadAll() ([]T, error) {
	var value []T
	return value, k.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BUCKET)
		c := b.Cursor()
		for key, val := c.First(); key != nil; key, val = c.Next() {
			var obj T
			err := json.Unmarshal(val, &obj)
			if err != nil {
				return err
			}
			value = append(value, obj)
		}

		return nil
	})
}

func (k *DB[T]) Upsert(key string, value T) error {
	return k.db.Update(
		func(tx *bbolt.Tx) error {
			val, err := json.Marshal(value)
			if err != nil {
				return err
			}
			return tx.Bucket(BUCKET).Put([]byte(key), val)
		})
}

func HashId(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
