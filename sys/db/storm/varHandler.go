package storm

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type VarsMaintenance struct {
	ID   string `json:"id" storm:"id"` //form id
	Vars []string
}

type Var struct {
	ID      string          `storm:"id"` //var
	VarRefs map[string]bool //ids that contain this var
}

func initVars(db *storm.DB) {
	db.Init(&VarsMaintenance{})
	db.Init(&Var{})
}

func updateVarsOf(auth model.Authorization, id string, allNewVars []string, tx storm.Node) error {
	var oldForm VarsMaintenance
	err := tx.One("ID", id, &oldForm)
	if err == storm.ErrNotFound {
		//add vars
		err = nil
		oldForm = VarsMaintenance{}
		oldForm.Vars = allNewVars
		oldForm.ID = id
		err = tx.Save(&oldForm)
		if err != nil {
			return err
		}
		for _, v := range allNewVars {
			err = putVar(id, v, tx)
			if err != nil {
				return err
			}
		}
	}
	if err != nil {
		return err
	}
	//update vars
	varsMap := make(map[string]bool)
	for _, v := range oldForm.Vars {
		varsMap[v] = false
	}
	for _, v := range allNewVars {
		varsMap[v] = true
	}
	//clean up all vars mapped with false
	for v, keep := range varsMap {
		if !keep {
			remVar(id, v, tx)
		}
	}
	for _, v := range allNewVars {
		err = putVar(id, v, tx)
		if err != nil {
			return err
		}
	}
	//update var maintenance
	oldForm.Vars = allNewVars
	return tx.Save(&oldForm)
}

func remVars(auth model.Authorization, id string, tx storm.Node) error {
	var oldForm VarsMaintenance
	err := tx.One("ID", id, &oldForm)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil
		}
		return err
	}
	oldFormRef := &oldForm
	for _, v := range oldFormRef.Vars {
		remVar(id, v, tx)
	}
	return tx.DeleteStruct(oldFormRef)
}

func putVar(id, strVar string, tx storm.Node) error {
	var svar Var
	err := tx.One("ID", strVar, &svar)
	svarRef := &svar
	if err == storm.ErrNotFound {
		err = tx.Save(&Var{ID: strVar, VarRefs: map[string]bool{id: true}})
		if err != nil {
			return err
		}
	} else {
		svarRef.VarRefs[id] = true
		err = tx.Save(svarRef)
		if err != nil {
			return err
		}
	}
	return nil
}

func remVar(id, strVar string, tx storm.Node) error {
	var svar Var
	err := tx.One("ID", strVar, &svar)
	svarRef := &svar
	if err == storm.ErrNotFound {
		return nil
	} else {
		delete(svarRef.VarRefs, id)
		if len(svarRef.VarRefs) == 0 {
			err = tx.DeleteStruct(svarRef)
		} else {
			err = tx.Save(svarRef)
		}
		return err
	}
}

func getVars(contains string, limit, index int, tx storm.Node) ([]string, error) {
	items := make([]string, 0)
	err := tx.Select(
		q.Re("ID", contains)).
		Limit(limit).
		Skip(index).
		OrderBy("ID").
		Each(new(Var), func(record interface{}) error {
			item := record.(*Var)
			items = append(items, item.ID)
			return nil
		})
	return items, err
}
