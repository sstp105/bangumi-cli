package bangumi

func paginate[T any](fetch func(offset int) ([]T, int, error)) ([]T, error) {
	var result []T
	total := 1

	for offset := 0; offset < total; offset += defaultPaginationLimit {
		data, pageTotal, err := fetch(offset)
		if err != nil {
			return nil, err
		}
		result = append(result, data...)
		total = pageTotal
	}

	return result, nil
}
