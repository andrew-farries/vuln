// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package govulncheck

import (
	"encoding/json"
	"io"

	"github.com/andrew-farries/vuln/internal/osv"
)

// Handler handles messages to be presented in a vulnerability scan output
// stream.
type Handler interface {
	// Config communicates introductory message to the user.
	Config(config *Config) error

	// Progress is called to display a progress message.
	Progress(progress *Progress) error

	// OSV is invoked for each osv Entry in the stream.
	OSV(entry *osv.Entry) error

	// Finding is called for each vulnerability finding in the stream.
	Finding(finding *Finding) error
}

// HandleJSON reads the json from the supplied stream and hands the decoded
// output to the handler.
func HandleJSON(from io.Reader, to Handler) error {
	dec := json.NewDecoder(from)
	for dec.More() {
		msg := Message{}
		// decode the next message in the stream
		if err := dec.Decode(&msg); err != nil {
			return err
		}
		// dispatch the message
		var err error
		if msg.Config != nil {
			err = to.Config(msg.Config)
		}
		if msg.Progress != nil {
			err = to.Progress(msg.Progress)
		}
		if msg.OSV != nil {
			err = to.OSV(msg.OSV)
		}
		if msg.Finding != nil {
			err = to.Finding(msg.Finding)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
