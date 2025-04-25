package internal

//go:generate go tool counterfeiter -generate

//counterfeiter:generate . Transformer
type Transformer interface {
	Transform(int) int
}

func Sum(a, b int, transformer Transformer) int {
	return transformer.Transform(a + b)
}
