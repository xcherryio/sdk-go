package xdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertPrimitiveGlobalAttribute(t *testing.T) {
	str := "hello"
	dbVal, err := GetDefaultObjectEncoder().FromGlobalAttributeToDbValue(str, nil)
	assert.Nil(t, err)
	assert.Equal(t, "hello", dbVal)
	str2 := ""
	err = GetDefaultObjectEncoder().FromDbValueToGlobalAttribute(dbVal, nil, &str2)
	assert.Nil(t, err)
	assert.Equal(t, str, str2)

	integer := 1234
	dbVal, err = GetDefaultObjectEncoder().FromGlobalAttributeToDbValue(integer, nil)
	assert.Nil(t, err)
	assert.Equal(t, "1234", dbVal)
	var integer2 int
	err = GetDefaultObjectEncoder().FromDbValueToGlobalAttribute(dbVal, nil, &integer2)
	assert.Nil(t, err)
	assert.Equal(t, integer, integer2)

	boolean := true
	dbVal, err = GetDefaultObjectEncoder().FromGlobalAttributeToDbValue(boolean, nil)
	assert.Nil(t, err)
	assert.Equal(t, "true", dbVal)
	var boolean2 bool
	err = GetDefaultObjectEncoder().FromDbValueToGlobalAttribute(dbVal, nil, &boolean2)
	assert.Nil(t, err)
	assert.Equal(t, boolean, boolean2)

	float := 3.14
	dbVal, err = GetDefaultObjectEncoder().FromGlobalAttributeToDbValue(float, nil)
	assert.Nil(t, err)
	assert.Equal(t, "3.14", dbVal)
	var float2 float32
	err = GetDefaultObjectEncoder().FromDbValueToGlobalAttribute(dbVal, nil, &float2)
	assert.Nil(t, err)
	assert.Equal(t, float2, float2)
}

func TestConvertNilGlobalAttribute(t *testing.T) {
	dbVal, err := GetDefaultObjectEncoder().FromGlobalAttributeToDbValue(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "null", dbVal)

	var obj interface{}
	err = GetDefaultObjectEncoder().FromDbValueToGlobalAttribute(dbVal, nil, &obj)
	assert.Nil(t, err)
	assert.Equal(t, nil, obj)
}

func TestConvertMapGlobalAttribute(t *testing.T) {
	m := map[string]string{
		"a": "b",
		"c": "d",
	}
	dbVal, err := GetDefaultObjectEncoder().FromGlobalAttributeToDbValue(m, nil)
	assert.Nil(t, err)
	assert.Equal(t, "{\"a\":\"b\",\"c\":\"d\"}", dbVal)

	var m2 map[string]string
	err = GetDefaultObjectEncoder().FromDbValueToGlobalAttribute(dbVal, nil, &m2)
	assert.Nil(t, err)
	assert.Equal(t, m, m2)
}

func TestConvertArrayGlobalAttribute(t *testing.T) {
	m := []string{
		"a", "b", "c", "d",
	}
	dbVal, err := GetDefaultObjectEncoder().FromGlobalAttributeToDbValue(m, nil)
	assert.Nil(t, err)
	assert.Equal(t, "[\"a\",\"b\",\"c\",\"d\"]", dbVal)

	var m2 []string
	err = GetDefaultObjectEncoder().FromDbValueToGlobalAttribute(dbVal, nil, &m2)
	assert.Nil(t, err)
	assert.Equal(t, m, m2)
}
