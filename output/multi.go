package output

import (
	"github.com/dasciam/bedrockscanner/scanner"
	"net"
)

type Multi struct {
	outputs []scanner.Output
}

func NewMulti(outputs ...scanner.Output) *Multi {
	return &Multi{outputs: outputs}
}

func (m Multi) Write(addr net.Addr, pong scanner.Pong) {
	for _, output := range m.outputs {
		output.Write(addr, pong)
	}
}
