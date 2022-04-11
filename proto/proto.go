package proto

// Converter is a struct containing configuration to convert a parsed protoc
// file to Go types.
type Converter struct {
	// PackageNames maps file paths to the package name of the file.
	PackageNames map[string]string
	// SkipPrefixes is a list of prefixes to skip when generating table names.
	SkipPrefixes []string
}
