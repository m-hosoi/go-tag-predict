package lambda

// MapString : string[] => string[]のmap関数
func MapString(ar []string, f func(string) string) []string {
	res := make([]string, len(ar))
	for i, s := range ar {
		res[i] = f(s)
	}
	return res
}

// MapStringInt : string[] => int[]のmap関数
func MapStringInt(ar []string, f func(string) int) []int {
	res := make([]int, len(ar))
	for i, s := range ar {
		res[i] = f(s)
	}
	return res
}

// MapStringInt : int[] => string[]のmap関数
func MapIntString(ar []int, f func(int) string) []string {
	res := make([]string, len(ar))
	for i, s := range ar {
		res[i] = f(s)
	}
	return res
}

// FilterString : string[] => string[]のfilter関数
func FilterString(ar []string, f func(string) bool) []string {
	res := make([]string, 0, len(ar))
	for _, s := range ar {
		if f(s) {
			res = append(res, s)
		}
	}
	return res
}
