package cloning

import (
	"reflect"
	"strconv"
	"strings"
)

type Cloner struct {
	cloners map[reflect.Kind]func(reflect.Value, string, map[string]reflect.Value) reflect.Value
}

func NewCloner() *Cloner {
	cloner := &Cloner{}
	cloner.initCloners()
	return cloner
}

func (this *Cloner) initCloners() {
	this.cloners = make(map[reflect.Kind]func(reflect.Value, string, map[string]reflect.Value) reflect.Value)
	this.cloners[reflect.Int] = this.intCloner
	this.cloners[reflect.Int32] = this.int32Cloner
	this.cloners[reflect.Ptr] = this.ptrCloner
	this.cloners[reflect.Struct] = this.structCloner
	this.cloners[reflect.String] = this.stringCloner
	this.cloners[reflect.Slice] = this.sliceCloner
	this.cloners[reflect.Map] = this.mapCloner
	this.cloners[reflect.Bool] = this.boolCloner
	this.cloners[reflect.Int64] = this.int64Cloner
	this.cloners[reflect.Uint] = this.uintCloner
	this.cloners[reflect.Uint32] = this.uint32Cloner
	this.cloners[reflect.Uint64] = this.uint64Cloner
	this.cloners[reflect.Float32] = this.float32Cloner
	this.cloners[reflect.Float64] = this.float64Cloner
}

func (this *Cloner) Clone(any interface{}) interface{} {
	value := reflect.ValueOf(any)
	stopLoop := make(map[string]reflect.Value)
	valueClone := this.clone(value, "", stopLoop)
	return valueClone.Interface()
}

func (this *Cloner) clone(value reflect.Value, fieldName string, stopLoop map[string]reflect.Value) reflect.Value {
	if !value.IsValid() {
		return value
	}
	kind := value.Kind()
	cloner := this.cloners[kind]
	if cloner == nil {
		panic("No cloner for kind:" + kind.String() + ":" + fieldName)
	}
	return cloner(value, fieldName, stopLoop)
}

func (this *Cloner) sliceCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}
	newSlice := reflect.MakeSlice(reflect.SliceOf(value.Index(0).Type()), value.Len(), value.Len())
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		elemClone := this.clone(elem, name, stopLoop)
		newSlice.Index(i).Set(elemClone)
	}
	return newSlice
}

func (this *Cloner) ptrCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}

	p := strconv.Itoa(int(value.Pointer()))
	exist, ok := stopLoop[p]
	if ok {
		return exist
	}

	newPtr := reflect.New(value.Elem().Type())
	stopLoop[p] = newPtr

	newPtr.Elem().Set(this.clone(value.Elem(), name, stopLoop))

	return newPtr
}

func (this *Cloner) structCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	cloneStruct := reflect.New(value.Type()).Elem()
	structType := value.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldName := structType.Field(i).Name
		if SkipFieldByName(fieldName) {
			continue
		}
		cloned := this.clone(fieldValue, structType.Field(i).Name, stopLoop)
		if cloned.Kind() == reflect.Int32 {
			cloneStruct.Field(i).SetInt(cloned.Int())
		} else {
			cloneStruct.Field(i).Set(cloned)
		}
	}
	return cloneStruct
}

func (this *Cloner) mapCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}
	mapKeys := value.MapKeys()
	mapClone := reflect.MakeMapWithSize(value.Type(), len(mapKeys))
	for _, key := range mapKeys {
		mapElem := value.MapIndex(key)
		mapElemClone := this.clone(mapElem, name, stopLoop)
		mapClone.SetMapIndex(key, mapElemClone)
	}
	return mapClone
}

func (this *Cloner) intCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Int()
	return reflect.ValueOf(int(i))
}

func (this *Cloner) uintCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Uint()
	return reflect.ValueOf(uint(i))
}

func (this *Cloner) uint32Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Uint()
	return reflect.ValueOf(uint32(i))
}

func (this *Cloner) uint64Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Uint()
	return reflect.ValueOf(uint64(i))
}

func (this *Cloner) float32Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Float()
	return reflect.ValueOf(float32(i))
}

func (this *Cloner) float64Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Float()
	return reflect.ValueOf(float64(i))
}

func (this *Cloner) boolCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	b := value.Bool()
	return reflect.ValueOf(b)
}

func (this *Cloner) int32Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Int()
	return reflect.ValueOf(int32(i))
}

func (this *Cloner) int64Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Int()
	return reflect.ValueOf(int64(i))
}

func (this *Cloner) stringCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	s := value.String()
	return reflect.ValueOf(s)
}

func SkipFieldByName(fieldName string) bool {
	if fieldName == "DoNotCompare" {
		return true
	}
	if fieldName == "DoNotCopy" {
		return true
	}
	if len(fieldName) > 3 && fieldName[0:3] == "XXX" {
		return true
	}
	if fieldName[0:1] == strings.ToLower(fieldName[0:1]) {
		return true
	}
	return false
}
