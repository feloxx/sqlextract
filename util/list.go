package util

func ListDistinct(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func ListFilterTrim(list []string) []string {
	var result []string
	for _, value := range list {
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}