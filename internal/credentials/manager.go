package credentials

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/dgraph-io/badger/v4"
	"github.com/vmihailenco/msgpack"

	"github.com/Dafaque/sshaman/v2/internal/config"
)

const (
	badgerDir string = "db"
)

type Manager struct {
	db       *badger.DB
	location string
}

func New(cfg *config.Config) (*Manager, error) {
	dbLocation := path.Join(cfg.Home, badgerDir)
	opts := badger.DefaultOptions(dbLocation)
	opts.Logger = nil
	opts.MetricsEnabled = false
	opts.VerifyValueChecksum = true
	opts.IndexCacheSize = 100 << 20
	opts.NumVersionsToKeep = 1
	opts.CompactL0OnClose = true
	opts.ValueLogFileSize = 1024 * 1024 * 10
	opts.InMemory = false

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	m := &Manager{
		db:       db,
		location: dbLocation,
	}

	return m, nil
}

func (m *Manager) Get(name string) (*Credentials, error) {
	tx := m.db.NewTransaction(false)
	item, err := tx.Get([]byte(name))
	if err != nil {
		return nil, err
	}
	var creds Credentials
	err = item.Value(func(val []byte) error {
		return m.decodeCredentials(val, &creds)
	})
	return &creds, err
}

func (m *Manager) Set(cred *Credentials, force bool) error {
	existedCred, err := m.Get(cred.Name)
	if existedCred != nil && !force {
		return NewCredentialExistsError(cred.Name)
	} else if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
		return err
	}
	return m.db.Update(func(txn *badger.Txn) error {
		buf := bytes.NewBuffer(nil)
		encoder := msgpack.NewEncoder(buf)
		if err := encoder.Encode(cred); err != nil {
			return err
		}
		return txn.Set([]byte(cred.Name), buf.Bytes())
	})
}

func (m *Manager) Del(name string) error {
	return m.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(name))
	})
}

func (m *Manager) List() ([]*Credentials, error) {
	var creds []*Credentials = make([]*Credentials, 0)
	err := m.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var cred Credentials
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

func (m *Manager) Drop() error {
	return m.db.DropAll()
}

func (m *Manager) Done() error {
	if m.db.IsClosed() {
		return nil
	}
	if err := m.db.RunValueLogGC(0.5); err != nil {
		if !errors.Is(err, badger.ErrNoRewrite) {
			fmt.Println(err)
		}
	}
	return m.db.Close()
}

func (m *Manager) decodeCredentials(val []byte, creds *Credentials) error {
	buf := bytes.NewBuffer(val)
	decoder := msgpack.NewDecoder(buf)
	return decoder.Decode(creds)
}
