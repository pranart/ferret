package values

import (
	"crypto/sha512"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	ObjectPredicate = func(value core.Value, key string) bool
	ObjectProperty  struct {
		name  string
		value core.Value
	}
	Object struct {
		value map[string]core.Value
	}
)

func NewObjectProperty(name string, value core.Value) *ObjectProperty {
	return &ObjectProperty{name, value}
}

func NewObject() *Object {
	return &Object{make(map[string]core.Value)}
}

func NewObjectWith(props ...*ObjectProperty) *Object {
	obj := NewObject()

	for _, prop := range props {
		obj.value[prop.name] = prop.value
	}

	return obj
}

func (t *Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func (t *Object) Type() core.Type {
	return core.ObjectType
}

func (t *Object) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "{}"
	}

	return string(marshaled)
}

func (t *Object) Compare(other core.Value) int {
	switch other.Type() {
	case core.ObjectType:
		arr := other.(*Object)

		if t.Length() == 0 && arr.Length() == 0 {
			return 0
		}

		var res = 1

		for _, val := range t.value {
			arr.ForEach(func(otherVal core.Value, key string) bool {
				res = val.Compare(otherVal)

				return res != -1
			})
		}

		return res
	default:
		return 1
	}
}

func (t *Object) Unwrap() interface{} {
	obj := make(map[string]interface{})

	for key, val := range t.value {
		obj[key] = val.Unwrap()
	}

	return obj
}

func (t *Object) Hash() int {
	bytes, err := t.MarshalJSON()

	if err != nil {
		return 0
	}

	h := sha512.New()

	out, err := h.Write(bytes)

	if err != nil {
		return 0
	}

	return out
}

func (t *Object) Clone() core.Value {
	c := NewObject()

	for k, v := range t.value {
		c.Set(NewString(k), v)
	}

	return c
}

func (t *Object) Length() Int {
	return Int(len(t.value))
}

func (t *Object) Keys() []string {
	keys := make([]string, 0, len(t.value))

	for k := range t.value {
		keys = append(keys, k)
	}

	return keys
}

func (t *Object) ForEach(predicate ObjectPredicate) {
	for key, val := range t.value {
		if predicate(val, key) == false {
			break
		}
	}
}

func (t *Object) Get(key String) (core.Value, bool) {
	val, found := t.value[string(key)]

	if found {
		return val, found
	}

	return None, found
}

func (t *Object) GetIn(path []core.Value) (core.Value, error) {
	return GetIn(t, path)
}

func (t *Object) Set(key String, value core.Value) {
	if core.IsNil(value) == false {
		t.value[string(key)] = value
	} else {
		t.value[string(key)] = None
	}
}

func (t *Object) Remove(key String) {
	delete(t.value, string(key))
}

func (t *Object) SetIn(path []core.Value, value core.Value) error {
	return SetIn(t, path, value)
}
