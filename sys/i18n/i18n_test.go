package i18n

import (
	"encoding/json"
	"log"
	"testing"
)

type A struct {
	A1 int
	A2 int
}

type B struct {
	A
	B1 int
	B2 int
}

type C struct {
	B
	A1 int
	B1 int
	C1 int
	C2 int
	C3 int
}

func (me *A) Me1() {
	log.Println("from A")
}
func (me *B) Me() {
	log.Println("from B")
}
func (me *C) Me() {
	log.Println("from C")
}

func TestPutAndGet(t *testing.T) {
	c := &C{}
	bts, err := json.Marshal(c)
	log.Println(string(bts), err)
	c.Me()
	c.Me1()
}
