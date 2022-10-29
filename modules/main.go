package modules

type timekeeping struct{}
type vcaccess struct{}

// Exports modules
var (
	VCAccess    = &vcaccess{}
	Timekeeping = &timekeeping{}
)
