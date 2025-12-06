package data

import (
	"fmt"
	"github.com/dasciam/bedrockscanner/scanner"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"iter"
	"net"
	"time"
)

type ModelServer struct {
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

	OnlineFlag bool
	LastUpdate time.Time
}

type Base struct {
	db *gorm.DB
}

func New(file string) *Base {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}
	_ = db.AutoMigrate(ModelServer{})

	return &Base{db: db}
}

func (d *Base) Write(addr net.Addr, pong scanner.Pong) {
	d.db.Save(ModelServer{
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
		LastUpdate:        time.Now(),
		OnlineFlag:        true,
	})
}

func (d *Base) Read() iter.Seq2[string, scanner.Pong] {
	return func(yield func(string, scanner.Pong) bool) {
		var result []ModelServer
		if d.db.Find(&result).Error != nil {
			return
		}
		for _, r := range result {
			if !yield(r.Address, scanner.Pong{
				MOTD:              r.MOTD,
				Protocol:          r.Protocol,
				VersionString:     r.VersionString,
				OnlinePlayerCount: r.OnlinePlayerCount,
				MaxPlayerCount:    r.MaxPlayerCount,
				GUID:              r.GUID,
				SubMOTD:           r.SubMOTD,
				GameModeString:    r.GameModeString,
				GameModeNumber:    r.GameModeNumber,
				IPV4Port:          r.IPV4Port,
				IPV6Port:          r.IPV6Port,
			}) {
				return
			}
		}
	}
}

func (d *Base) FlagOffline() {
	d.db.Model(d.db).
		Update("online_flag", false)
}
