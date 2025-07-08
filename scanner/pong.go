package scanner

import (
	"bytes"
	"log"
	"strconv"
)

type Pong struct {
	MOTD              string
	Protocol          int
	VersionString     string
	OnlinePlayerCount int
	MaxPlayerCount    int
	GUID              int64
	SubMOTD           string
	GameModeString    string
	GameModeNumber    int
	IPV4Port          uint16
	IPV6Port          uint16
}

func PongFromBytes(data []byte) (pong Pong, ok bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("PongFromBytes panic: %v", x)
		}
	}()

	seq := bytes.Split(data, []byte(";"))
	seq = seq[1:] // MCPE;

	var (
		protocol             = -1
		onlineCount          = -1
		maxCount             = -1
		guid           int64 = -1
		subMotd        string
		gameModeString string
		gameModeNumber = -1
		ipV4Port       uint16
		ipV6Port       uint16
	)
	protocol, _ = strconv.Atoi(string(seq[1]))
	onlineCount, _ = strconv.Atoi(string(seq[3]))
	maxCount, _ = strconv.Atoi(string(seq[4]))
	guid, _ = strconv.ParseInt(string(seq[5]), 10, 64)
	func() {
		// Some pongs can be from
		// servers with version lower
		// than 1.12, and they're
		// shit shit truncated.
		// Example: MCPE;Paola;82;0.15.7;1;5;18446744073564749294;
		defer func() {
			_ = recover()
			// Do nothing.
		}()
		subMotd = string(seq[6])
		gameModeString = string(seq[7])
		gameModeNumber, _ = strconv.Atoi(string(seq[8]))

		var (
			v4port uint64
			v6port uint64
		)
		v4port, _ = strconv.ParseUint(string(seq[9]), 10, 16)
		v6port, _ = strconv.ParseUint(string(seq[10]), 10, 16)

		ipV4Port = uint16(v4port)
		ipV6Port = uint16(v6port)
	}()

	return Pong{
		MOTD:              string(seq[0]),
		Protocol:          protocol,
		VersionString:     string(seq[2]),
		OnlinePlayerCount: onlineCount,
		MaxPlayerCount:    maxCount,
		GUID:              guid,
		SubMOTD:           subMotd,
		GameModeString:    gameModeString,
		GameModeNumber:    gameModeNumber,
		IPV4Port:          ipV4Port,
		IPV6Port:          ipV6Port,
	}, true
}
