package tests

import (
	"fmt"
	"github.com/saichler/reflect/go/reflect/clone"
	"github.com/saichler/reflect/go/reflect/inspect"
	"github.com/saichler/reflect/go/reflect/updater"
	"github.com/saichler/reflect/go/tests/utils"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/tests"
	"testing"
)

func TestUpdater(t *testing.T) {
	in := inspect.NewIntrospect(registry.NewRegistry())
	_, err := in.Inspect(&tests.TestProto{})
	if err != nil {
		log.Fail(t, err.Error())
		return
	}
	upd := updater.NewUpdater(in, false)
	aside := utils.CreateTestModelInstance(0)
	zside := &tests.TestProto{MyString: "updated"}
	uside := in.Clone(aside).(*tests.TestProto)
	err = upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	changes := upd.Changes()

	if len(changes) != 1 {
		t.Fail()
		fmt.Println("Expected 1 change but got ", len(upd.Changes()))
		for _, c := range changes {
			fmt.Println(c.String())
		}
		return
	}

	if aside.MyString != zside.MyString {
		t.Fail()
		fmt.Println("1 Expected ", zside.MyString, " got ", aside.MyString)
		return
	}

	for _, change := range changes {
		change.Apply(uside)
	}

	if uside.MyString != aside.MyString {
		fmt.Println("2 Expected ", aside.MyString, " got ", uside.MyString)
		t.Fail()
		return
	}
}

func TestEnum(t *testing.T) {
	in := inspect.NewIntrospect(registry.NewRegistry())
	_, err := in.Inspect(&tests.TestProto{})
	if err != nil {
		log.Fail(t, err.Error())
		return
	}
	upd := updater.NewUpdater(in, false)
	aside := utils.CreateTestModelInstance(0)
	zside := clone.NewCloner().Clone(aside).(*tests.TestProto)
	zside.MyEnum = tests.TestEnum_ValueTwo

	err = upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}
	if aside.MyEnum != zside.MyEnum {
		log.Fail(t, aside.MyEnum)
		return
	}
}
