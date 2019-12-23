package storm

import (
	"fmt"
	"testing"

	"github.com/asdine/storm"
)

func TestVars(t *testing.T) {
	db, err := NewUserDataDB("./")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		db.Close()
		db.remove()
	}()
	initVars(db.db)
	pVars(db.db, "123", []string{"var1", "var2", "var3", "var4"})
	pVars(db.db, "123", []string{"var1", "var2", "var5"})
	pVars(db.db, "123", []string{"var1"})
	pVars(db.db, "123", []string{"var1", "var2", "var5"})
	pVars(db.db, "1234", []string{"var1", "var2", "var3", "var4"})
	pVars(db.db, "12356", []string{"var1", "var2", "var3", "var4"})
	pVars(db.db, "12357", []string{"var1", "var2", "var3", "var4"})
	pVars(db.db, "12358", []string{"var1", "var2", "var3", "var4"})
	pVars(db.db, "12359", []string{"aa", "abs", "var3", "abc"})
	rVars(db.db, "1234")
	rVars(db.db, "12356")
	rVars(db.db, "12357")
	rVars(db.db, "12358")
	tx, err := db.db.Begin(false)
	varsRes, err := getVars("", 10, 0, tx)
	if err != nil {
		t.Error(err)
		return
	}
	tx.Rollback()
	if fmt.Sprint(varsRes) != "[aa abc abs var1 var2 var3 var5]" {
		t.Error(varsRes)
	}
}

func pVars(db *storm.DB, id string, vars []string) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	err = updateVarsOf(nil, id, vars, tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
	}
	return err
}

func rVars(db *storm.DB, id string) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	err = remVars(nil, id, tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
	}
	return err
}
