package fileutil

import "fmt"

func ExampleExists() {
	fmt.Println(Exist("/etc/hosts"))
	fmt.Println(Exist("/etc/hosts.notexists"))
	// Output:
	// true
	// false
}
func ExampleGetGoPath() {
	ar := GetGopath()
	fmt.Println(len(ar))
	fmt.Println(Exist(ar[0]))
	// Output:
	// 1
	// true
}
