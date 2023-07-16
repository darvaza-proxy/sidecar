// Package flags provides a flag mapper for configuration structs
package flags

// Mapper is a generic flag/config mapper
type Mapper interface {
	Parse() error

	Var(v any, longOpt string, message string) Mapper
	VarP(v any, longOpt string, shortOpt rune, message string) Mapper
}

var registry = make(map[any]Mapper)

// PutMapper stores a Mapper to be accessed by GetMapper later
func PutMapper(ref any, m Mapper) {
	if ref == nil || m == nil {
		panic("invalid arguments")
	}

	registry[ref] = m
}

// GetMapper retrieves a previously stored Mapper
func GetMapper(ref any) Mapper {
	if m, ok := registry[ref]; ok {
		return m
	}

	return nil
}
