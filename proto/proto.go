package proto

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Converter is a struct containing configuration to convert a parsed protoc
// file to Go types.
type Converter struct {
	// Packages maps file paths to the package of the file.
	Packages map[string]*protogen.File
}

// searchForType searches for the exact type matching the provided suffix.
//
// To match, suffix must exactly match the type name or form the final
// identifier segments of the type name. "xo.sample" will match
// "com.github.xo.xo.sample" and "xo.sample" but not "xoxo.sample".
//
// This function returns an error if the number of messages found is not
// exactly one.
func (c Converter) searchForType(suffix string) (*protogen.Message, error) {
	// Walk through all messages and find matches.
	var results []*protogen.Message
	var processMsg func(*protogen.Message)
	processMsg = func(msg *protogen.Message) {
		for _, nested := range msg.Messages {
			processMsg(nested)
		}
		fullName := string(msg.Desc.FullName())
		if strings.HasSuffix(fullName, "."+suffix) || fullName == suffix {
			results = append(results, msg)
		}
	}
	for _, pkg := range c.Packages {
		for _, m := range pkg.Messages {
			processMsg(m)
		}
	}
	// Return errors if applicable.
	switch len(results) {
	case 0:
		return nil, fmt.Errorf("no types with suffix %q found", suffix)
	case 1:
		return results[0], nil
	default:
		names := make([]string, 0, len(results))
		for _, r := range results {
			names = append(names, string(r.Desc.FullName()))
		}
		return nil, fmt.Errorf(
			"too many names matching suffix %q:\n"+
				"\t- %s",
			suffix, strings.Join(names, "\n\t- "),
		)
	}
}

// pkgOf returns the protogen.File containing the object declared by the
// descriptor.
func (c Converter) pkgOf(d protoreflect.Descriptor) (*protogen.File, error) {
	path := d.ParentFile().Path()
	if f, ok := c.Packages[path]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("file with path %q not provided in Packages", path)
}
