package constant

import "time"

const (
	// AppBanner is the application banner
	AppBanner = `
  ▄▄▌ ▐ ▄▌ ▄▄▄· .▄▄ ·  ▄▄▄· ▄▄▄▄· ▪  
  ██· █▌▐█▐█ ▀█ ▐█ ▀. ▐█ ▀█ ▐█ ▀█▪██ 
  ██▪▐█▐▐▌▄█▀▀█ ▄▀▀▀█▄▄█▀▀█ ▐█▀▀█▄▐█·
  ▐█▌██▐█▌▐█ ▪▐▌▐█▄▪▐█▐█ ▪▐▌██▄▪▐█▐█▌
   ▀▀▀▀ ▀▪ ▀  ▀  ▀▀▀▀  ▀  ▀ ·▀▀▀▀ ▀▀▀`
	// AppName is the name of the application
	AppName = "wasabi"
	// DefaultTimeout is the default timeout for http request
	DefaultTimeout = 5 * time.Second
)

var (
	// AppVersion is the version of the application.
	AppVersion = "0.0.0"
	// Branch is the git branch of the application
	Branch = "main"
	// Buildtime is the time when the application was built
	Buildtime = "2006 Jan 02 15:04:05"
	// CommitHash is the git commit hash of the application
	CommitHash = "0b00b135"
	// CommitMsg is the git commit message of the application.
	CommitMsg = "n/a"
	// ReleaseVersion is the release version of the application.
	ReleaseVersion = ""
)
