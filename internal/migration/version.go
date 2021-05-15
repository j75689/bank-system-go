package migration

// Version is a migrate version of database
type Version struct {
	ID   int64
	Name string
}
