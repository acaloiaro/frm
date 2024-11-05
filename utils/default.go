package utils

// SafelyDereference safely dereferences pointer by using a default when the pointer is nil
func SafelyDereference[T any](pointer *T, defaultValue T) T {
	if pointer != nil {
		return *pointer
	} else {
		return defaultValue
	}
}
