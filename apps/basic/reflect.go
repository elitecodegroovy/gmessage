package main

import (
	"reflect"
	"fmt"
	"strconv"
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

type User struct {
	UserName  string  `tag_name:"tag 1"`
	NickName  string  `tag_name:"tag 2"`
	Age       int    ` tag_name:"tag 3"`
}

func (u *User) printFields() {
	val := reflect.ValueOf(u).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n",
			typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
}
func doReflectStructBasic(){
	u := &User{
		UserName: "梦放飞",
		NickName:  "梦子",
		Age:       30,
	}
	u.printFields()
}


// formatAtom formats a value without inspecting its internal structure.
// It is a copy of the the function in gopl.io/ch11/format.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
		// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr,
		reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

//print value of the reflect value
func printValue(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			printValue(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			printValue(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			printValue(fmt.Sprintf("%s[%s]", path,
				formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			printValue(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			printValue(path+".value", v.Elem())
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

func changeValue(){
	a := "C programming"
	d := reflect.ValueOf(&a).Elem()		    //变量d，拥有地址
	pa := d.Addr().Interface().(*string)   //获取string指针
	*pa = "Go Programming"                 //赋值给指针的值
	fmt.Println(a)

	isAddr()
}

func isAddr(){
	x := "immutable"                       	// 不是常量
	a := reflect.ValueOf("immutable")    // 不是常量
	b := reflect.ValueOf(x)      			// 不是常量
	c := reflect.ValueOf(&x)     			// 不是常量
	d := c.Elem()							// 是常量

	fmt.Println("a是常量：", a.CanAddr()) // "false"
	fmt.Println("b是常量：",b.CanAddr())  // "false"
	fmt.Println("c是常量：",c.CanAddr())  // "false"
	fmt.Println("d是常量：",d.CanAddr())  // "true"

}
