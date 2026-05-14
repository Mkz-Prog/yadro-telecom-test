package domain

import (
	"time"
)

// Config описывает структуру файла config.json
type Config struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

type PlayerState string

const (
	StateInProcess PlayerState = ""
	StateSuccess   PlayerState = "SUCCESS"
	StateFail      PlayerState = "FAIL"
	StateDisqual   PlayerState = "DISQUAL"
)

// Константы для входящих событий
const (
	EvRegistered       = 1
	EvEnteredDungeon   = 2
	EvKilledMonster    = 3
	EvNextFloor        = 4
	EvPrevFloor        = 5
	EvEnteredBossFloor = 6
	EvKilledBoss       = 7
	EvLeftDungeon      = 8
	EvCannotContinue   = 9
	EvHealthRestored   = 10
	EvDamageReceived   = 11
)

// Константы для исходящих событий
const (
	OutEvDisqualified   = 31
	OutEvDead           = 32
	OutEvImpossibleMove = 33
)

type Event struct {
	Time       time.Time
	TimeRaw    string
	ID         int
	PlayerID   int
	ExtraParam string
}

type Player struct {
	ID    int
	State PlayerState
	HP    int

	CurrentFloor   int
	MonstersKilled int

	DungeonEnterTime time.Time
	DungeonLeaveTime time.Time
	FloorEnterTime   time.Time
	BossEnterTime    time.Time

	TotalFloorTime time.Duration
	ClearedFloors  int
	BossKillTime   time.Duration
	TotalTimeSpent time.Duration

	IsRegistered bool
	HasEntered   bool
	HasLeft      bool
	IsDead       bool
}
