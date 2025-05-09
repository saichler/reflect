package properties

import (
	"errors"
	"github.com/saichler/reflect/go/reflect/introspecting"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
	"reflect"
)

func (this *Property) Set(any interface{}, value interface{}) (interface{}, interface{}, error) {
	if this == nil {
		return nil, nil, errors.New("property is nil, cannot instantiate")
	}
	if this.parent == nil {
		if any == nil {
			info, err := this.introspector.Registry().Info(this.node.TypeName)
			if err != nil {
				return nil, nil, err
			}
			newAny, err := info.NewInstance()
			if err != nil {
				return nil, nil, err
			}
			any = newAny
		}
		if this.key != nil {
			this.SetPrimaryKey(this.node, any, this.key)
		}
		return any, any, nil
	}
	parent, root, err := this.parent.Set(any, value)
	if err != nil {
		return nil, nil, err
	}
	if any == nil {
		any = root
	}
	parentValue := reflect.ValueOf(parent)
	if parentValue.Kind() == reflect.Ptr {
		parentValue = parentValue.Elem()
	}

	//Special case for setting a value to the map
	if this.node.IsMap && parentValue.Kind() == reflect.Map {
		if this.IsLeaf() {
			parentValue.SetMapIndex(reflect.ValueOf(this.key), reflect.ValueOf(this.value))
		}
		return this.value, any, nil
	} else if parentValue.Kind() == reflect.Map {
		parentValue = parentValue.MapIndex(reflect.ValueOf(this.key))
	}

	myValue := parentValue.FieldByName(this.node.FieldName)
	info, err := this.introspector.Registry().Info(this.node.TypeName)
	if err != nil {
		return nil, nil, err
	}
	typ := info.Type()
	if this.node.IsMap {
		v, e := this.mapSet(myValue, reflect.ValueOf(value))
		return v, any, e
	} else if this.node.IsSlice {
		v, e := this.sliceSet(myValue, reflect.ValueOf(value))
		return v, any, e
	} else if this.introspector.Kind(this.node) == reflect.Struct {
		if !myValue.IsValid() || myValue.IsNil() {
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Ptr &&
				!v.IsNil() && v.Elem().Type().Name() == typ.Name() {
				myValue.Set(reflect.ValueOf(value))
			} else {
				newInstance := reflect.New(typ)
				if v.Kind() == reflect.String {
					serializer := info.Serializer(ifs.STRING)
					if serializer != nil {
						inst, _ := serializer.Unmarshal([]byte(v.String()), this.introspector.Registry())
						if inst != nil {
							newInstance = reflect.ValueOf(inst)
						}
					}
				}
				myValue.Set(newInstance)
			}
		}
		return myValue.Interface(), any, err
	} else if reflect.ValueOf(value).Kind() == reflect.Int32 || myValue.Kind() == reflect.Int32 {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.String {
			value = this.introspector.Registry().Enum(value.(string))
		}
		myValue.SetInt(reflect.ValueOf(value).Int())
		return value, any, err
	} else {
		if value != nil {
			myValue.Set(reflect.ValueOf(value))
		}
		return value, any, err
	}
}

func (this *Property) SetPrimaryKey(node *types.RNode, any interface{}, anyKey interface{}) {
	if anyKey == nil {
		return
	}
	var fieldsValues []interface{}
	if reflect.ValueOf(anyKey).Kind() == reflect.Slice {
		fieldsValues = anyKey.([]interface{})
	} else {
		fieldsValues = []interface{}{anyKey}
	}
	value := reflect.ValueOf(any)
	if !value.IsValid() {
		return
	}
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}

	f := introspecting.PrimaryKeyDecorator(node)
	if f != nil {
		fields := f.([]string)
		for i, attr := range fields {
			fld := value.FieldByName(attr)
			fld.Set(reflect.ValueOf(fieldsValues[i]))
		}
	}
}
