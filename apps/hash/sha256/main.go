package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func hashSHA256Str(s string) {

	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	fmt.Printf("origin: %s, sha256 hash: %x\n", s, bs)
}

func testSha256Str() {
	s := "sha256 芳华"
	hashSHA256Str(s)
}

func hashSHA256File(filePath string) (string, error) {
	var hashValue string
	file, err := os.Open(filePath)
	if err != nil {
		return hashValue, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return hashValue, err
	}
	hashInBytes := hash.Sum(nil)
	hashValue = hex.EncodeToString(hashInBytes)
	return hashValue, nil

}

func main() {
	//testSha256Str()
	filePath := "D:/VPN/advancedconf.jpg"
	if hash, err := hashSHA256File(filePath); err != nil {
		fmt.Printf(" %s, sha256 value: %s ", filePath, hash)
	} else {
		fmt.Printf(" %s, sha256 hash: %s ", filePath, err.Error())
	}
}
