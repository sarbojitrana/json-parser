package utils

type Resolver struct {
	mappings map[string]string
}

func NewResolver(mappings map[string]string) *Resolver {
	return &Resolver{
		mappings: mappings,
	}
}

func (r *Resolver) Resolve(data any) any {
	switch v := data.(type) {

	case map[string]any:
		return r.resolveMap(v)

	case []any:
		return r.resolveSlice(v)

	default:
		return v
	}
}

func (r *Resolver) resolveMap(data map[string]any) map[string]any {
	result := make(map[string]any, len(data))

	for key, value := range data {
		resolvedKey := key

		if mapped, ok := r.mappings[key]; ok {
			resolvedKey = mapped
		}

		result[resolvedKey] = r.Resolve(value)
	}

	return result
}

func (r *Resolver) resolveSlice(data []any) []any {
	result := make([]any, len(data))

	for i, value := range data {
		result[i] = r.Resolve(value)
	}

	return result
}