package entities


type DBScanner interface {
	Scan(dest ...any) error
}

type Entity interface {
	Columns() []string
	ScanRow(row DBScanner) error
}