package main

import (
	"fmt"
)

func main() {
	rc4_Reslut := RC4("rc4_2006", string("ABC123😀"))
	vv := RC4("rc4_2006",string(rc4_Reslut))
	fmt.Println(string(vv))
}
