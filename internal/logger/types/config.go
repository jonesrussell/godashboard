package types

// Config holds the configuration for the logger
type Config struct {
	Level      string
	OutputPath string
	MaxSize    int  // megabytes
	MaxBackups int  // number of backups
	MaxAge     int  // days
	Compress   bool // compress old files
	Debug      bool // development mode
}
