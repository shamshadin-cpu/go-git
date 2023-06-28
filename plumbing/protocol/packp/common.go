package packp

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/hash"
)

type stateFn func() stateFn

const (
	// common
	hashSize = 40

	// advrefs
	head   = "HEAD"
	noHead = "capabilities^{}"
)

var (
	// common
	sp  = []byte(" ")
	eol = []byte("\n")

	// advertised-refs
	null       = []byte("\x00")
	peeled     = []byte("^{}")
	noHeadMark = []byte(" capabilities^{}\x00")

	// upload-request
	want            = []byte("want ")
	shallow         = []byte("shallow ")
	deepen          = []byte("deepen")
	deepenCommits   = []byte("deepen ")
	deepenSince     = []byte("deepen-since ")
	deepenReference = []byte("deepen-not ")

	// shallow-update
	unshallow = []byte("unshallow ")

	// server-response
	ack = []byte("ACK")
	nak = []byte("NAK")

	// updreq
	shallowNoSp = []byte("shallow")
)

func isFlush(payload []byte) bool {
	return len(payload) == 0
}

func isEmpty(payload []byte) bool {
	if isFlush(payload) {
		return true
	}

	if len(payload) > hash.HexSize &&
		string(payload[:hash.HexSize]) == plumbing.ZeroHash.String() {
		return true
	}

	return false
}

// ErrUnexpectedData represents an unexpected data decoding a message
type ErrUnexpectedData struct {
	Msg  string
	Data []byte
}

// NewErrUnexpectedData returns a new ErrUnexpectedData containing the data and
// the message given
func NewErrUnexpectedData(msg string, data []byte) error {
	return &ErrUnexpectedData{Msg: msg, Data: data}
}

func (err *ErrUnexpectedData) Error() string {
	if len(err.Data) == 0 {
		return err.Msg
	}

	return fmt.Sprintf("%s (%s)", err.Msg, err.Data)
}
