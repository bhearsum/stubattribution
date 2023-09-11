package dmg

import (
	"io"
)

type ReaderSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type ReaderAtSeeker interface {
	io.ReaderAt
	io.Seeker
}
