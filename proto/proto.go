package proto

import "google.golang.org/protobuf/compiler/protogen"

// Converter is a struct containing configuration to convert a parsed protoc
// file to Go types.
type Converter struct {
	// Packages maps file paths to the package of the file.
	Packages map[string]*protogen.File
}
