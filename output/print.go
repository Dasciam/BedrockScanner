package output

import (
	"github.com/dasciam/bedrockscanner/scanner"
	"log"
	"net"
)

type Print struct{}

func (Print) Write(addr net.Addr, pong scanner.Pong) {
	log.Printf(
		"Server found!\n"+
			" | Address: %s\n"+
			" | MOTD: %s | %s\n"+
			" | Online: %d/%d\n"+
			" | Version: %s (%d)\n"+
			" | Game Mode: %s (%d)",
		addr,
		pong.MOTD, pong.SubMOTD,
		pong.OnlinePlayerCount, pong.MaxPlayerCount,
		pong.VersionString, pong.Protocol,
		pong.GameModeString, pong.GameModeNumber,
	)
}
