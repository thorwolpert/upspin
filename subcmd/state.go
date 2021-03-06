// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package subcmd // import "upspin.io/subcmd"

import (
	"fmt"
	"os"

	"upspin.io/bind"
	"upspin.io/client"
	"upspin.io/upspin"
)

// State describes the state of a subcommand.
// See the comments for Exitf to see how Interactive is used.
// It allows a program to run multiple commands.
type State struct {
	Name        string        // Name of the subcommand we are running.
	Config      upspin.Config // Config; may be nil.
	Client      upspin.Client // Client; may be nil.
	Interactive bool          // Whether the command is line-by-line.
	ExitCode    int           // Exit with non-zero status for minor problems.
}

// NewState returns a new State for the named subcommand.
func NewState(name string) *State {
	return &State{Name: name}
}

// Init initializes the config and client for the State.
func (s *State) Init(config upspin.Config) {
	var cl upspin.Client
	if config != nil {
		cl = client.New(config)
	}
	s.Config = config
	s.Client = cl
}

// Exitf prints the error and exits the program.
// If we are interactive, it calls panic("exit"), which is intended to be recovered
// from by the calling interpreter.
// We don't use log (although the packages we call do) because the errors
// are for regular people.
func (s *State) Exitf(format string, args ...interface{}) {
	format = fmt.Sprintf("upspin: %s: %s\n", s.Name, format)
	fmt.Fprintf(os.Stderr, format, args...)
	if s.Interactive {
		panic("exit")
	}
	s.Cleanup()
	os.Exit(1)
}

// Exit calls s.Exitf with the error.
func (s *State) Exit(err error) {
	s.Exitf("%s", err)
}

// Failf logs the error and sets the exit code. It does not exit the program.
func (s *State) Failf(format string, args ...interface{}) {
	format = fmt.Sprintf("upspin: %s: %s\n", s.Name, format)
	fmt.Fprintf(os.Stderr, format, args...)
	s.ExitCode = 1
}

// Fail calls s.Failf with the error.
func (s *State) Fail(err error) {
	s.Failf("%v", err)
}

// Cleanup terminates any necessary state.
// Keep this around because there might be some one day.
func (s *State) Cleanup() {
}

// KeyServer returns the KeyServer for the root of the name, or exits on failure.
func (s *State) KeyServer() upspin.KeyServer {
	key, err := bind.KeyServer(s.Config, s.Config.KeyEndpoint())
	if err != nil {
		s.Exit(err)
	}
	return key
}

// DirServer returns the DirServer for the root of the name, or exits on failure.
func (s *State) DirServer(name upspin.PathName) upspin.DirServer {
	dir, err := s.Client.DirServer(name)
	if err != nil {
		s.Exit(err)
	}
	return dir
}
