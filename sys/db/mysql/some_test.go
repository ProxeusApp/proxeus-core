package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	d "git.proxeus.com/core/central/sys/db"
	"git.proxeus.com/core/central/sys/model"
)

func TestCompatibilityConvert(t *testing.T) {
	fmt.Println(fmt.Sprintf(""))
}

/**
version
i18n.json
workflow.json
form.json
template.json
*/
func TestNewFormStore(t *testing.T) {
	dbName := "bcdocs"
	dbUser := "root"
	dbPw := "root"
	baseFilePath := ""

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPw, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	usrStore, err := NewUserStore(db, baseFilePath)
	if err != nil {
		panic(err)
	}

	pk := "123"
	usr, err := usrStore.Login("walletSign", "", pk)
	log.Println("end", usr, err)
	if err == d.ErrNotFound {
		it := &model.User{
			EthereumAddr: pk,
			Role:         model.CREATOR,
		}
		err = usrStore.Put(nil, it)
		if err != nil {
			log.Println("put", err)
		}
		usr, err = usrStore.Login("walletSign", "", pk)
		log.Println("end", usr, err)
	}
}
