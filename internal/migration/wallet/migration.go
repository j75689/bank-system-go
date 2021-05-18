package wallet

import "github.com/go-gormigrate/gormigrate/v2"

// Migrations is a collection of storage migration patterns
var Migrations = []*gormigrate.Migration{
	v202105162211,
	v202105181826,
	v202105190000,
}
