package xc

import (
	"encoding/json"
	"github.com/xcherryio/apis/goapi/xcapi"
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

func (b *builtinJsonEncoder) Encode(obj interface{}) (*xcapi.EncodedObject, error) {
	if obj == nil {
		return &xcapi.EncodedObject{}, nil
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return &xcapi.EncodedObject{
		Encoding: encodingType,
		Data:     string(data),
	}, nil
}

func (b *builtinJsonEncoder) Decode(encodedObj *xcapi.EncodedObject, resultPtr interface{}) error {
	if encodedObj == nil || resultPtr == nil || encodedObj.GetData() == "" {
		return nil
	}
	return json.Unmarshal([]byte(encodedObj.GetData()), resultPtr)
}
