package database

import (
	"bytes"
	"testing"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/model"

	. "github.com/onsi/gomega"
)

func TestUser(t *testing.T) {
	RegisterTestingT(t)
	us := testDBSet.User

	options := storage.IndexOptions(0)
	const email = "a@example.com"
	item := &model.User{
		ID:           dummy.UserID(),
		Name:         "some user",
		EthereumAddr: "0xdeee",
	}
	item2 := &model.User{
		ID:           dummySuperAdmin.UserID(),
		Name:         "some admin",
		EthereumAddr: "0xaaaaa",
	}

	Expect(us.GetBaseFilePath()).ToNot(BeEmpty())

	// put
	Expect(us.Put(dummy, item)).To(Succeed())
	Expect(us.Put(dummy, item2)).To(Succeed())
	Expect(us.UpdateEmail(item.ID, email)).To(Succeed())
	item.Email = email
	Expect(us.PutPw(item.ID, "pass")).To(Succeed())

	// get
	Expect(us.Login(email, "pass")).To(equalJSON(item))
	Expect(us.Get(dummy, item.ID)).To(equalJSON(item))
	Expect(us.GetByBCAddress(item.EthereumAddr)).To(equalJSON(item))
	Expect(us.GetByEmail(email)).To(equalJSON(item))
	Expect(us.Count()).To(Equal(2))
	Expect(us.List(dummySuperAdmin, email, options)).To(equalJSON([]*model.User{item}))
	Expect(us.List(dummy, email, options)).To(equalJSON([]*model.User{item}))

	// profile photo
	phData := []byte("img_data")
	buf := bytes.NewBuffer(phData)
	Expect(us.PutProfilePhoto(dummySuperAdmin, item.ID, buf)).To(Succeed())
	buf.Reset()
	Expect(us.GetProfilePhoto(dummySuperAdmin, item.ID, buf)).To(Succeed())
	Expect(buf.Bytes()).To(Equal(phData))
	_, err := us.List(dummySuperAdmin, "", options)
	Expect(err).To(Succeed())

	// api keys
	const keyName = "test key"
	val, err := us.CreateApiKey(dummySuperAdmin, item2.ID, keyName)
	Expect(err).To(Succeed())
	Expect(val).ToNot(Equal(keyName))
	testKeyItem, err := us.GetByApiKey(val, item2.ID)
	Expect(err).To(Succeed())

	// api key without user id provided
	testKeyItem2, err := us.GetByApiKey(val, "")
	Expect(err).To(Succeed())
	Expect(testKeyItem2.ID).To(Equal(item2.ID))
	item2.ApiKeys = []*model.ApiKey{{Name: keyName, Key: testKeyItem.ApiKeys[0].Key}}
	Expect(testKeyItem).To(equalJSON(item2))

	// delete api key
	Expect(us.DeleteApiKey(dummySuperAdmin, item2.ID, item2.ApiKeys[0].Key)).To(Succeed())
}
