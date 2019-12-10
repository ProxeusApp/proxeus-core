// +build integration

package eio

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var l sync.Mutex
var success int
var errCount int

func TestDocumentServiceClient_Compile(t *testing.T) {
	ds := &DocumentServiceClient{Url: "http://localhost:2115/"}
	start := time.Now()
	for c := 0; c < 80; c++ {
		fire := rand.Intn(2000-100) + 100
		log.Println("firing", fire)
		for i := 0; i < fire; i++ {
			go func() {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(60-15)+15))
				call(ds)
			}()
		}
		log.Println("sleeping a bit")
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100-50)+50))

		//rand.Intn(max - min) + min
	}
	end := time.Now()
	time.Sleep(40 * time.Second)
	log.Println("errors:", errCount, "successes:", success, "in:", start.Sub(end)*time.Millisecond, "ms")
}

func call(ds *DocumentServiceClient) {
	_, err := ds.Compile(Template{
		TemplatePath: "/home/ave/Downloads/Subworkflow2.odt",
	})
	if err != nil {
		l.Lock()
		errCount++
		l.Unlock()
		log.Println("error", err)
	} else {
		l.Lock()
		success++
		l.Unlock()
		log.Println("success")
	}
}
