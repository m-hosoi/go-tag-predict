package lambda

import "fmt"

func ExampleMapString() {
	ar := []string{"A", "B", "C"}
	ar = MapString(ar, func(s string) string {
		return s + s
	})
	for _, s := range ar {
		fmt.Println(s)
	}
	// Output:
	// AA
	// BB
	// CC
}
func ExampleFilterString() {
	ar := []string{"A", "B", "C"}
	ar = FilterString(ar, func(s string) bool {
		return s != "B"
	})
	for _, s := range ar {
		fmt.Println(s)
	}
	// Output:
	// A
	// C
}
func ExampleMapIntString() {
	m := map[int]string{
		1: "A",
		2: "B",
		3: "C",
	}
	ar := []int{1, 2, 3}
	for _, s := range MapIntString(ar, func(i int) string {
		return m[i]
	}) {
		fmt.Println(s)
	}

	// Output:
	// A
	// B
	// C
}
func ExampleMapStringInt() {
	m := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
	}
	ar := []string{"A", "B", "C"}
	for _, i := range MapStringInt(ar, func(s string) int {
		return m[s]
	}) {
		fmt.Println(i)
	}

	// Output:
	// 1
	// 2
	// 3
}
