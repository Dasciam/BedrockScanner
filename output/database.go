package output

import (
	"github.com/dasciam/bedrockscanner/scanner"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
)

type modelServer struct {
	Address string `gorm:"primaryKey"`

	MOTD              string
	CleanMOTD         string
	Protocol          int
	VersionString     string
	OnlinePlayerCount int
	MaxPlayerCount    int
	GUID              int64
	SubMOTD           string
	CleanSubMOTD      string
	GameModeString    string
	GameModeNumber    int
	IPV4Port          uint16
	IPV6Port          uint16
}

type Database struct {
	db *gorm.DB
}

func NewDatabase(file string) *Database {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	_ = db.AutoMigrate(modelServer{})

	return &Database{db: db}
}

func (d Database) Write(addr net.Addr, pong scanner.Pong) {
	d.db.Save(modelServer{
		Address:           addr.String(),
		MOTD:              pong.MOTD,
		CleanMOTD:         text.Clean(pong.MOTD),
		Protocol:          pong.Protocol,
		VersionString:     pong.VersionString,
		OnlinePlayerCount: pong.OnlinePlayerCount,
		MaxPlayerCount:    pong.MaxPlayerCount,
		GUID:              pong.GUID,
		SubMOTD:           pong.SubMOTD,
		CleanSubMOTD:      text.Clean(pong.SubMOTD),
		GameModeString:    pong.GameModeString,
		GameModeNumber:    pong.GameModeNumber,
		IPV4Port:          pong.IPV4Port,
		IPV6Port:          pong.IPV6Port,
	})
}
