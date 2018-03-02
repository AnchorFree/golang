package kv

import (
	"errors"
	bolt "github.com/coreos/bbolt"
	"strings"
)

type Bolt struct {
	db      *bolt.DB
	buckets map[string]bool
}

func parseAddress(prefix string) (bucket, key string) {

	i := strings.LastIndex(prefix, "/")
	if i > 0 {
		bucket = prefix[0:i]
		if (i + 1) < len(prefix) {
			key = prefix[i+1:]
		}
	} else {
		key = prefix
		bucket = "default"
	}
	return

}

func (b *Bolt) bucketExists(bucket string) bool {

	_, ok := b.buckets[bucket]
	return ok

}

func (b *Bolt) Init(opts []string) error {

	if len(opts) < 1 {
		return errors.New("filename required to init bolt store")
	}

	db, err := bolt.Open(opts[0], 0600, nil)
	if err != nil {
		return err
	}
	b.db = db
	b.buckets = map[string]bool{}
	_ = b.CreateBucket("default")

	bkts, err := b.Buckets()
	if err != nil {
		for _, bkt := range bkts {
			b.buckets[bkt] = true
		}
	}
	return nil

}

func (b *Bolt) CreateBucket(bucket string) error {

	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	if err == nil {
		b.buckets[bucket] = true
	}
	return err

}

func (b *Bolt) Buckets() ([]string, error) {

	buckets := []string{}
	err := b.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
	})
	return buckets, err

}

func (b *Bolt) Get(prefix string) ([]byte, error) {

	var val []byte
	bucket, key := parseAddress(prefix)
	err := b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if v != nil {
			val = make([]byte, len(v))
			copy(val, v)
			return nil
		}
		return errors.New("Key not found")
	})
	if err != nil {
		return nil, err
	}
	return val, nil

}

func (b *Bolt) Delete(prefix string) error {

	bucket, key := parseAddress(prefix)
	err := b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Delete([]byte(key))
	})
	return err
}

func (b *Bolt) Put(prefix string, value []byte) error {

	bucket, key := parseAddress(prefix)
	if !b.bucketExists(bucket) {
		err := b.CreateBucket(bucket)
		if err != nil {
			return err
		}
	}
	err := b.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))
		err := bkt.Put([]byte(key), value)
		return err
	})
	return err

}

func (b *Bolt) List(bucket string) ([]string, error) {

	if bucket == "" {
		return b.Buckets()
	}

	list := []string{}

	err := b.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bkt := tx.Bucket([]byte(bucket))

		bkt.ForEach(func(k, v []byte) error {
			list = append(list, string(k))
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil

}

func (b *Bolt) DeleteTree(prefix string) error {

	err := b.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(prefix))
		return err
	})
	return err

}

func (b *Bolt) Close() error {
	return b.db.Close()
}
