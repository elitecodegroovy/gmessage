package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func doReflectStruct() {
	type Contact struct {
		Name string "check:len(3,40)"
		Id   int    "check:range(1,999999)"
	}
	person := Contact{"Bjork", 0xDEEDED}
	personType := reflect.TypeOf(person)
	if nameField, ok := personType.FieldByName("Name"); ok {
		fmt.Printf("%q %q %q\n", nameField.Type, nameField.Name, nameField.Tag)
	}
}

// Data Model
type Dish struct {
	Id     int
	Name   string
	Origin string
	Query  func()
}

// Example of how to use Go's reflection
// Print the attributes of a Data Model
func attributes(m interface{}) map[string]reflect.Type {
	typ := reflect.TypeOf(m)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
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

func doReflectSlice() {
	nums := []string{"123", "456", "789"}
	sliceValue := reflect.ValueOf(nums)
	value := sliceValue.Index(len(nums) - 1)
	value.SetString("John")
	fmt.Println(nums)
}

func doReplaceValue() {
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
	UserName string `tag_name:"tag 1"`
	NickName string `tag_name:"tag 2"`
	Age      int    ` tag_name:"tag 3"`
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
func doReflectStructBasic() {
	u := &User{
		UserName: "梦放飞",
		NickName: "梦子",
		Age:      30,
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

func changeValue() {
	a := "C programming"
	d := reflect.ValueOf(&a).Elem()      //变量d，拥有地址
	pa := d.Addr().Interface().(*string) //获取string指针
	*pa = "Go Programming"               //赋值给指针的值
	fmt.Println(a)

	isAddr()
}

func isAddr() {
	x := "immutable"                  // 不是常量
	a := reflect.ValueOf("immutable") // 不是常量
	b := reflect.ValueOf(x)           // 不是常量
	c := reflect.ValueOf(&x)          // 不是常量
	d := c.Elem()                     // 是常量

	fmt.Println("a是常量：", a.CanAddr()) // "false"
	fmt.Println("b是常量：", b.CanAddr()) // "false"
	fmt.Println("c是常量：", c.CanAddr()) // "false"
	fmt.Println("d是常量：", d.CanAddr()) // "true"
}

type TypeVariable struct {
	noImitation struct{}
}

// tyvarUnderlyingType is used to discover types that are type variables.
// Namely, any type variable must be convertible to `TypeVariable`.
var tyvarUnderlyingType = reflect.TypeOf(TypeVariable{})

type A TypeVariable
type B TypeVariable
type C TypeVariable
type D TypeVariable
type E TypeVariable
type F TypeVariable
type G TypeVariable
type TypeError string

func (te TypeError) Error() string {
	return string(te)
}

func pe(format string, v ...interface{}) TypeError {
	return TypeError(fmt.Sprintf(format, v...))
}

func ppe(format string, v ...interface{}) {
	panic(pe(format, v...))
}

// Typed corresponds to the information returned by `Check`.
type Typed struct {
	// In correspondence with the `as` parameter to `Check`.
	Args []reflect.Value

	// In correspondence with the return types of `f` in `Check`.
	Returns []reflect.Type

	// The type environment generated via unification in `Check`.
	// (Its usefulness in the public API is questionable.)
	TypeEnv map[string]reflect.Type
}

func Check(f interface{}, as ...interface{}) *Typed {
	rf := reflect.ValueOf(f)
	tf := rf.Type()

	if tf.Kind() == reflect.Ptr {
		rf = reflect.Indirect(rf)
		tf = rf.Type()
	}
	if tf.Kind() != reflect.Func {
		ppe("The type of `f` must be a function, but it is a '%s'.", tf.Kind())
	}
	if tf.NumIn() != len(as) {
		ppe("`f` expects %d arguments, but only %d were given.",
			tf.NumIn(), len(as))
	}

	// Populate the argument value list.
	args := make([]reflect.Value, len(as))
	for i := 0; i < len(as); i++ {
		args[i] = reflect.ValueOf(as[i])
	}

	// Populate our type variable environment through unification.
	tyenv := make(tyenv)
	for i := 0; i < len(args); i++ {
		tp := typePair{tyenv, tf.In(i), args[i].Type()}

		// Mutates the type variable environment.
		if err := tp.unify(tp.param, tp.input); err != nil {
			argTypes := make([]string, len(args))
			for i := range args {
				argTypes[i] = args[i].Type().String()
			}
			ppe("\nError type checking\n\t%s\nwith argument types\n\t(%s)\n%s",
				tf, strings.Join(argTypes, ", "), err)
		}
	}

	// Now substitute those types into the return types of `f`.
	retTypes := make([]reflect.Type, tf.NumOut())
	for i := 0; i < tf.NumOut(); i++ {
		retTypes[i] = (&returnType{tyenv, tf.Out(i)}).tysubst(tf.Out(i))
	}
	return &Typed{args, retTypes, map[string]reflect.Type(tyenv)}
}

func (tp typePair) unify(param, input reflect.Type) error {
	if tyname := tyvarName(input); len(tyname) > 0 {
		return tp.error("Type variables are not allowed in the types of " +
			"arguments.")
	}
	if tyname := tyvarName(param); len(tyname) > 0 {
		if cur, ok := tp.tyenv[tyname]; ok && cur != input {
			return tp.error("Type variable %s expected type '%s' but got '%s'.",
				tyname, cur, input)
		} else if !ok {
			tp.tyenv[tyname] = input
		}
		return nil
	}
	if param.Kind() != input.Kind() {
		return tp.error("Cannot unify different kinds of types '%s' and '%s'.",
			param, input)
	}

	switch param.Kind() {
	case reflect.Array:
		return tp.unify(param.Elem(), input.Elem())
	case reflect.Chan:
		if param.ChanDir() != input.ChanDir() {
			return tp.error("Cannot unify '%s' with '%s' "+
				"(channel directions are different: '%s' != '%s').",
				param, input, param.ChanDir(), input.ChanDir())
		}
		return tp.unify(param.Elem(), input.Elem())
	case reflect.Func:
		if param.NumIn() != input.NumIn() || param.NumOut() != input.NumOut() {
			return tp.error("Cannot unify '%s' with '%s'.", param, input)
		}
		for i := 0; i < param.NumIn(); i++ {
			if err := tp.unify(param.In(i), input.In(i)); err != nil {
				return err
			}
		}
		for i := 0; i < param.NumOut(); i++ {
			if err := tp.unify(param.Out(i), input.Out(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		if err := tp.unify(param.Key(), input.Key()); err != nil {
			return err
		}
		return tp.unify(param.Elem(), input.Elem())
	case reflect.Ptr:
		return tp.unify(param.Elem(), input.Elem())
	case reflect.Slice:
		return tp.unify(param.Elem(), input.Elem())
	}

	// The only other container types are Interface and Struct.
	// I am unsure about what to do with interfaces. Mind is fuzzy.
	// Structs? I don't think it really makes much sense to use type
	// variables inside of them.
	return nil
}

// tyenv maps type variable names to their inferred Go type.
type tyenv map[string]reflect.Type

// typePair represents a pair of types to be unified. They act as a way to
// report sensible error messages from within the unification algorithm.
//
// It also includes a type environment, which is mutated during unification.
type typePair struct {
	tyenv tyenv
	param reflect.Type
	input reflect.Type
}

func (tp typePair) error(format string, v ...interface{}) error {
	return pe("Type error when unifying type '%s' and '%s': %s",
		tp.param, tp.input, fmt.Sprintf(format, v...))
}

// returnType corresponds to the type of a single return value of a function,
// in which the type may be parametric. It also contains a type environment
// constructed from unification.
type returnType struct {
	tyenv tyenv
	typ   reflect.Type
}

func (rt returnType) panic(format string, v ...interface{}) {
	ppe("Error substituting in return type '%s': %s",
		rt.typ, fmt.Sprintf(format, v...))
}

// tysubst attempts to substitute all type variables within a single return
// type with their corresponding Go type from the type environment.
//
// tysubst will panic if a type variable is unbound, or if it encounters a
// type that cannot be dynamically created. Such types include arrays,
// functions and structs. (A limitation of the `reflect` package.)
func (rt returnType) tysubst(typ reflect.Type) reflect.Type {
	if tyname := tyvarName(typ); len(tyname) > 0 {
		if thetype, ok := rt.tyenv[tyname]; !ok {
			rt.panic("Unbound type variable %s.", tyname)
		} else {
			return thetype
		}
	}

	switch typ.Kind() {
	case reflect.Array:
		rt.panic("Cannot dynamically create Array types.")
	case reflect.Chan:
		return reflect.ChanOf(typ.ChanDir(), rt.tysubst(typ.Elem()))
	case reflect.Func:
		rt.panic("Cannot dynamically create Function types.")
	case reflect.Interface:
		// rt.panic("TODO")
		// Not sure if this is right.
		return typ
	case reflect.Map:
		return reflect.MapOf(rt.tysubst(typ.Key()), rt.tysubst(typ.Elem()))
	case reflect.Ptr:
		return reflect.PtrTo(rt.tysubst(typ.Elem()))
	case reflect.Slice:
		return reflect.SliceOf(rt.tysubst(typ.Elem()))
	case reflect.Struct:
		rt.panic("Cannot dynamically create Struct types.")
	case reflect.UnsafePointer:
		rt.panic("Cannot dynamically create unsafe.Pointer types.")
	}

	// We've covered all the composite types, so we're only left with
	// base types.
	return typ
}

func tyvarName(t reflect.Type) string {
	if !t.ConvertibleTo(tyvarUnderlyingType) {
		return ""
	}
	return t.Name()
}

// Memo has a parametric type:
//
//	func Memo(f func(A) B) func(A) B
//
// Memo memoizes any function of a single argument that returns a single value.
// The type `A` must be a Go type for which the comparison operators `==` and
// `!=` are fully defined (this rules out functions, maps and slices).
func Memo(f interface{}) interface{} {
	chk := Check(new(func(func(A) B)), f)
	vf := chk.Args[0]

	saved := make(map[interface{}]reflect.Value)
	memo := func(in []reflect.Value) []reflect.Value {
		val := in[0].Interface()
		ret, ok := saved[val]
		if ok {
			return []reflect.Value{ret}
		}

		ret = call1(vf, in[0])
		saved[val] = ret
		return []reflect.Value{ret}
	}
	return reflect.MakeFunc(vf.Type(), memo).Interface()
}

/**
目标函数模型为：
func CallTypeFunc(f func(A) B, xs []A) []B
*/
//Only support Go 1.1以上的发布版本
func CallTypeFunc(f interface{}, ps interface{}) interface{} {
	vf := reflect.ValueOf(f)
	vps := reflect.ValueOf(ps)

	// 3) Map's return type must be `[]B1` where `B == B1`.
	tys := reflect.SliceOf(vf.Type().Out(0))

	vys := reflect.MakeSlice(tys, vps.Len(), vps.Len())
	for i := 0; i < vps.Len(); i++ {
		y := vf.Call([]reflect.Value{vps.Index(i)})[0]
		vys.Index(i).Set(y)
	}
	return vys.Interface()
}

func CallTypeFunc1(f, xs interface{}) interface{} {
	chk := Check(new(func(func(A) B, []A) []B), f, xs)
	vf, vxs, tys := chk.Args[0], chk.Args[1], chk.Returns[0]

	xsLen := vxs.Len()
	vys := reflect.MakeSlice(tys, xsLen, xsLen)
	for i := 0; i < xsLen; i++ {
		vy := vf.Call([]reflect.Value{vxs.Index(i)})[0]
		vys.Index(i).Set(vy)
	}
	return vys.Interface()
}
