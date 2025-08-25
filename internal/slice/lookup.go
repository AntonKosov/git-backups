package slice

func Lookup[Item any, Key comparable, Value any](input []Item, transform func(Item) (Key, Value)) (output map[Key]Value) {
	output = make(map[Key]Value, len(input))

	if input == nil {
		return output
	}

	for _, item := range input {
		k, v := transform(item)
		output[k] = v
	}

	return output
}
