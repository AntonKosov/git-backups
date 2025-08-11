package slice

func Map[T1, T2 any](input []T1, transform func(T1) T2) (output []T2) {
	if input == nil {
		return nil
	}

	output = make([]T2, len(input))
	for i, item := range input {
		output[i] = transform(item)
	}

	return output
}
