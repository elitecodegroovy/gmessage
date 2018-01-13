package main

import (
	"fmt"
	"unsafe"
)

func doArithmetic() {
	intArray := [...]int{10, 20}

	fmt.Printf("\nintArray: %v\n", intArray)

	intPtr := &intArray[0]
	fmt.Printf("\nintPtr=%p, *intPtr=%d.\n", intPtr, *intPtr)

	nextIntPtr := uintptr(unsafe.Pointer(intPtr)) + unsafe.Sizeof(intArray[0])

	intPtr = (*int)(unsafe.Pointer(nextIntPtr))

	fmt.Printf("\nintPtr=%p, *intPtr=%d.\n\n", intPtr, *intPtr)
}

func doConvert() {
	type MyStr string

	a := []MyStr{"10", "100", "1000"}
	// b := ([]MyStr)(a) // error: cannot convert a (type []MyStr) to type []MyStr
	b := *(*[]MyStr)(unsafe.Pointer(&a))

	b[0] = "0"

	fmt.Println("a =", a) // a = [0 100 1000]
	fmt.Println("b =", b) // b = [0 100 1000]

	a[2] = "CONVERT"

	fmt.Println("a =", a) // a = [0 100 CONVERT]
	fmt.Println("b =", b) // b = [0 100 CONVERT]
}

func main() {
	doArithmetic()
	doConvert()
}
