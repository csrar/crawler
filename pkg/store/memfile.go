package store

//go:generate mockgen -source=memfile.go -destination=mocks/memfile_mock.go
type ImFile interface {
	Bytes() []byte
	Truncate(n int64) error
	WriteAt(b []byte, offset int64) (int, error)
}
