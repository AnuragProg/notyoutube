package database

import (
	"io"
)

type Database interface {

	// to make sure the handler closes the connection properly
	io.Closer

}
