package xc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertPrimitiveGlobalAttribute(t *testing.T) {
	str := "hello"
	dbVal, err := NewBasicDBConverter().ToDBValue(str, nil)
	assert.Nil(t, err)
	assert.Equal(t, "hello", dbVal)
	str2 := ""
	err = NewBasicDBConverter().FromDBValue(dbVal, nil, &str2)
	assert.Nil(t, err)
	assert.Equal(t, str, str2)

	integer := 1234
	dbVal, err = NewBasicDBConverter().ToDBValue(integer, nil)
	assert.Nil(t, err)
	assert.Equal(t, "1234", dbVal)
	var integer2 int
	err = NewBasicDBConverter().FromDBValue(dbVal, nil, &integer2)
	assert.Nil(t, err)
	assert.Equal(t, integer, integer2)

	boolean := true
	dbVal, err = NewBasicDBConverter().ToDBValue(boolean, nil)
	assert.Nil(t, err)
	assert.Equal(t, "true", dbVal)
	var boolean2 bool
	err = NewBasicDBConverter().FromDBValue(dbVal, nil, &boolean2)
	assert.Nil(t, err)
	assert.Equal(t, boolean, boolean2)

	float := 3.14
	dbVal, err = NewBasicDBConverter().ToDBValue(float, nil)
	assert.Nil(t, err)
	assert.Equal(t, "3.14", dbVal)
	var float2 float32
	err = NewBasicDBConverter().FromDBValue(dbVal, nil, &float2)
	assert.Nil(t, err)
	assert.Equal(t, float2, float2)
}

func TestConvertNilGlobalAttribute(t *testing.T) {
	dbVal, err := NewBasicDBConverter().ToDBValue(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "null", dbVal)

	var obj interface{}
	err = NewBasicDBConverter().FromDBValue(dbVal, nil, &obj)
	assert.Nil(t, err)
	assert.Equal(t, nil, obj)
}

func TestConvertMapGlobalAttribute(t *testing.T) {
	m := map[string]string{
		"a": "b",
		"c": "d",
	}
	dbVal, err := NewBasicDBConverter().ToDBValue(m, nil)
	assert.Nil(t, err)
	assert.Equal(t, "{\"a\":\"b\",\"c\":\"d\"}", dbVal)

	var m2 map[string]string
	err = NewBasicDBConverter().FromDBValue(dbVal, nil, &m2)
	assert.Nil(t, err)
	assert.Equal(t, m, m2)
}

func TestConvertArrayGlobalAttribute(t *testing.T) {
	m := []string{
		"a", "b", "c", "d",
	}
	dbVal, err := NewBasicDBConverter().ToDBValue(m, nil)
	assert.Nil(t, err)
	assert.Equal(t, "[\"a\",\"b\",\"c\",\"d\"]", dbVal)

	var m2 []string
	err = NewBasicDBConverter().FromDBValue(dbVal, nil, &m2)
	assert.Nil(t, err)
	assert.Equal(t, m, m2)
}
