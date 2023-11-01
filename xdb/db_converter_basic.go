package xdb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type basicDBConverter struct {
}

func NewBasicDBConverter() DBConverter {
	return &basicDBConverter{}
}

func (b basicDBConverter) ToDBValue(
	val interface{}, hint *DBHint,
) (dbValue string, err error) {
	// TODO convert for binary, datetime, etc using hint

	data, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	dbVal := strings.Trim(string(data), "\"")
	return dbVal, nil
}

func (b basicDBConverter) FromDBValue(
	dbQueryValue string, hint *DBHint, resultPtr interface{},
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
