package main

import (
	"github.com/dgraph-io/badger/v3"
)

type DBClient struct {
	db *badger.DB
}

func NewDBClient(storeDir string) (*DBClient, error) {
	db, err := badger.Open(badger.DefaultOptions(storeDir))
	if err != nil {
		return nil, err
	}

	dc := &DBClient {
		db,
	}

	return dc, nil
}

func (dc *DBClient) Close() error {
	return dc.db.Close()
}

func (dc *DBClient) ListKeys() ([]RawScheduleRequest, error) {
	res := make([]RawScheduleRequest, 0)
	err := dc.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			tag := string(item.Key())

			var cron string
			err := item.Value(func(v []byte) error {
				cron = string(v)
				return nil
			})

			if err != nil {
				return err
			}

			r := RawScheduleRequest {
				tag,
				cron,
			}

			res = append(res, r)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (dc *DBClient) InsertEntry(tag string, cronstring string) error {
	err := dc.db.Update(func (txn *badger.Txn) error {
		err := txn.Set([]byte(tag), []byte(cronstring))
		return err
	})

	return err
}

func (dc *DBClient) RemoveEntry(tag string) error {
	err := dc.db.Update(func (txn *badger.Txn) error {
		err := txn.Delete([]byte(tag))
		return err
	})

	return err
}
