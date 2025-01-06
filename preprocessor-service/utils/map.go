package utils


func Map[A, B any](items []A, fun func(A) B) []B {
	result := make([]B, 0, len(items))
	for _, item := range items {
		result = append(result, fun(item))
	}
	return result
}
