package libs

func RemoveElements[T comparable](arr []T, remove []T) []T {
	var result []T
	set := NewSet[T]()

	for _, v := range remove {
		set.Add(v)
	}

	for _, v := range arr {
		if exist := set.Contains(v); !exist {
			result = append(result, v)
		}
	}

	return result
}
