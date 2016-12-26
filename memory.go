package govirtual

import (
	"fmt"
	"strings"
)

//A Value is an allocation of memory of undefined type
type Value interface {
	Get() interface{}
	Set(interface{})
}

//An Address is something that can be used to address a Value in Memory
type Address interface{}

//Memory maps Addresses to Values
type Memory map[Address]Value

// A Variable is a named Value
type Variable struct {
	Value
	Name string
}

func (this *Variable) Get() interface{} {
	return this.Value.Get()
}

func (this *Variable) Set(value interface{}) {
	this.Value.Set(value)
}

func (this *Variable) String() string {
	return this.Name
}

//A Literal Value
type Literal struct {
	Value interface{}
}

func (this *Literal) Get() interface{} {
	return this.Value
}

func (this *Literal) Set(value interface{}) {
	this.Value = value
}

func (this *Literal) String() string {
	switch x := this.Value.(type) {
	case string:
		if strings.HasPrefix(x, ":") {
			return x
		} else {
			return fmt.Sprintf("\"%v\"", x)
		}
	default:
		return fmt.Sprintf("%v", x)
	}
}

//A Reference points at another value
type Reference struct {
	Value *Value
}

func (this *Reference) Get() interface{} {
	val := *(this.Value)
	return val.Get()
}

func (this *Reference) Set(value interface{}) {
	val := *(this.Value)
	val.Set(value)
}

func (this *Reference) String() string {
	return fmt.Sprintf("&%v", this.Value)
}
