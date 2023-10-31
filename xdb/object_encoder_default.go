package xdb

import (
	"encoding/json"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
)

func GetDefaultObjectEncoder() ObjectEncoder {
	return &builtinJsonEncoder{}
}

type builtinJsonEncoder struct {
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
