package config

type ConnectionDB struct {
	URL string
}

func NewConnectionDB() ConnectionDB {
	return ConnectionDB{
		URL: "postgres://postgres:postgres@localhost:5432/go_app_db?sslmode=disable",
	}
}
