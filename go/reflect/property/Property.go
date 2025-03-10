package property

import (
	"errors"
	"github.com/saichler/reflect/go/reflect/common"
	strings2 "github.com/saichler/shared/go/share/strings"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"strings"
)

type Property struct {
	parent       *Property
	node         *types.RNode
	key          interface{}
	value        interface{}
	id           string
	introspector common.IIntrospector
}

func NewProperty(node *types.RNode, parent *Property, key interface{}, value interface{}, introspector common.IIntrospector) *Property {
	property := &Property{}
	property.parent = parent
	property.node = node
	property.key = key
	property.value = value
	property.introspector = introspector
	return property
}

func PropertyOf(propertyId string, introspector common.IIntrospector) (*Property, error) {
	propertyKey := common.PropertyNodeKey(propertyId)
	node, ok := introspector.Node(propertyKey)
	if !ok {
		return nil, errors.New("Unknown attribute " + propertyKey)
	}
	return newProperty(node, propertyId, introspector)
}

func (this *Property) Parent() *Property {
	return this.parent
}

func (this *Property) Node() *types.RNode {
	return this.node
}

func (this *Property) Key() interface{} {
	return this.key
}

func (this *Property) Value() interface{} {
	return this.value
}

func (this *Property) setKeyValue(propertyId string) (string, error) {
	id := propertyId
	dIndex := strings.LastIndex(propertyId, ".")
	if dIndex == -1 {
		return "", nil
	}
	beIndex := strings.LastIndex(propertyId, ">")
	if beIndex == -1 {
		return "", nil
	}

	if dIndex > beIndex {
		prefix := propertyId[0:dIndex]
		return prefix, nil
	}

	bsIndex := strings.LastIndex(propertyId, "<")
	if dIndex > bsIndex {
		id = propertyId[:bsIndex]
		dIndex = strings.LastIndex(id, ".")
	}
	prefix := propertyId[0:dIndex]
	suffix := propertyId[dIndex+1:]
	bbIndex := strings.LastIndex(suffix, "<")
	if bbIndex == -1 {
		return prefix, nil
	}

	v := suffix[bbIndex+1 : len(suffix)-1]
	k, e := strings2.FromString(v, this.introspector.Registry())
	if e != nil {
		return "", e
	}
	this.key = k.Interface()
	return prefix, nil
}

func (this *Property) PropertyId() (string, error) {
	if this.id != "" {
		return this.id, nil
	}
	buff := strings2.New()
	if this.parent == nil {
		buff.Add(strings.ToLower(this.node.TypeName))
		buff.Add(this.node.CachedKey)
	} else {
		pi, err := this.parent.PropertyId()
		if err != nil {
			return "", err
		}
		buff.Add(pi)
		buff.Add(".")
		buff.Add(strings.ToLower(this.node.FieldName))
	}
	if this.key != nil {
		keyStr := strings2.New()
		keyStr.TypesPrefix = true
		buff.Add("<")
		buff.Add(keyStr.StringOf(this.key))
		buff.Add(">")
	}
	this.id = buff.String()
	return this.id, nil
}

func newProperty(node *types.RNode, propertyPath string, introspector common.IIntrospector) (*Property, error) {
	property := &Property{}
	property.node = node
	property.introspector = introspector
	if node.Parent != nil {
		prefix, err := property.setKeyValue(propertyPath)
		if err != nil {
			return nil, err
		}
		pi, err := newProperty(node.Parent, prefix, introspector)
		if err != nil {
			return nil, err
		}
		property.parent = pi
	} else {
		index1 := strings.Index(propertyPath, "<")
		index2 := strings.Index(propertyPath, ">")
		if index1 != -1 && index2 != -1 && index2 > index1 {
			k, e := strings2.FromString(propertyPath[index1+1:index2], property.introspector.Registry())
			if e != nil {
				return nil, e
			}
			property.key = k.Interface()
		}
	}
	return property, nil
}
