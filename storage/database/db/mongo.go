package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/asdine/storm/q"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoShim struct {
	db          *mongo.Database
	collections map[string]*mongo.Collection // type to collection mapping
	context     context.Context
	sess        mongo.Session
}

const initializedKey = "initialized"
const ttlKey = "ttl"

// OpenMongo connects to the database using the specified URI and database name and returns a handle for accessing it
func OpenMongo(dbURI string, dbName string) (*MongoShim, error) {
	spl := strings.Split(dbName, "/")
	dbName = spl[len(spl)-1]
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &MongoShim{
		db:          db,
		collections: map[string]*mongo.Collection{},
		context:     context.TODO()}, nil
}

func (s *MongoShim) ctx() context.Context {
	inTransaction := s.sess != nil
	if inTransaction {
		var sc mongo.SessionContext
		mongo.WithSession(s.context, s.sess, func(sessionContext mongo.SessionContext) error {
			sc = sessionContext
			return nil
		})
		return sc
	}
	return s.context
}

func (s *MongoShim) nameToCollection(n string) (*mongo.Collection, error) {
	if c, ok := s.collections[n]; ok {
		return c, nil
	}
	c := s.db.Collection(n)
	s.collections[n] = c

	var err error
	fr := c.FindOne(s.ctx(), bson.M{"K": initializedKey})
	if fr.Err() == mongo.ErrNoDocuments {
		err = s.assureUniqueIndex(c)
		if err != nil {
			return c, err
		}
		// just to force collection creation as implicit creation doesn't work for transactions
		_, err = c.UpdateOne(s.ctx(), bson.M{"K": initializedKey}, bson.M{"$set": bson.M{"V": true}},
			options.Update().SetUpsert(true))
	}
	return c, err
}

func (s *MongoShim) assureUniqueIndex(c *mongo.Collection) error {
	i := mongo.IndexModel{
		Keys: bson.D{
			{"K", 1},
			{"V.id", 1},
		},
		Options: options.Index().SetUnique(true).SetName("custom-keys"),
	}
	i2 := mongo.IndexModel{
		Keys: bson.D{
			{"expireAt", 1},
		},
		Options: options.Index().SetName("ttl-index").SetSparse(true).
			SetExpireAfterSeconds(0),
	}
	_, err := c.Indexes().CreateMany(s.ctx(), []mongo.IndexModel{i, i2})
	return err
}

func (s *MongoShim) objectToCollection(o interface{}) (*mongo.Collection, error) {
	t := reflect.Indirect(reflect.ValueOf(o)).Type()
	return s.typeToCollection(t)
}

func (s *MongoShim) typeToCollection(t reflect.Type) (*mongo.Collection, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return s.nameToCollection(t.Name())
}

// Get a value from a bucket
func (s *MongoShim) Get(bucketName string, key interface{}, to interface{}, opts ...*GetOptions) error {
	c, err := s.nameToCollection(bucketName)
	if err != nil {
		return err
	}
	var raw bson.Raw
	err = c.FindOne(s.ctx(), bson.M{"K": key}).Decode(&raw)
	if err != nil {
		return err
	}

	if ttlRefresh(opts) {
		tv, err := raw.LookupErr(ttlKey)
		if err == nil {
			// item had ttl set
			var ttl time.Duration
			err = tv.Unmarshal(&ttl)
			if err != nil {
				return err
			}
			_, err := c.UpdateOne(s.ctx(), bson.M{"K": key},
				bson.M{"$set": bson.M{"expireAt": time.Now().UTC().Add(ttl)}})
			if err != nil {
				return err
			}
		}
	}
	return raw.Lookup("V").Unmarshal(to)
}

// Set a key/value pair into a bucket
func (s *MongoShim) Set(bucketName string, key interface{}, value interface{}, opts ...*SetOptions) error {
	if t := reflect.ValueOf(key).Type().Kind(); t != reflect.String {
		return fmt.Errorf("required string type got %s", t.String())
	}
	c, err := s.nameToCollection(bucketName)
	if err != nil {
		return err
	}
	v := bson.M{"V": value}

	ttl := ttlDuration(opts)
	if ttl > 0 {
		v["expireAt"] = time.Now().UTC().Add(ttl)
		v[ttlKey] = ttl
	}
	_, err = c.UpdateOne(s.ctx(), bson.M{"K": key}, bson.M{"$set": v},
		options.Update().SetUpsert(true))
	return err
}

// Delete deletes a key from a bucket
func (s *MongoShim) Delete(bucketName string, key interface{}) error {
	c, err := s.nameToCollection(bucketName)
	if err != nil {
		return err
	}
	_, err = c.DeleteOne(s.ctx(), bson.M{"K": key})
	if err != nil {
		return err
	}
	return nil
}

// Begin starts a new transaction
func (s *MongoShim) Begin(_ bool) (DB, error) {
	var err error
	s.sess, err = s.db.Client().StartSession()
	if err != nil {
		return s, err
	}
	err = s.sess.StartTransaction()
	return s, err
}

func (s *MongoShim) WithBatch(enabled bool) DB {
	// Not implemented
	return s
}

// Rollback aborts the current session's uncommited transactions
func (s *MongoShim) Rollback() error {
	if s.sess == nil {
		return nil
	}
	err := s.sess.AbortTransaction(s.ctx())
	if err != nil {
		return err
	}
	s.sess.EndSession(s.ctx())
	s.sess = nil
	return nil
}

// Commit applies yet uncommited transactions of the current session
func (s *MongoShim) Commit() error {
	if s.sess == nil {
		return nil
	}
	err := s.sess.CommitTransaction(s.ctx())
	if err != nil {
		return err
	}
	s.sess.EndSession(s.ctx())
	s.sess = nil
	return nil
}

// Select a list of records that match a list of matchers. Doesn't use indexes.
func (s *MongoShim) Select(matchers ...q.Matcher) Query {
	// TODO: translate matchers into bson query to improve performance
	return &MongoQuery{matcher: q.And(matchers...), s: s}
}

// Init creates the indexes and buckets for a given structure
func (s *MongoShim) Init(data interface{}) error {
	var err error
	if str, ok := data.(string); ok {
		_, err = s.nameToCollection(str)
	} else {
		_, err = s.objectToCollection(data)
	}
	return err
}

// ReIndex rebuilds all the indexes of a bucket
func (s *MongoShim) ReIndex(data interface{}) error {
	return nil
}

func (s *MongoShim) bucketAndKey(data interface{}) (bucket string, key interface{}, err error) {
	c, err := s.objectToCollection(data)
	id := reflect.Indirect(reflect.ValueOf(data)).FieldByName("ID")
	if !id.IsValid() {
		return "", "", errors.New("no ID field provided")
	}
	return c.Name(), id.Interface(), err
}

// Save a structure
func (s *MongoShim) Save(data interface{}) error {
	bucket, key, err := s.bucketAndKey(data)
	if err != nil {
		return err
	}
	return s.Set(bucket, key, data)
}

// Update a structure
func (s *MongoShim) Update(data interface{}) error {
	return s.Save(data)
}

// DeleteStruct deletes a structure from the associated bucket
func (s *MongoShim) DeleteStruct(data interface{}) error {
	bucket, key, err := s.bucketAndKey(data)
	if err != nil {
		return err
	}
	return s.Delete(bucket, key)
}

// One returns one record by the specified index
func (s *MongoShim) One(fieldName string, value interface{}, to interface{}) error {
	c, err := s.objectToCollection(to)
	if err != nil {
		return err
	}
	var raw bson.Raw
	err = c.FindOne(s.ctx(), bson.M{"V." + strings.ToLower(fieldName): value}).Decode(&raw)
	if err != nil {
		return err
	}
	return raw.Lookup("V").Unmarshal(to)
}

// Count all the matching records
func (s *MongoShim) Count(data interface{}) (int, error) {
	c, err := s.objectToCollection(data)
	if err != nil {
		return 0, err
	}
	count, err := c.CountDocuments(s.ctx(), bson.D{})
	count-- // subtract initializedKey
	return int(count), err
}

func (s *MongoShim) findWithConstraints(filter interface{}, to interface{}, q *MongoQuery) (int64, error) {
	// *to* can be a pointer to slice of elements or pointer to an element
	var resultsToSlice bool
	var resultsCount int64
	typ := reflect.Indirect(reflect.ValueOf(to)).Type()
	if typ.Kind() == reflect.Slice {
		resultsToSlice = true
		// clear result slice
		toV := reflect.Indirect(reflect.ValueOf(to))
		toV.Set(reflect.MakeSlice(typ, 0, toV.Cap()))
		// set slice element type
		typ = typ.Elem()
	}
	c, err := s.typeToCollection(typ)
	if err != nil {
		return 0, err
	}

	cur, err := c.Find(s.ctx(), filter, q.opts())
	if err != nil {
		return resultsCount, err
	}
	defer cur.Close(s.ctx())

	var skipFirstN int64
	if q.skip != nil {
		skipFirstN = *q.skip
	}
	limitN := int64(-1)
	if q.limit != nil {
		limitN = *q.limit
	}

	for cur.Next(s.ctx()) {
		if cur.Err() != nil {
			return resultsCount, cur.Err()
		}
		var raw bson.Raw
		err = cur.Decode(&raw)
		if err != nil {
			return resultsCount, err
		}

		if raw.Lookup("K").StringValue() == initializedKey {
			continue
		}

		sliceElemPtr := reflect.New(typ).Interface()
		sliceElemV := func() reflect.Value {
			return reflect.Indirect(reflect.ValueOf(sliceElemPtr))
		}

		err = raw.Lookup("V").Unmarshal(sliceElemPtr)
		if err != nil {
			return resultsCount, err
		}

		if q.matcher != nil {
			success, err := q.matcher.Match(sliceElemV().Interface())
			if err != nil {
				return resultsCount, err
			}
			if !success {
				// skip this element
				continue
			}
		}

		if skipFirstN > 0 {
			skipFirstN--
			continue
		}
		if limitN == 0 {
			return resultsCount, nil
		}
		limitN--

		// got valid result
		resultsCount++
		if q.callback != nil {
			err = q.callback(sliceElemV().Addr().Interface())
			if err != nil {
				return resultsCount, err
			}
		} else {
			if resultsToSlice {
				// append to result slice
				reflect.Indirect(reflect.ValueOf(to)).Set(
					reflect.Append(reflect.Indirect(reflect.ValueOf(to)), sliceElemV()))
			} else {
				// set to result element
				reflect.Indirect(reflect.ValueOf(to)).Set(sliceElemV())
			}
		}
	}
	return resultsCount, nil
}

// All gets all the records of a bucket. If there are no records it returns no error and the 'to' parameter is set to an empty slice.
func (s *MongoShim) All(to interface{}) error {
	c, err := s.findWithConstraints(bson.M{}, to, &MongoQuery{})
	if err != nil {
		return err
	}
	// compatibility quirk
	if c == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Close disconnects the db client from the database
func (s *MongoShim) Close() error {
	return s.db.Client().Disconnect(s.ctx())
}

type MongoQuery struct {
	orderBy []string
	reverse bool

	skip     *int64
	limit    *int64
	callback func(interface{}) error

	matcher q.Matcher
	s       *MongoShim
}

// Skip matching records by the given number
func (q *MongoQuery) Skip(i int) Query {
	v := int64(i)
	q.skip = &v
	return q
}

// Limit the results by the given number
func (q *MongoQuery) Limit(i int) Query {
	v := int64(i)
	q.limit = &v
	return q
}

// Order by the given fields, in descending precedence, left-to-right.
func (q *MongoQuery) OrderBy(str ...string) Query {
	q.orderBy = str
	return q
}

// Reverse the order of the results
func (q *MongoQuery) Reverse() Query {
	q.reverse = true
	return q
}

func (q *MongoQuery) opts() *options.FindOptions {
	opts := options.Find()
	// TODO: can be enabled when matchers are fully in bson otherwise results mismatch
	//if q.limit != nil {
	//	opts.SetLimit(*q.limit)
	//}
	//if q.skip != nil {
	//	opts.SetSkip(*q.skip)
	//}
	if len(q.orderBy) > 0 {
		sortOrder := 1 // ascending
		if q.reverse {
			sortOrder = -1
		}
		var sort bson.D
		for _, field := range q.orderBy {
			sort = append(sort, bson.E{
				Key:   "V." + strings.ToLower(field),
				Value: sortOrder,
			})
		}
		opts.SetSort(sort)
	}
	return opts
}

// Find a list of matching records
func (q *MongoQuery) Find(to interface{}) error {
	c, err := q.s.findWithConstraints(bson.M{}, to, q)
	if err != nil {
		return err
	}
	// compatibility quirk
	if c == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// First gets the first matching record
func (q *MongoQuery) First(to interface{}) error {
	q.Limit(1)
	c, err := q.s.findWithConstraints(bson.M{}, to, q)
	if err != nil {
		return err
	}
	// compatibility quirk
	if c == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Execute the given function for each element
func (q *MongoQuery) Each(kind interface{}, fn func(interface{}) error) error {
	q.callback = fn
	_, err := q.s.findWithConstraints(bson.M{}, kind, q)
	return err
}
