// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"io"

	"github.com/pkg/errors"
	"github.com/tsuru/tsuru/event"
	"github.com/tsuru/tsuru/provision"
)

const defaultDockerBuilder = "docker"

var DefaultBuilder = defaultDockerBuilder

type BuildOpts struct {
	BuildFromFile  bool
	Rebuild        bool
	Redeploy       bool
	ArchiveURL     string
	ArchiveFile    io.Reader
	ArchiveTarFile io.ReadCloser
	ArchiveSize    int64
	ImageID        string
}

// Builder is the basic interface of this package.
type Builder interface {
	Build(p provision.BuilderDeploy, app provision.App, evt *event.Event, opts BuildOpts) (string, error)
}

var builders = make(map[string]Builder)

// PlatformBuilder is a builder where administrators can manage
// platforms (automatically adding, removing and updating platforms).
type PlatformBuilder interface {
	PlatformAdd(PlatformOptions) error
	PlatformUpdate(PlatformOptions) error
	PlatformRemove(name string) error
}

// PlatformOptions is the set of options provided to PlatformAdd and
// PlatformUpdate, in PlatformBuilder.
type PlatformOptions struct {
	Name   string
	Args   map[string]string
	Input  io.Reader
	Output io.Writer
}

// Register registers a new builder in the Builder registry.
func Register(name string, builder Builder) {
	builders[name] = builder
}

// Get gets the named builder from the registry.
func Get(name string) (Builder, error) {
	b, ok := builders[name]
	if !ok {
		return nil, errors.Errorf("unknown builder: %q", name)
	}
	return b, nil
}

func GetDefault() (Builder, error) {
	return Get(DefaultBuilder)
}

// Registry returns the list of registered builders.
func Registry() ([]Builder, error) {
	registry := make([]Builder, 0, len(builders))
	for _, b := range builders {
		registry = append(registry, b)
	}
	return registry, nil
}

func PlatformAdd(opts PlatformOptions) error {
	builders, err := Registry()
	if err != nil {
		return err
	}
	for _, b := range builders {
		if platformBuilder, ok := b.(PlatformBuilder); ok {
			return platformBuilder.PlatformAdd(opts)
		}
	}
	return errors.New("No builder available")
}

func PlatformUpdate(opts PlatformOptions) error {
	builders, err := Registry()
	if err != nil {
		return err
	}
	for _, b := range builders {
		if platformBuilder, ok := b.(PlatformBuilder); ok {
			return platformBuilder.PlatformUpdate(opts)
		}
	}
	return errors.New("No builder available")
}

func PlatformRemove(name string) error {
	builders, err := Registry()
	if err != nil {
		return err
	}
	for _, b := range builders {
		if platformBuilder, ok := b.(PlatformBuilder); ok {
			return platformBuilder.PlatformRemove(name)
		}
	}
	return errors.New("No builder available")
}
