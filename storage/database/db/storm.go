package db

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"
)

type StormShim struct {
	db *storm.DB
	tx storm.Node

	hasTTL     bool
	once       sync.Once
	ticker     *time.Ticker
	tickerDone chan struct{}
}

const ttlBucket = "ttlAttribute"

type ttlAttribute struct {
	Expires int64 // seconds since epoch
	TTL     time.Duration
	Bucket  string
	Key     interface{}
}

func OpenStorm(path string) (*StormShim, error) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return nil, err
	}
	db, err := storm.Open(path, storm.Codec(msgpack.Codec))
	s := &StormShim{db: db, tx: db}
	var hasTTL bool
	s.Get(ttlBucket, "ttl", &hasTTL)
	s.hasTTL = hasTTL
	if hasTTL {
		s.once.Do(s.fireTTLCleaner)
	}
	return s, err
}

// Get a value from a bucket
func (s *StormShim) Get(bucketName string, key interface{}, to interface{}, opts ...*GetOptions) error {
	if s.hasTTL && ttlRefresh(opts) {
		var a ttlAttribute
		err := s.tx.Get(ttlBucket, key, &a)
		if err == nil {
			// TTL entry found
			now := time.Now().Unix()
			if a.Expires <= now {
				s.elementExpired(&a)
			} else {
				a.Expires = now + (a.TTL.Milliseconds() / 1e3)
				s.tx.Set(ttlBucket, key, a)
			}
		}
	}
	return s.tx.Get(bucketName, key, to)
}

// Set a key/value pair into a bucket
func (s *StormShim) Set(bucketName string, key interface{}, value interface{}, opts ...*SetOptions) error {
	ttl := ttlDuration(opts)
	if ttl > 0 {
		if !s.hasTTL {
			s.hasTTL = true
			err := s.tx.Set(ttlBucket, "ttl", true)
			if err != nil {
				return err
			}
			s.once.Do(s.fireTTLCleaner)
		}
		err := s.tx.Set(ttlBucket, key,
			&ttlAttribute{
				Expires: time.Now().UTC().Add(ttl).Unix(),
				TTL:     ttl,
				Bucket:  bucketName,
				Key:     key,
			})
		if err != nil {
			return err
		}
	}
	return s.tx.Set(bucketName, key, value)
}

// Delete deletes a key from a bucket
func (s *StormShim) Delete(bucketName string, key interface{}) error {
	return s.tx.Delete(bucketName, key)
}

// Begin starts a new transaction
func (s *StormShim) Begin(writable bool) (DB, error) {
	tx, err := s.db.Begin(writable)
	s.tx = tx
	return s, err
}

func (s *StormShim) WithBatch(enabled bool) DB {
	s.tx = s.tx.WithBatch(enabled)
	return s
}

func (s *StormShim) Rollback() error {
	err := s.tx.Rollback()
	s.tx = s.db
	return err
}

func (s *StormShim) Commit() error {
	err := s.tx.Commit()
	s.tx = s.db
	return err
}

// Select a list of records that match a list of matchers. Doesn't use indexes.
func (s *StormShim) Select(matchers ...q.Matcher) Query {
	return StormQueryShim{q: s.tx.Select(matchers...)}
}

// Init creates the indexes and buckets for a given structure
func (s *StormShim) Init(data interface{}) error {
	return s.tx.Init(data)
}

// ReIndex rebuilds all the indexes of a bucket
func (s *StormShim) ReIndex(data interface{}) error {
	return s.tx.ReIndex(data)
}

// Save a structure
func (s *StormShim) Save(data interface{}) error {
	return s.tx.Save(data)
}

// Update a structure
func (s *StormShim) Update(data interface{}) error {
	return s.tx.Update(data)
}

// DeleteStruct deletes a structure from the associated bucket
func (s *StormShim) DeleteStruct(data interface{}) error {
	return s.tx.DeleteStruct(data)
}

// One returns one record by the specified index
func (s *StormShim) One(fieldName string, value interface{}, to interface{}) error {
	return s.tx.One(fieldName, value, to)
}

// Count all the matching records
func (s *StormShim) Count(data interface{}) (int, error) {
	return s.tx.Count(data)
}

// All gets all the records of a bucket. If there are no records it returns no error and the 'to' parameter is set to an empty slice.
func (s *StormShim) All(to interface{}) error {
	return s.tx.All(to)
}

func (s *StormShim) Close() error {
	if s.ticker != nil {
		s.ticker.Stop()
		close(s.tickerDone)
	}
	return s.db.Close()
}

type StormQueryShim struct {
	q storm.Query
}

func (s StormQueryShim) wrap(q storm.Query) Query {
	s.q = q
	return s
}

// Skip matching records by the given number
func (s StormQueryShim) Skip(i int) Query {
	return s.wrap(s.q.Skip(i))
}

// Limit the results by the given number
func (s StormQueryShim) Limit(i int) Query {
	return s.wrap(s.q.Limit(i))
}

// Order by the given fields, in descending precedence, left-to-right
func (s StormQueryShim) OrderBy(str ...string) Query {
	return s.wrap(s.q.OrderBy(str...))
}

// Reverse the order of the results
func (s StormQueryShim) Reverse() Query {
	return s.wrap(s.q.Reverse())
}

// Find a list of matching records
func (s StormQueryShim) Find(to interface{}) error {
	return s.q.Find(to)
}

// First gets the first matching record
func (s StormQueryShim) First(to interface{}) error {
	return s.q.First(to)
}

// Execute the given function for each element
func (s StormQueryShim) Each(kind interface{}, fn func(interface{}) error) error {
	return s.q.Each(kind, fn)
}

func (s *StormShim) elementExpired(a *ttlAttribute) {
	s.db.Delete(a.Bucket, a.Key)
}

func (s *StormShim) fireTTLCleaner() {
	s.ticker = time.NewTicker(1 * time.Second)
	s.tickerDone = make(chan struct{})
	go func() {
		for {
			select {
			case t := <-s.ticker.C:
				var deleteTTLs []interface{}
				s.db.Select(q.Lte("Expires", t.Unix())).
					Each(new(ttlAttribute), func(v interface{}) error {
						if a, ok := v.(*ttlAttribute); ok {
							s.elementExpired(a)
							deleteTTLs = append(deleteTTLs, a.Key)
						}
						return nil
					})
				for _, k := range deleteTTLs {
					s.db.Delete(ttlBucket, k)
				}
			case <-s.tickerDone:
				return
			}
		}
	}()
}
