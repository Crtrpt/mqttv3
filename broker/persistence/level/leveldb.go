package level

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/asdine/storm"
	"github.com/crtrpt/mqtt/server/persistence"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (

	// defaultPath is the default file to use to store the data.
	defaultPath = "leveldb.data"

	// defaultTimeout is the default timeout of the file lock.
	defaultTimeout = 250 * time.Millisecond
)

var (
	// ErrDBNotOpen indicates the bolt db file is not open for reading.
	ErrDBNotOpen  = fmt.Errorf("leveldb not opened")
	ErrJsonEncode = fmt.Errorf("json encode error")
	ErrJsonDecode = fmt.Errorf("json decode error")
)

// Store is a backend for writing and reading to leveldb persistent storage.
type Store struct {
	path        string      // the path on which to store the db file.
	opts        any         // options for configuring the leveldb instance.
	db          *leveldb.DB // the leveldb instance.
	inflightTTL int64       // the number of seconds an inflight message should be retained before being dropped.
}

// New returns a configured instance of the leveldb store.
func New(path string, opts *opt.Options) *Store {
	if path == "" || path == "." {
		path = defaultPath
	}

	if opts == nil {
		opts = opts
	}

	return &Store{
		path: path,
		opts: opts,
	}
}

// SetInflightTTL sets the number of seconds an inflight message should be kept
// before being dropped, in the event it is not delivered. Unless you have a good reason,
// you should allow this to be called by the server (in AddStore) instead of directly.
func (s *Store) SetInflightTTL(seconds int64) {
	s.inflightTTL = seconds
}

// Open opens the boltdb instance.
func (s *Store) Open() error {
	var err error
	s.db, err = leveldb.OpenFile(s.path, nil)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the boltdb instance.
func (s *Store) Close() {
	s.db.Close()
}

// WriteServerInfo writes the server info to the leveldb instance.
func (s *Store) WriteServerInfo(v persistence.ServerInfo) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	data, err := json.Marshal(&v)
	if err != nil {
		return ErrJsonEncode
	}

	err = s.db.Put([]byte("WriteServerInfo"), []byte(data), nil)
	if err != nil {
		return err
	}
	return nil
}

// WriteSubscription writes a single subscription to the leveldb instance.
func (s *Store) WriteSubscription(v persistence.Subscription) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	data, err := json.Marshal(&v)
	if err != nil {
		return ErrJsonEncode
	}
	err = s.db.Put([]byte("WriteSubscription"), []byte(data), nil)
	if err != nil {
		return err
	}
	return nil
}

// WriteInflight writes a single inflight message to the leveldb instance.
func (s *Store) WriteInflight(v persistence.Message) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	data, err := json.Marshal(&v)
	if err != nil {
		return ErrJsonEncode
	}
	err = s.db.Put([]byte("WriteInflight"), []byte(data), nil)
	if err != nil {
		return err
	}
	return nil
}

// WriteRetained writes a single retained message to the boltdb instance.
func (s *Store) WriteRetained(v persistence.Message) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	data, err := json.Marshal(&v)
	if err != nil {
		return ErrJsonEncode
	}
	err = s.db.Put([]byte("WriteRetained"), []byte(data), nil)
	if err != nil {
		return err
	}
	return nil
}

// WriteClient writes a single client to the leveldb instance.
func (s *Store) WriteClient(v persistence.Client) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	data, err := json.Marshal(&v)
	if err != nil {
		return ErrJsonEncode
	}
	err = s.db.Put([]byte("WriteClient"), []byte(data), nil)
	if err != nil {
		return err
	}
	return nil
}

// DeleteSubscription deletes a subscription from the boltdb instance.
func (s *Store) DeleteSubscription(id string) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	err := s.db.DeleteStruct(&persistence.Subscription{
		ID: id,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteClient deletes a client from the boltdb instance.
func (s *Store) DeleteClient(id string) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	err := s.db.DeleteStruct(&persistence.Client{
		ID: id,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteInflight deletes an inflight message from the boltdb instance.
func (s *Store) DeleteInflight(id string) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	err := s.db.DeleteStruct(&persistence.Message{
		ID: id,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteRetained deletes a retained message from the boltdb instance.
func (s *Store) DeleteRetained(id string) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	err := s.db.DeleteStruct(&persistence.Message{
		ID: id,
	})
	if err != nil {
		return err
	}

	return nil
}

// ReadSubscriptions loads all the subscriptions from the boltdb instance.
func (s *Store) ReadSubscriptions() (v []persistence.Subscription, err error) {
	if s.db == nil {
		return v, ErrDBNotOpen
	}

	err = s.db.Find("T", persistence.KSubscription, &v)
	if err != nil && err != storm.ErrNotFound {
		return
	}

	return v, nil
}

// ReadClients loads all the clients from the boltdb instance.
func (s *Store) ReadClients() (v []persistence.Client, err error) {
	if s.db == nil {
		return v, ErrDBNotOpen
	}

	err = s.db.Find("T", persistence.KClient, &v)
	if err != nil && err != storm.ErrNotFound {
		return
	}

	return v, nil
}

// ReadInflight loads all the inflight messages from the boltdb instance.
func (s *Store) ReadInflight() (v []persistence.Message, err error) {
	if s.db == nil {
		return v, ErrDBNotOpen
	}

	err = s.db.Find("T", persistence.KInflight, &v)
	if err != nil && err != storm.ErrNotFound {
		return
	}

	return v, nil
}

// ReadRetained loads all the retained messages from the boltdb instance.
func (s *Store) ReadRetained() (v []persistence.Message, err error) {
	if s.db == nil {
		return v, ErrDBNotOpen
	}

	err = s.db.Find("T", persistence.KRetained, &v)
	if err != nil && err != storm.ErrNotFound {
		return
	}

	return v, nil
}

// ReadServerInfo loads the server info from the boltdb instance.
func (s *Store) ReadServerInfo() (v persistence.ServerInfo, err error) {
	if s.db == nil {
		return v, ErrDBNotOpen
	}

	err = s.db.One("ID", persistence.KServerInfo, &v)
	if err != nil && err != storm.ErrNotFound {
		return
	}

	return v, nil
}

// ClearExpiredInflight deletes any inflight messages older than the provided unix timestamp.
func (s *Store) ClearExpiredInflight(expiry int64) error {
	if s.db == nil {
		return ErrDBNotOpen
	}

	var v []persistence.Message
	err := s.db.Find("T", persistence.KInflight, &v)
	if err != nil && err != storm.ErrNotFound {
		return err
	}

	for _, m := range v {
		if m.Created < expiry || m.Created == 0 {
			err := s.db.DeleteStruct(&persistence.Message{ID: m.ID})
			if err != nil {
				return err
			}
		}
	}

	return nil
}