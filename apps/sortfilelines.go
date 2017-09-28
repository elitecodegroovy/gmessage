package main
import (
	"bytes"
	"io/ioutil"
	"sort"
)

// allow [][]byte to implement the sort.Interface interface
type sortLines [][]byte

// bytes.Compare compares the byte slices lexicographically (alphabetically)
func (l sortLines) Less(i, j int) bool { return bytes.Compare(l[i], l[j]) > 0 }
func (l sortLines) Len() int           { return len(l) }
func (l sortLines) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func SortFile(name string) error {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	lines := bytes.Split(content, []byte{'\n'})
	sort.Sort(sortLines(lines))

	content = bytes.Join(lines, []byte{'\n'})
	return ioutil.WriteFile(name, content, 0644)
}

//func main() {
//	SortFile("web2.txt")
//}


