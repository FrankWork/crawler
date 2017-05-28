package main

import (
	"crypto/sha1"
	"fmt"
)

func fingerPrint(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)

	return fmt.Sprintf("%x", bytes)
}
func fpExample() {
	fp := fingerPrint("hello world")
	fmt.Println(fp)
}
