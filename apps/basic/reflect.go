package main

import (
	"reflect"
	"fmt"
)

func doReflectStruct(){
	type Contact struct {
		Name string "check:len(3,40)"
		Id int "check:range(1,999999)"
	}
	person := Contact{"Bjork", 0xDEEDED}
	personType := reflect.TypeOf(person)
	if nameField, ok := personType.FieldByName("Name"); ok {
		fmt.Printf("%q %q %q\n", nameField.Type, nameField.Name, nameField.Tag)
	}
}

// Data Model
type Dish struct {
	Id  int
	Name string
	Origin string
	Query func()
}

// Example of how to use Go's reflection
// Print the attributes of a Data Model
func attributes(m interface{}) (map[string]reflect.Type) {
	typ := reflect.TypeOf(m)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr{
		typ = typ.Elem()
	}

	// create an attribute data structure as a map of types keyed by a string.
	attrs := make(map[string]reflect.Type)
	// Only structs are supported so return an empty result if the passed object
	// isn't a struct
	if typ.Kind() != reflect.Struct {
		fmt.Printf("%v type can't have attributes inspected\n", typ.Kind())
		return attrs
	}

	// loop through the struct's fields and set the map
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		if !p.Anonymous {
			attrs[p.Name] = p.Type
		}
	}

	return attrs
}

func doReflectSlice(){
	nums := []string{"123", "456", "789"}
	sliceValue := reflect.ValueOf(nums)
	value := sliceValue.Index(len(nums)-1)
	value.SetString("John")
	fmt.Println(nums)
}

func doReplaceValue(){
	count := 10
	//the value can't be set.
	if value := reflect.ValueOf(count); value.CanSet() {
		//value.SetInt(20) //it never reach to the point
	}
	fmt.Print(count, " ")
	value := reflect.ValueOf(&count)
	// Can't call SetInt() on value since value is a * int not an int
	pointee := value.Elem()
	pointee.SetInt(30) // OK. Can replace a pointed-to value.
	fmt.Println(count)
}

type Foo struct {
	FirstName string `tag_name:"tag 1"`
	LastName  string `tag_name:"tag 2"`
	Age       int    `tag_name:"tag 3"`
}

func (f *Foo) reflect() {
	val := reflect.ValueOf(f).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n",
			typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
}
func doReflectStructBasic(){
	f := &Foo{
		FirstName: "Drew",
		LastName:  "Olson",
		Age:       30,
	}

	f.reflect()
}
