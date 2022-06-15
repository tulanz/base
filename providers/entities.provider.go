package providers

import "github.com/tulanz/base/datasource"

type EntityMap struct {
}

func (m *EntityMap) GetEntities() []interface{} {
	return []interface{}{}
}

func NewEntityMap() datasource.EntityMap {
	return &EntityMap{}
}
