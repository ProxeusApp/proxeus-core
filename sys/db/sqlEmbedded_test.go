package db

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"testing"

	"github.com/cznic/ql"
)

var (
	compiledCommit        = ql.MustCompile("COMMIT;")
	compiledBeginTx       = ql.MustCompile("BEGIN TRANSACTION;")
	compiledCreate        = ql.MustCompile("CREATE TABLE IF NOT EXISTS t (i16 int16, s16 string, s string);")
	compiledCreateIndex   = ql.MustCompile("CREATE INDEX IF NOT EXISTS xi16 ON t (i16);CREATE INDEX IF NOT EXISTS xs16 ON t (s16);CREATE INDEX IF NOT EXISTS xs ON t (s);")
	compiledCreate2       = ql.MustCompile("BEGIN TRANSACTION; CREATE TABLE t (i16 int16, s16 string, s string); COMMIT;")
	compiledIns           = ql.MustCompile("INSERT INTO t VALUES($1, $2, $3);")
	compiledSelect        = ql.MustCompile("SELECT * FROM t;")
	compiledSelectOrderBy = ql.MustCompile("SELECT * FROM t ORDER BY i16, s16;")
	compiledTrunc         = ql.MustCompile("BEGIN TRANSACTION; TRUNCATE TABLE t; COMMIT;")
)

//TODO Optimize sorting performance
func TestSQLFileDB_IO(t *testing.T) {
	filePath := "./embeddedSql/sql.db"
	db, err := NewSQLFileDB(filePath)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	inf, err := db.Info()
	fmt.Println("info", inf)
	countQuery := ql.MustCompile("select count(*) as count from t;")
	ctx := ql.NewRWCtx()
	if _, _, err := db.Execute(ctx, compiledBeginTx); err != nil {
		//t.Error(i, err)
		//return
	}
	if _, _, err := db.Execute(ctx, compiledCreate); err != nil {
		fmt.Println(err)
		//t.Error(i, err)
		//return
	}
	fmt.Println("create index")
	if _, _, err := db.Execute(ctx, compiledCreateIndex); err != nil {
		fmt.Println(err)
		//t.Error(i, err)
		//return
	}
	if a, _, err := db.Execute(ctx, countQuery); err == nil {
		aa, bb := a[0].FirstRow()

		fmt.Println("info 2", aa, bb)
		//t.Error(i, err)
		//return
	}

	rng := rand.New(rand.NewSource(42))
	fmt.Println("inserting")
	for i := 0; i < 100; i++ {
		if _, _, err := db.Execute(ctx, compiledIns, int16(rng.Int()), rnds16(rng, 1), rnds16(rng, 63)); err != nil {
			t.Error(err)
			return
		}
	}
	fmt.Println("commit")
	if _, i, err := db.Execute(ctx, compiledCommit); err != nil {
		t.Error(i, err)
		return
	}

	runtime.GC()
	fmt.Println("selecting")
	for i := 0; i < 10; i++ {
		fmt.Println("execute select")
		rs, index, err := db.Execute(nil, compiledSelectOrderBy)
		if err != nil {
			t.Error(index, err)
			return
		}
		fmt.Println(index, rs)
		flds, err := rs[0].Fields()
		fmt.Println(flds)
		//count := 0
		fmt.Println("do..")
		//if err = rs[0].Do(false, func(record []interface{}) (bool, error) {
		//	if count > 100 {
		//		return false, nil
		//	}
		//	fmt.Println(record)
		//	count++
		//	return true, nil
		//}); err != nil {
		//	t.Errorf("%v %T(%#v)", err, err, err)
		//	return
		//}
		fmt.Println(rs[0].Rows(20, 1))
	}
}

func rnds16(rng *rand.Rand, n int) string {
	a := make([]string, n)
	for i := range a {
		a[i] = fmt.Sprintf("%016x", rng.Int63())
	}
	return strings.Join(a, "")
}
