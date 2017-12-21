package main

import "fmt"

type FuzzyBool struct{ value float32 }

func New(value interface{}) (*FuzzyBool, error) {
	amount, err := float32ForValue(value)
	return &FuzzyBool{amount}, err
}

func float32ForValue(value interface{}) (fuzzy float32, err error) {
	switch value := value.(type) { // shadow variable
	case float32:
		fuzzy = value
	case float64:
		fuzzy = float32(value)
	case int:
		fuzzy = float32(value)
	case bool:
		fuzzy = 0
		if value {
			fuzzy = 1
		}
	default:
		return 0, fmt.Errorf("float32ForValue(): %v is not a "+
			"number or Boolean", value)
	}
	if fuzzy < 0 {
		fuzzy = 0
	} else if fuzzy > 1 {
		fuzzy = 1
	}
	return fuzzy, nil
}

func (fuzzy *FuzzyBool) String() string {
	return fmt.Sprintf("%.0f%%", 100*fuzzy.value)
}

func (fuzzy *FuzzyBool) Set(value interface{}) (err error) {
	fuzzy.value, err = float32ForValue(value)
	return err
}

func (fuzzy *FuzzyBool) Copy() *FuzzyBool {
	return &FuzzyBool{fuzzy.value}
}

func (fuzzy *FuzzyBool) Not() *FuzzyBool {
	return &FuzzyBool{1 - fuzzy.value}
}

func (fuzzy *FuzzyBool) And(first *FuzzyBool, rest ...*FuzzyBool) *FuzzyBool {
	minimum := fuzzy.value
	rest = append(rest, first)
	for _, other := range rest {
		if minimum > other.value {
			minimum = other.value
		}
	}
	return &FuzzyBool{minimum}
}

func (fuzzy *FuzzyBool) Or(first *FuzzyBool, rest ...*FuzzyBool) *FuzzyBool {
	max := fuzzy.value
	rest = append(rest, first)
	for _, other := range rest {
		if max < other.value {
			max = other.value
		}
	}
	return &FuzzyBool{max}
}

func (fuzzy *FuzzyBool) Less(other *FuzzyBool) bool {
	return fuzzy.value < other.value
}

func (fuzzy *FuzzyBool) Equal(other *FuzzyBool) bool {
	return fuzzy.value == other.value
}

func (fuzzy *FuzzyBool) Bool() bool {
	return fuzzy.value >= .5
}
func (fuzzy *FuzzyBool) Float() float64 {
	return float64(fuzzy.value)
}

func process(a, b, c, d *FuzzyBool) {
	fmt.Println("Original:", a, b, c, d)
	fmt.Println("Not: ", a.Not(), b.Not(), c.Not(), d.Not())
	fmt.Println("Not Not: ", a.Not().Not(), b.Not().Not(), c.Not().Not(),
		d.Not().Not())
	fmt.Print("0.And(.25)→", a.And(b), "• .25.And(.75)→", b.And(c),
		"• .75.And(1)→", c.And(d), " • .25.And(.75,1)→", b.And(c, d), "\n")
	fmt.Print("0.Or(.25)→", a.Or(b), "• .25.Or(.75)→", b.Or(c),
		"• .75.Or(1)→", c.Or(d), " • .25.Or(.75,1)→", b.Or(c, d), "\n")
	fmt.Println("a < c, a == c, a > c:", a.Less(c), a.Equal(c), c.Less(a))
	fmt.Println("Bool: ", a.Bool(), b.Bool(), c.Bool(), d.Bool())
	fmt.Println("Float: ", a.Float(), b.Float(), c.Float(), d.Float())
}

func doFuzzy() {
	a, _ := New(0)   // Safe to ignore err value when using
	b, _ := New(.25) // known valid values; must check if using
	c, _ := New(.75) // variables though.
	d := c.Copy()
	if err := d.Set(1); err != nil {
		fmt.Println(err)
	}
	process(a, b, c, d)
	s := []*FuzzyBool{a, b, c, d}
	fmt.Println(s)
}

//func main() {
//	doFuzzy()
//}
