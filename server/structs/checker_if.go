package structs

// Checker defines a struct that can perform protocol-specific checks
type Checker interface {
	Check(name string, config *ServiceConfig) (CheckResult, error)
}
