package xdb

import (
	"encoding/json"
	"fmt"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"reflect"
	"strings"
)

func GetDefaultObjectEncoder() ObjectEncoder {
	return &builtinJsonEncoder{}
}

type builtinJsonEncoder struct {
}

func (b *builtinJsonEncoder) FromGlobalAttributeToDbValue(
	val interface{}, hint *string,
) (string, error) {
	// TODO convert for binary, datetime, etc using hint

	data, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	dbVal := strings.Trim(string(data), "\"")
	return dbVal, nil
}

func (b *builtinJsonEncoder) FromDbValueToGlobalAttribute(
	dbQueryValue string, hint *string, resultPtr interface{},
) error {
	if dbQueryValue == "" {
		return nil
	}
	rv := reflect.TypeOf(resultPtr)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("resultPtr must be a pointer")
	}
	ele := rv.Elem()
	if ele.Kind() == reflect.String {
		// append "" to make it a valid string json
		dbQueryValue = "\"" + dbQueryValue + "\""
	}
	return json.Unmarshal([]byte(dbQueryValue), resultPtr)
}

const encodingType = "golangJson"

func (b *builtinJsonEncoder) GetEncodingType() string {
	return encodingType
}

func (b *builtinJsonEncoder) Encode(obj interface{}) (*xdbapi.EncodedObject, error) {
	if obj == nil {
		return &xdbapi.EncodedObject{}, nil
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return &xdbapi.EncodedObject{
		Encoding: encodingType,
		Data:     string(data),
	}, nil
}

func (b *builtinJsonEncoder) Decode(encodedObj *xdbapi.EncodedObject, resultPtr interface{}) error {
	if encodedObj == nil || resultPtr == nil || encodedObj.GetData() == "" {
		return nil
	}
	return json.Unmarshal([]byte(encodedObj.GetData()), resultPtr)
}
