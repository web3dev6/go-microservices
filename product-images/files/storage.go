package files

import "io"

// Storage defines the behavior for file operations
// Implementations may be of the type -> local disk, or cloud storage, etc
type Storage interface {
	Save(path string, file io.Reader) error
}
