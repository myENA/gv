package gv

import (
	// core
	"encoding/base64"

	// Glide YAML parsing
	glide "github.com/Masterminds/glide/cfg"
)

// BuildInfo contains application version information
type BuildInfo struct {
	Name          string          `json:"Name,omitempty"`
	Version       string          `json:"Version,omitempty"`
	Date          string          `json:"BuildDate,omitempty"`
	Branch        string          `json:"BuildBranch,omitempty"`
	Commit        string          `json:"BuildCommit,omitempty"`
	Build         string          `json:"BuildNumber,omitempty"`
	GlideData     string          `json:"-"`
	GlideLockData string          `json:"-"`
	GlideConfig   *glide.Config   `json:"Glide,omitempty"`
	GlideLockfile *glide.Lockfile `json:"GlideLock,omitempty"`
}

// Init initializes a BuildInfo struct
func (b *BuildInfo) Init() error {
	// check glide data
	if b.GlideData != "" && b.GlideConfig == nil {
		// decode data
		dec, err := base64.StdEncoding.DecodeString(b.GlideData)
		if err != nil {
			return err
		}
		// parse data
		cfg, err := glide.ConfigFromYaml(dec)
		if err != nil {
			return err
		}
		b.GlideConfig = cfg
	}

	// check glide lock data
	if b.GlideLockData != "" && b.GlideLockfile == nil {
		// decode data
		dec, err := base64.StdEncoding.DecodeString(b.GlideLockData)
		if err != nil {
			return err
		}
		// parse data
		cfg, err := glide.LockfileFromYaml(dec)
		if err != nil {
			return err
		}
		b.GlideLockfile = cfg
	}

	// return
	return nil
}
