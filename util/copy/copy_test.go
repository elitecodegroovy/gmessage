package copy

import (
	"fmt"
	"testing"
)

func TestCopyMap(t *testing.T) {
	input := map[string]interface{}{
		"bob": map[string]interface{}{
			"emails": []string{"a", "b"},
		},
	}

	dup, err := Copy(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", dup)
	fmt.Println("\n", dup.(map[string]interface{})["bob"])
}
