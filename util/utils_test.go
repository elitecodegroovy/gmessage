package util

import (
	"fmt"
	"testing"
	"time"
)

func TestSwapCase(t *testing.T) {
	input, expected := "Hello, World", "hELLO, wORLD"
	result := SwapCase(input)
	if result != expected {
		t.Errorf("SwapCase(%q) == %q, expected %q", input, result, expected)
	}
}

func TestReverse(t *testing.T) {
	input, expected := "Hello, World", "dlroW ,olleH"
	result := Reverse(input)
	if result != expected {
		t.Errorf("Reverse (%q) == %q, expected %q", input, result, expected)
	}
}

// -------------------benchmark ------------------------------------
//Benchmark for SwapCase function
func BenchmarkSwapCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SwapCase("Hello, World")
	}
}

//Benchmark for Reverse function
func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reverse("Hello, World")
	}
}

// -------------------verifyting ------------------------------------
//Example code for Reverse function
func ExampleReverse() {
	fmt.Println(Reverse("Hello, World"))
	// Output: dlroW ,olleH
}

//Example code for Reverse function
func ExampleSwapCase() {
	fmt.Println(SwapCase("Hello, World"))
	// Output: hELLO, wORLD
}

func TestSwapCaseInParallel(t *testing.T) {
	t.Parallel()
	time.Sleep(1 * time.Second)
	input, expected := "Hello, World", "hELLO, wORLD"
	result := SwapCase(input)
	if result != expected {
		t.Errorf("SwapCase(%q) == %q, expected %q", input, result, expected)
	}
}

// Test case for the Reverse function to execute in parallel
func TestReverseInParallel(t *testing.T) {
	t.Parallel()
	// Delaying 2 seconds for the sake of demonstration
	time.Sleep(2 * time.Second)
	input, expected := "Hello, World", "dlroW ,olleH"
	result := Reverse(input)
	if result != expected {
		t.Errorf("Reverse(%q) == %q, expected %q", input, result, expected)
	}
}

func TestParsingTime(t *testing.T){
	form := "2006-01-02 15:04:05"
	time , err := time.Parse(form, "2017-03-02 19:04:05")
	if err != nil {
		fmt.Println("parsing time error", err)
	}
	fmt.Println("time :", time)
}
