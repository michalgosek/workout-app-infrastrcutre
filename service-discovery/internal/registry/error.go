package registry

import "fmt"

type Error struct {
	message   string
	code      int
	temporary bool
	network   bool
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Error() string {
	return fmt.Sprintf("msg: %s code: %d", e.message, e.code)
}

func (e Error) Temporary() bool {
	return e.temporary
}

func (e Error) IsNetwork() bool {
	return e.network
}

func IsTemporary(err error) bool {
	type temporary interface {
		Temporary() bool
	}
	te, ok := err.(temporary)
	return ok && te.Temporary()
}

func IsNetwork(err error) bool {
	type network interface {
		Network() bool
	}
	te, ok := err.(network)
	return ok && te.Network()
}
