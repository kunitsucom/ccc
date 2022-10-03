package config

// nolint: deadcode,gochecknoglobals,unused,varcheck
var (
	subcommandVersion bool
	version           string
	revision          string
	branch            string
	timestamp         string
)

func SubcommandVersion() bool { return subcommandVersion }
func Version() string         { return version }
func Revision() string        { return revision }
func Branch() string          { return branch }
func Timestamp() string       { return timestamp }
