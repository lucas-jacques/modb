module github.com/lucasjacques/modb/drivers/repo.pgx

go 1.23.2

replace github.com/lucasjacques/modb => ../..

require (
	github.com/jackc/pgx/v5 v5.7.2
	github.com/lucasjacques/modb v0.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
