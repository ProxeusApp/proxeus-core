package externalnode

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"

	"io/ioutil"
	"net/http"

	"fmt"
	"time"
)

type (
	ExternalNode struct {
		ID     string `json:"id" storm:"id"`
		Name   string `json:"name"`
		Detail string `json:"detail"`
		Url    string `json:"url"`
		Secret string `json:"secret"`
	}

	ExternalNodeInstance struct {
		ID       string      `json:"id" storm:"id"`
		NodeName string      `json:"nodeName"`
		Config   map[string]interface{} `json:"nodeConfig"`
	}

	ExternalQuery struct {
		*ExternalNode
		*ExternalNodeInstance
	}
)

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "I'm ok")
}

func Register(proxeusUrl, name, serviceUrl, jwtSecret, description string, retryInterval int) error {
	client := http.Client{Timeout: 5 * time.Second}
	var err error
	for {
		n := ExternalNode{
			ID:     name,
			Name:   name,
			Detail: description,
			Url:    serviceUrl,
			Secret: jwtSecret,
		}
		buf, err := json.Marshal(n)
		if err != nil {
			panic(err.Error())
		}
		r, err := client.Post(proxeusUrl+"/api/admin/external/register",
			"application/json", bytes.NewBuffer(buf))
		if err == nil && r.StatusCode == http.StatusOK {
			log.Print("[nodeservice] ", n.Name, " registered")
			return nil
		}

		log.Print("[nodeservice] error registering ", n.Name, " err ", err)
		time.Sleep(time.Duration(retryInterval) * time.Second)
	}
	return err
}

func Nop(_ echo.Context) error {
	return nil
}

func NodeID(c echo.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", errors.New("empty id")
	}
	t := c.Get("user").(*jwt.Token)
	if id != t.Claims.(jwt.MapClaims)["jti"].(string) {
		return "", errors.New("id mismatch")
	}
	return id, nil
}

func SetStoredConfig(c echo.Context, proxeusUrl string, conf interface{}) error {
	id, err := NodeID(c)
	if err != nil {
		return err
	}

	n := ExternalNodeInstance{
		ID:     id,
		Config: conf.(map[string]interface{}),
	}
	buf, err := json.Marshal(n)
	if err != nil {
		log.Print("[nodeservice] error marshalling config ", id, conf, " err ", err)
		return err
	}
	client := http.Client{Timeout: 5 * time.Second}
	_, err = client.Post(proxeusUrl+"/api/admin/external/config/"+n.ID,
		"application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Print("[nodeservice] error updating config ", id, " err ", err)
		return err
	}

	log.Print("[nodeservice] ", id, " config updated")
	return nil
}

func GetStoredConfig(c echo.Context, proxeusUrl string) ([]byte, error) {
	id, err := NodeID(c)
	if err != nil {
		return nil, err
	}
	client := http.Client{Timeout: 5 * time.Second}
	r, err := client.Get(proxeusUrl + "/api/admin/external/config/" + id)
	if err == nil && r.StatusCode == http.StatusOK {
		jsonBody, _ := ioutil.ReadAll(r.Body)
		return jsonBody, nil
	}
	return nil, errors.New("no valid config received")
}

func (e *ExternalNode) HealthUrl() string {
	return fmt.Sprintf("%s/health", e.Url)
}

func (e ExternalQuery) jwtToken() string {
	claims := jwt.StandardClaims{
		Id:        e.ExternalNodeInstance.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(e.Secret))
	if err != nil {
		log.Println("Could not sign token with secret", token.Raw)
		return ""
	}
	return t
}

func (e ExternalQuery) nodeUrl(method string) string {
	return fmt.Sprintf("%s/node/%s/%s?auth=%s",
		e.Url,
		e.ExternalNodeInstance.ID,
		method,
		e.jwtToken(),
	)
}

func (e ExternalQuery) ConfigUrl() string {
	return e.nodeUrl("config")
}

func (e ExternalQuery) NextUrl() string {
	return e.nodeUrl("next")
}

func (e ExternalQuery) RemoveUrl() string {
	return e.nodeUrl("remove")
}

func (e ExternalQuery) CloseUrl() string {
	return e.nodeUrl("close")
}
