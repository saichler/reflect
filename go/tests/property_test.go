package tests

import (
	"github.com/saichler/reflect/go/reflect/inspect"
	"github.com/saichler/reflect/go/reflect/property"
	"github.com/saichler/reflect/go/reflect/updater"
	"github.com/saichler/reflect/go/tests/utils"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/tests"
	"github.com/saichler/shared/go/types"
	"testing"
)

var _introspect interfaces.IIntrospector

func propertyOf(id string, root interface{}, t *testing.T) (interface{}, bool) {
	ins, err := property.PropertyOf(id, _introspect)
	if err != nil {
		log.Fail(t, "failed with id: ", id, err.Error())
		return nil, false
	}

	v, err := ins.Get(root)
	if err != nil {
		log.Fail(t, "failed with get: ", id, err.Error())
		return nil, false
	}
	return v, true
}

func TestPrimaryKey(t *testing.T) {
	_introspect = inspect.NewIntrospect(registry.NewRegistry())
	node, err := _introspect.Inspect(&tests.TestProto{})
	if err != nil {
		log.Fail(t, "failed with inspect: ", err.Error())
		return
	}
	_introspect.AddDecorator(types.DecoratorType_Primary, []string{"MyString"}, node)
	aside := utils.CreateTestModelInstance(1)
	zside := utils.CreateTestModelInstance(1)
	zside.MyEnum = tests.TestEnum_ValueTwo

	upd := updater.NewUpdater(_introspect, false)
	err = upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, "failed with update: ", err.Error())
		return
	}
	if len(upd.Changes()) != 1 {
		log.Fail(t, "wrong number of changes: ", len(upd.Changes()))
		return
	}

	pid := upd.Changes()[0].PropertyId()
	n := upd.Changes()[0].NewValue()

	p, e := property.PropertyOf(pid, _introspect)
	if e != nil {
		log.Fail(t, "failed with property: ", e.Error())
		return
	}

	_, root, e := p.Set(nil, n)
	if e != nil {
		log.Fail(t, "failed with set: ", e.Error())
		return
	}

	yside := root.(*tests.TestProto)
	if yside.MyEnum != aside.MyEnum {
		log.Fail(t, "wrong enum: ", yside.MyEnum)
		return
	}
	if yside.MyString != aside.MyString {
		log.Fail(t, "wrong string: ", yside.MyString)
		return
	}
}

func TestInstance(t *testing.T) {
	_introspect = inspect.NewIntrospect(registry.NewRegistry())
	node, err := _introspect.Inspect(&tests.TestProto{})
	if err != nil {
		log.Fail(t, "failed with inspect: ", err.Error())
		return
	}
	_introspect.AddDecorator(types.DecoratorType_Primary, []string{"MyString"}, node)

	id := "testproto<{24}Hello>"
	v, ok := propertyOf(id, nil, t)
	if !ok {
		return
	}

	mytest := v.(*tests.TestProto)
	if mytest.MyString != "Hello" {
		log.Fail(t, "wrong string: ", mytest.MyString)
		return
	}

	mytest.MyFloat64 = 128.128
	id = "testproto.myfloat64"
	v, ok = propertyOf(id, mytest, t)
	if !ok {
		return
	}

	f := v.(float64)
	if f != mytest.MyFloat64 {
		log.Fail(t, "wrong float64: ", f)
		return
	}

	mytest.MySingle = &tests.TestProtoSub{MyString: "Hello"}

	id = "testproto.mysingle.mystring"
	v, ok = propertyOf(id, mytest, t)
	if !ok {
		return
	}
	s := v.(string)
	if s != mytest.MySingle.MyString {
		log.Fail(t, "wrong string: ", s)
		return
	}

	/*
		myInstsnce:=model.MyTestModel{
			MyString: "Hello",
			MySingle: &model.MyTestSubModelSingle{MyString: "World"},
		}

		instance,_:=instance.propertyOf("mytestmodel.mysingle.mystring",introspect.DefaultIntrospect)

		//Getting a value
		v,_:=instance.Get(myInstsnce)
		//Creating another instance
		myOtherInstance:=model.MyTestModel{}
		//Setting the value we fetched from the original instance
		instance.Set(myOtherInstance,"Metadata")

	*/
}
