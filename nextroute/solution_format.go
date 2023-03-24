package nextroute

import (
	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// The Formatter interface is used to create custom JSON output.
type Formatter interface {
	ToOutput(Solution) any
}

// NewDefaultFormatter creates a new NewDefaultFormatter.
func NewDefaultFormatter() Formatter {
	connect.Connect(con, &newDefaultFormatter)
	return newDefaultFormatter()
}

// NewExtensiveFormatter creates a NewExtensiveFormatter which outputs
// additional solution information.
func NewExtensiveFormatter(p []alns.ProgressionEntry) Formatter {
	connect.Connect(con, &newExtensiveFormatter)
	return newExtensiveFormatter(p)
}
