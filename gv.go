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
	GlideYAML     string          `json:"-"`
	GlideLock     string          `json:"-"`
	GlideConfig   *glide.Config   `json:"Glide,omitempty"`
	GlideLockfile *glide.Lockfile `json:"GlideLock,omitempty"`
}

// Init initializes a BuildInfo struct
func (b *BuildInfo) Init() error {
	// check glide data
	if b.GlideYAML != "" && b.GlideConfig == nil {
		// decode data
		dec, err := base64.StdEncoding.DecodeString(b.GlideYAML)
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

	// check glide lock
	if b.GlideLock != "" && b.GlideLockfile == nil {
		// decode data
		dec, err := base64.StdEncoding.DecodeString(b.GlideLock)
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
