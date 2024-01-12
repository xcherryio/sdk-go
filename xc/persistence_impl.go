package xc

import "github.com/xcherryio/apis/goapi/xcapi"

type persistenceImpl struct {

	// for local attributes
	localAttrKeys         map[string]bool
	currLocalAttrs        map[string]xcapi.EncodedObject
	currUpdatedLocalAttrs map[string]xcapi.EncodedObject
}

func NewPersistenceImpl(
	localAttrKeys map[string]bool,
	currLocalAttrs *xcapi.LoadLocalAttributesResponse,
) Persistence {

	currLocalAttrsMap := map[string]xcapi.EncodedObject{}
	if currLocalAttrs != nil {
		for _, kv := range currLocalAttrs.Attributes {
			if _, ok := localAttrKeys[kv.Key]; !ok {
				panic("local attribute not found " + kv.Key)
			}
			currLocalAttrsMap[kv.Key] = kv.Value
		}
	}

	return &persistenceImpl{
		localAttrKeys:         localAttrKeys,
		currLocalAttrs:        currLocalAttrsMap,
		currUpdatedLocalAttrs: map[string]xcapi.EncodedObject{},
	}
}

func (p *persistenceImpl) GetLocalAttribute(key string, resultPtr interface{}) {
	_, ok := p.localAttrKeys[key]
	if !ok {
		panic("local attribute not found " + key)
	}

	curVal, ok := p.currLocalAttrs[key]
	if !ok {
		return
	}

	err := GetDefaultObjectEncoder().Decode(&curVal, resultPtr)
	if err != nil {
		panic(err)
	}
}

func (p *persistenceImpl) SetLocalAttribute(key string, value interface{}) {
	_, ok := p.localAttrKeys[key]
	if !ok {
		panic("local attribute is not defined/registered in the PersistenceSchema: " + key)
	}

	encodedVal, err := GetDefaultObjectEncoder().Encode(value)
	if err != nil {
		panic(err)
	}

	p.currLocalAttrs[key] = *encodedVal
	p.currUpdatedLocalAttrs[key] = *encodedVal
}

func (p *persistenceImpl) getLocalAttributesToUpdate() []xcapi.KeyValue {
	var res []xcapi.KeyValue
	for k, v := range p.currUpdatedLocalAttrs {
		res = append(res, xcapi.KeyValue{
			Key:   k,
			Value: v,
		})
	}
	return res
}
