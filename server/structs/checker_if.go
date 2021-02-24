package structs

// Checker defines a protocol check
type Checker interface {
	Check(c Config) (Result, error)
}
