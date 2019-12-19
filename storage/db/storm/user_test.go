package storm

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"

	"github.com/ProxeusApp/proxeus-core/sys/model"

	//"reflect"
	//"git.proxeus.com/core/form"
	"regexp"
	"testing"
	"time"
)

//func TestVal(t *testing.T) {
//	validate = validator.New()
//
//	address := &Address{
//		Street: "Eavesdown Docks",
//		Planet: "Persphone",
//		Phone:  "none",
//	}
//
//	user := &User{
//		FirstName: "",
//		LastName:  "",
//		Age:       135,
//		Email:     "Badger.Smith@gmail.com",
//		//FavouriteColor: "#000-",
//		Addresses: []*Address{address},
//	}
//
//	// returns nil or ValidationErrors ( []FieldError )
//	ve := make(map[string][]form.ValidateError)
//	err := validate.Struct(user)
//	if err != nil {
//		bts, er := json.Marshal(err)
//		if er != nil {
//			panic(er)
//		}
//		log.Println(string(bts))
//		// this check is only needed when your code could produce
//		// an invalid value for validation such as interface with nil
//		// value most including myself do not usually have code like this.
//		if _, ok := err.(*validator.InvalidValidationError); ok {
//			fmt.Println(err)
//			return
//		}
//		log.Println(err)
//		for _, err := range err.(validator.ValidationErrors) {
//			uv := reflect.TypeOf(User{})
//			log.Println(err.Field())
//			fieldName := err.Field()
//			field, ok := uv.FieldByName(err.Field())
//			log.Println(field, ok)
//			if ok {
//				v := field.Tag.Get("json")
//				if v != "" {
//					fieldName = v
//				}
//			}
//			var arr []form.ValidateError
//			if a, b := ve[fieldName]; b {
//				arr = a
//			} else {
//				arr = make([]form.ValidateError, 0)
//				ve[fieldName] = arr
//			}
//			log.Println("append", err.Tag())
//			arr = append(arr, form.ValidateError{I: 0, Msg: err.Tag()})
//			ve[fieldName] = arr
//			log.Println(err)
//			fmt.Println(err.Namespace())
//			fmt.Println(err.Field())
//			fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
//			fmt.Println(err.StructField())     // by passing alt name to ReportError like below
//			fmt.Println(err.Tag())
//			fmt.Println(err.ActualTag())
//			fmt.Println(err.Kind())
//			fmt.Println(err.Type())
//			fmt.Println(err.Value())
//			fmt.Println(err.Param())
//			fmt.Println()
//		}
//		log.Println(ve)
//		// from here you can create your own error messages in whatever language you wish
//		return
//	}
//
//}

func TestUserDB_List(t *testing.T) {
	udb, err := NewUserDB("")
	if err != nil {
		log.Println(err)
	}
	orgID := "1"
	orgEmail := "artan.veliju@gmail.com"
	u := model.User{}
	u.Name = "Artan Veliju"
	insert(udb, err, orgID, orgEmail, u)
	orgID = "n1"
	orgEmail = "nadije.veliju@gmail.com"
	u = model.User{}
	u.Name = "Nadije Veliju"
	insert(udb, err, orgID, orgEmail, u)
	udb.db.Close()
}

func insert(udb *UserDB, err error, orgID, orgEmail string, u model.User) {
	i := 0
	max := 10000

	var tx storm.Node

	txSet := false
	defer func() {
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			log.Println(err, u)
		}
	}()
	for {
		if !txSet {
			tx, err = udb.db.Begin(true)
			if err != nil {
				log.Println(err, u)
				break
			}
			txSet = true
		}
		u.ID = fmt.Sprintf("%s%v", orgID, i)
		u.Email = fmt.Sprintf("%v%s", i, orgEmail)
		u.Created = time.Now()
		u.Updated = time.Now()
		err = tx.Save(&u)
		if err != nil {
			log.Println(err, u)
			break
		}
		if i%1000 == 0 {
			txSet = false
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				log.Println(err, u)
				break
			}
		}
		if i >= max {
			break
		}
		i++
		log.Println("after put", i)
	}
	tx.Commit()
}

func TestUserDB_List1(t *testing.T) {
	udb, err := NewUserDB("")
	if err != nil {
		log.Println(err)
	}
	var users []model.User
	udb.db.All(&users)
	//for i, u := range users {
	//	log.Println(i, u)
	//}
	udb.db.Close()
}

func TestNewUserDB(t *testing.T) {
	a := base64.StdEncoding.EncodeToString([]byte("abcksjlkfjslkdjflksdjf"))
	a = "data:image/jpeg;base64," + a
	log.Println(a)
	abts := []byte(a)
	beginHere := []byte(";base64,")
	abts = abts[bytes.Index(abts, beginHere)+len(beginHere):]
	log.Println(string(abts))
	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(abts)))
	n, err := base64.StdEncoding.Decode(dbuf, abts)
	if err != nil {
		log.Println(err)
	}
	abts = dbuf[:n]
	log.Println(string(abts))
}

func TestUserDB_List2(t *testing.T) {
	udb, err := NewUserDB("")
	if err != nil {
		log.Println(err)
	}

	//udb.db.All(&users)
	em := regexp.QuoteMeta("vel")
	index := 0
	for {
		var users []*model.User
		limit := 1000
		skip := limit * index
		log.Println("skip", skip)
		err = udb.db.Select(q.Re("Email", em)).Limit(limit).Skip(skip).OrderBy("Updated").Reverse().Find(&users)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("-------------------------------------------------------------")
		for i, u := range users {
			log.Println(i, u)
		}
		log.Println("-------------------------------------------------------------")
		log.Println("-------------------------------------------------------------")
		index++
	}

}
