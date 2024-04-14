package local

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/dgraph-io/badger/v4"
	"github.com/vmihailenco/msgpack"

	"github.com/Dafaque/sshaman/internal/config"
	"github.com/Dafaque/sshaman/internal/credentials"
)

const (
	badgerDir string = "db"
)

type manager struct {
	db       *badger.DB
	location string
}

func NewManager(cfg *config.Config) (credentials.Manager, error) {
	dbLocation := path.Join(cfg.Home, badgerDir)
	opts := badger.DefaultOptions(dbLocation)
	opts.Logger = nil
	opts.MetricsEnabled = false
	opts.VerifyValueChecksum = true
	opts.EncryptionKey = cfg.EncryptionKey
	opts.IndexCacheSize = 100 << 20
	opts.NumVersionsToKeep = 1
	opts.CompactL0OnClose = true
	opts.ValueLogFileSize = 1024 * 1024 * 10
	opts.InMemory = false

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &manager{
		db:       db,
		location: dbLocation,
	}, nil
}

func (m *manager) Get(name string) (*credentials.Credentials, error) {
	tx := m.db.NewTransaction(false)
	item, err := tx.Get([]byte(name))
	if err != nil {
		return nil, err
	}
	var creds credentials.Credentials
	err = item.Value(func(val []byte) error {
		return m.decodeCredentials(val, &creds)
	})
	return &creds, err
}

func (m *manager) Set(name string, cred *credentials.Credentials, force bool) error {
	_, err := m.Get(name)
	if errors.Is(err, badger.ErrKeyNotFound) && !force {
		return fmt.Errorf("connection %s is already exists; use -force to override", name)
	} else if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
		return err
	}
	return m.db.Update(func(txn *badger.Txn) error {
		buf := bytes.NewBuffer(nil)
		encoder := msgpack.NewEncoder(buf)
		if err := encoder.Encode(cred); err != nil {
			return err
		}
		return txn.Set([]byte(name), buf.Bytes())
	})
}

func (m *manager) Del(alias string) error {
	return m.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(alias))
	})
}

func (m *manager) List() ([]*credentials.Credentials, error) {
	var creds []*credentials.Credentials = make([]*credentials.Credentials, 0)
	err := m.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var cred credentials.Credentials
				if err := m.decodeCredentials(val, &cred); err != nil {
					return err
				}
				creds = append(creds, &cred)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return creds, err
}

func (m *manager) Drop() error {
	return m.db.DropAll()
}

func (m *manager) Done() error {
	if err := m.db.RunValueLogGC(0.5); err != nil {
		if !errors.Is(err, badger.ErrNoRewrite) {
			fmt.Println(err)
		}
	}
	return m.db.Close()
}

func (m *manager) decodeCredentials(val []byte, creds *credentials.Credentials) error {
	creds.Source = m.location
	buf := bytes.NewBuffer(val)
	decoder := msgpack.NewDecoder(buf)
	return decoder.Decode(creds)
}
