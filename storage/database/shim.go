package database

import (
	"github.com/asdine/storm/q"
)

// Shim represents common set of the database API.
type Shim interface {
	// Get a value from a bucket
	Get(bucketName string, key interface{}, to interface{}) error
	// Set a key/value pair into a bucket
	Set(bucketName string, key interface{}, value interface{}) error
	// Delete deletes a key from a bucket
	Delete(bucketName string, key interface{}) error

	// Begin starts a new transaction
	Begin(writable bool) (Shim, error)
	WithBatch(enabled bool) Shim
	Rollback() error
	Commit() error

	// Select a list of records that match a list of matchers. Doesn't use indexes.
	Select(matchers ...q.Matcher) QueryShim

	// Init creates the indexes and buckets for a given structure
	Init(data interface{}) error
	// ReIndex rebuilds all the indexes of a bucket
	ReIndex(data interface{}) error
	// Save a structure
	Save(data interface{}) error
	// Update a structure
	Update(data interface{}) error
	// DeleteStruct deletes a structure from the associated bucket
	DeleteStruct(data interface{}) error

	// One returns one record by the specified index
	One(fieldName string, value interface{}, to interface{}) error
	// Count all the matching records
	Count(data interface{}) (int, error)
	// All gets all the records of a bucket. If there are no records it returns no error and the 'to' parameter is set to an empty slice.
	All(to interface{}) error

	Close() error
}

// QueryShim allows to operate searches.
type QueryShim interface {
	// Skip matching records by the given number
	Skip(int) QueryShim
	// Limit the results by the given number
	Limit(int) QueryShim
	// Order by the given fields, in descending precedence, left-to-right
	OrderBy(str ...string) QueryShim
	// Reverse the order of the results
	Reverse() QueryShim
	// Find a list of matching records
	Find(to interface{}) error
	// First gets the first matching record
	First(to interface{}) error
	// Execute the given function for each element
	Each(kind interface{}, fn func(interface{}) error) error
}
