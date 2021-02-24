package checkers

import "status-api/structs"

// HTTP -
type HTTP struct {
}

// Check -
func (HTTP) Check(name string, c *structs.ServiceConfig) (structs.Result, error) {
	return structs.Result{
		Status: "up",
		URL:    c.FriendlyURL,
	}, nil
}

// Interface guard
var _ structs.Checker = (*HTTP)(nil)
