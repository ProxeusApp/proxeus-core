package database

import (
	"reflect"
	"testing"

	"github.com/ProxeusApp/proxeus-core/storage/database"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"

	"github.com/asdine/storm/q"
)

type myStruct struct {
	ID   string `storm:"id"`
	Val  int
	Name string
}

var objs = map[int]myStruct{
	1: {ID: "1", Val: 11, Name: "gop"},
	2: {ID: "2", Val: 22, Name: "abe"},
	3: {ID: "3", Val: 33, Name: "abe"},
}

var insertData = []myStruct{objs[1], objs[3], objs[2]}

func testCRUD(t *testing.T, db db.DB) {
	initDB(t, db)
	defer db.Close()
	// get
	for _, o := range insertData {
		var to myStruct
		err := db.Get("myStruct", o.ID, &to)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(to, o) {
			t.Errorf("got obj '%v' expected '%v'", to, o)
		}
	}

	// count
	count, err := db.Count(&myStruct{})
	if err != nil {
		t.Error(err)
	}
	if count < len(insertData) {
		t.Error("expected count >= 3 got", count)
	}

	// all with pointer type
	var all []*myStruct
	err = db.All(&all)
	if err != nil {
		t.Error(err)
	}
	if len(all) != count {
		t.Errorf("expected len(all) = %v got %v", count, len(all))
	}
	for _, o := range insertData {
		objPtrInSlice(t, &o, all)
	}

	// all without pointer type
	var all2 []myStruct
	err = db.All(&all2)
	if err != nil {
		t.Error(err)
	}
	if len(all2) != count {
		t.Errorf("expected len(all2) = %v got %v", count, len(all2))
	}
	for _, o := range insertData {
		objInSlice(t, o, all2)
	}

	deleteData(t, db)
}

func testGetQuirks(t *testing.T, db db.DB) {
	// each for non-existing
	err := db.Select(q.Eq("ID", "non-existing")).Each(new(myStruct), func(v interface{}) error {
		t.Error("callback shouldn't be called")
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	// find for non-existing
	var obj myStruct
	var objS []myStruct
	err = db.Select(q.Eq("ID", "non-existing")).Find(&objS)
	if !database.NotFound(err) {
		t.Error("expected not found err, got", err)
	}
	// first for non-existing
	err = db.Select(q.Eq("ID", "non-existing")).First(&obj)
	if !database.NotFound(err) {
		t.Error("expected not found err, got", err)
	}
	// get for non-existing
	err = db.Get("myStruct", "non-existing", &obj)
	if !database.NotFound(err) {
		t.Error("expected not found err, got", err)
	}
}

func testAdvancedFetching(t *testing.T, db db.DB) {
	initDB(t, db)
	defer db.Close()

	{
		// one by id
		var to myStruct
		err := db.One("ID", "2", &to)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(to, objs[2]) {
			t.Errorf("expected obj '%v' got '%v'", objs[2], to)
		}

		// one by int val
		err = db.One("Val", 33, &to)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(to, objs[3]) {
			t.Errorf("expected obj '%v' got '%v'", objs[3], to)
		}
	}

	// select
	{
		var toSlice []myStruct
		err := db.Select(q.In("ID", []string{"2", "1", "3"})).Find(&toSlice)
		if err != nil {
			t.Error(err)
		}
		if len(toSlice) != len(insertData) {
			t.Errorf("Expected len %v, got %v", len(toSlice), len(insertData))
		}
		for _, o := range insertData {
			objInSlice(t, o, toSlice)
		}
	}

	// advanced selects
	{
		// with orderBy
		var to myStruct
		err := db.Select(q.Gt("Val", 13)).OrderBy("Val").First(&to)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(to, objs[2]) {
			t.Errorf("got obj '%v' expected '%v'", to, objs[2])
		}

		// with orderBy reversed
		err = db.Select(q.Gt("Val", 13)).OrderBy("Val").Reverse().First(&to)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(to, objs[3]) {
			t.Errorf("got obj '%v' expected '%v'", to, objs[3])
		}

		// find with limit 3 reversed
		var toSlice []myStruct
		err = db.Select(q.Gt("Val", 13)).OrderBy("Val").Reverse().Limit(3).Find(&toSlice)
		if err != nil {
			t.Error(err)
		}
		expected := []myStruct{objs[3], objs[2]}
		if !reflect.DeepEqual(toSlice, expected) {
			t.Errorf("got objs '%v' expected '%v'", toSlice, expected)
		}

		// find with limit 1 not reversed
		err = db.Select(q.Gt("Val", 13)).OrderBy("Val").Limit(1).Find(&toSlice)
		if err != nil {
			t.Error(err)
		}
		expected = []myStruct{objs[2]}
		if !reflect.DeepEqual(toSlice, expected) {
			t.Errorf("got objs '%v' expected '%v'", toSlice, expected)
		}

		// find with limit and skip
		err = db.Select(q.Gt("Val", 5)).OrderBy("Val").Limit(1).Skip(2).Find(&toSlice)
		if err != nil {
			t.Error(err)
		}
		expected = []myStruct{objs[3]}
		if !reflect.DeepEqual(toSlice, expected) {
			t.Errorf("got objs '%v' expected '%v'", toSlice, expected)
		}
	}

	// each
	{
		// orderBy string and int
		var toSlice []myStruct
		err := db.Select().OrderBy("Name", "Val").Each(new(myStruct), func(v interface{}) error {
			toSlice = append(toSlice, *v.(*myStruct))
			return nil
		})
		if err != nil {
			t.Error(err)
		}
		expected := []myStruct{objs[2], objs[3], objs[1]}
		if !reflect.DeepEqual(toSlice, expected) {
			t.Errorf("got objs '%v' expected '%v'", toSlice, expected)
		}

		// orderBy string and int, reversed
		toSlice = toSlice[:0]
		err = db.Select().OrderBy("Name", "Val").Reverse().Each(new(myStruct), func(v interface{}) error {
			toSlice = append(toSlice, *v.(*myStruct))
			return nil
		})
		if err != nil {
			t.Error(err)
		}
		expected = []myStruct{objs[1], objs[3], objs[2]}
		if !reflect.DeepEqual(toSlice, expected) {
			t.Errorf("got objs '%v' expected '%v'", toSlice, expected)
		}
	}

	deleteData(t, db)
}

func testTransactions(t *testing.T, db db.DB) {
	initDB(t, db)
	defer db.Close()

	// get one
	{
		tx, err := db.Begin(true)
		if err != nil {
			t.Error(err)
		}

		var to myStruct
		err = tx.One("ID", "2", &to)
		if err != nil {
			t.Error(err)
		}

		err = tx.Commit()
		if err != nil {
			t.Error(err)
		}
		tx.Rollback() // nop after commit
	}

	// set with rollback
	{
		tx, err := db.Begin(true)
		if err != nil {
			t.Error(err)
		}

		objX := myStruct{ID: "2", Val: 777}
		err = tx.Save(&objX)
		if err != nil {
			t.Error(err)
		}

		objX = myStruct{ID: "2", Val: 777}
		err = tx.Save(&objX)
		if err != nil {
			t.Error(err)
		}

		objY := myStruct{}
		err = tx.Get("myStruct", "2", &objY)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(objX, objY) {
			t.Errorf("expected obj '%v' got '%v'", objY, objX)
		}

		err = tx.Rollback()
		if err != nil {
			t.Error(err)
		}

		objY = myStruct{}
		err = db.Get("myStruct", "2", &objY)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(objs[2], objY) {
			t.Errorf("expected obj '%v' got '%v'", objs[2], objY)
		}

	}

	deleteData(t, db)
}

func initDB(t *testing.T, db db.DB) {
	// explicit index building
	{
		err := db.Init(&myStruct{})
		if err != nil {
			t.Error(err)
		}
		db.ReIndex(&myStruct{}) // nop
	}
	// add initial data
	for _, o := range insertData {
		err := db.Save(&o)
		if err != nil {
			t.Error(err)
		}
	}
}

func deleteData(t *testing.T, db db.DB) {
	for _, o := range insertData {
		err := db.DeleteStruct(&o)
		if err != nil {
			t.Error(err)
		}
	}
}

func objInSlice(t *testing.T, o myStruct, s []myStruct) {
	for _, o2 := range s {
		if reflect.DeepEqual(o, o2) {
			return
		}
	}
	t.Errorf("object '%v' not found in slice '%v'", o, s)
}

func objPtrInSlice(t *testing.T, o *myStruct, s []*myStruct) {
	for _, o2 := range s {
		if reflect.DeepEqual(o, o2) {
			return
		}
	}
	t.Errorf("object ptr '%v' not found in slice '%v'", o, s)
}
