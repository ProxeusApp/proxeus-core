package compatability

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CarriedJsonRaw struct {
	json.RawMessage
}

type CarriedStringMap map[string]interface{}

func (c CarriedStringMap) MarshalBSON() ([]byte, error) {
	raw, err := json.Marshal(&c)
	if err != nil {
		return nil, err
	}
	v := bson.M{"d": raw}
	return bson.Marshal(v)
}

func (c *CarriedStringMap) UnmarshalBSON(raw []byte) error {
	v := bson.M{"d": primitive.Binary{}}
	err := bson.Unmarshal(raw, &v)
	if err != nil {
		return err
	}
	b := v["d"].(primitive.Binary).Data
	if len(b) == 0 || string(b) == "null" {
		b = []byte{'{', '}'}
	}
	return json.Unmarshal(b, c)
}

func ToMapStringIF(d interface{}) (map[string]interface{}, bool) {
	if dd, ok := d.(CarriedStringMap); ok {
		return dd, true
	}
	dd, ok := d.(map[string]interface{})
	return dd, ok
}
