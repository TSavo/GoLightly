package govirtual

import (
	"fmt"
	"strings"
)

//A Value is an allocation of memory of undefined type
type Value interface {
	Get() interface{}
	Set(interface{})
	String() string
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

//A MemoryPointer points at specific Address in Memory
type MemoryPointer struct {
	*Memory
	Address interface{}
	Name string
}

func (this *MemoryPointer) Get() interface{} {
	return this.Memory[this.Address]
}

func (this *MemoryPointer) Set(value interface{}) {
	this.Memory[this.Address].Set(value)
}

func (this *MemoryPointer) String() string {
	return fmt.Sprintf("*%v", this.Address)
}

func (this *MemoryPointer) Dereference() interface{} {
	return Dereference(this)
}

func (this *MemoryPointer) CompletelyDererence() interface{} {
	return CompletelyDereference(this)
}

func CompletelyDereference(value interface{}) interface{} {
	switch x := value.(type) {
	case Value:
		return CompletelyDereference(x.Get())
	default:
		return x
	}
}

func Dereference(value interface{}) interface{} {
	switch x := value.(type) {
	case Value:
		return x.Get()
	default:
		return x
	}
}