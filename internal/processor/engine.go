package processor

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/Mkz-Prog/yadro-telecom-test/internal/domain"
)

// Engine управляет состоянием всех игроков и обрабатывает входящие события.
type Engine struct {
	Config    *domain.Config
	Players   map[int]*domain.Player
	OpenTime  time.Time
	CloseTime time.Time
}

// NewEngine инициализирует движок.
func NewEngine(cfg *domain.Config) *Engine {
	openTime, _ := time.Parse("15:04:05", cfg.OpenAt)

	closeTime := openTime.Add(time.Duration(cfg.Duration) * time.Hour)

	return &Engine{
		Config:    cfg,
		Players:   make(map[int]*domain.Player),
		OpenTime:  openTime,
		CloseTime: closeTime,
	}
}

// ProcessEvent применяет одно событие к состоянию конкретного игрока.
func (e *Engine) ProcessEvent(event *domain.Event) {
	if event.Time.Before(e.OpenTime) {
		event.Time = event.Time.Add(24 * time.Hour)
	}

	player, exists := e.Players[event.PlayerID]
	if !exists {
		player = &domain.Player{
			ID:    event.PlayerID,
			State: domain.StateInProcess,
			HP:    100,
		}
		e.Players[event.PlayerID] = player
	}

	if !player.IsRegistered && event.ID != domain.EvRegistered {
		if player.State != domain.StateDisqual {
			player.State = domain.StateDisqual
			fmt.Printf("%s Player [%d] is disqualified\n", event.TimeRaw, player.ID)
		}
		return
	}

	if player.State == domain.StateDisqual || player.IsDead || player.HasLeft {
		return
	}

	switch event.ID {
	case domain.EvRegistered:
		player.IsRegistered = true
		player.HP = 100
		fmt.Printf("%s Player [%d] registered\n", event.TimeRaw, player.ID)

	case domain.EvEnteredDungeon:
		player.HasEntered = true
		player.DungeonEnterTime = event.Time
		player.FloorEnterTime = event.Time
		player.CurrentFloor = 1
		fmt.Printf("%s Player [%d] entered the dungeon\n", event.TimeRaw, player.ID)

	case domain.EvKilledMonster:
		if !player.HasEntered || player.CurrentFloor >= e.Config.Floors || player.MonstersKilled >= e.Config.Monsters {
			fmt.Printf("%s Player [%d] makes imposible move [%d]\n", event.TimeRaw, player.ID, event.ID)
			return
		}
		player.MonstersKilled++
		fmt.Printf("%s Player [%d] killed the monster\n", event.TimeRaw, player.ID)

		if player.MonstersKilled == e.Config.Monsters {
			player.TotalFloorTime += event.Time.Sub(player.FloorEnterTime)
			player.ClearedFloors++
		}

	case domain.EvNextFloor:
		if player.CurrentFloor >= e.Config.Floors {
			fmt.Printf("%s Player [%d] makes imposible move [%d]\n", event.TimeRaw, player.ID, event.ID)
			return
		}
		player.CurrentFloor++
		player.FloorEnterTime = event.Time
		player.MonstersKilled = 0 // Сбрасываем счетчик монстров для нового этажа
		fmt.Printf("%s Player [%d] went to the next floor\n", event.TimeRaw, player.ID)

	case domain.EvPrevFloor:
		if player.CurrentFloor <= 1 {
			fmt.Printf("%s Player [%d] makes imposible move [%d]\n", event.TimeRaw, player.ID, event.ID)
			return
		}
		player.CurrentFloor--
		player.FloorEnterTime = event.Time
		player.MonstersKilled = 0
		fmt.Printf("%s Player [%d] went to the previous floor\n", event.TimeRaw, player.ID)

	case domain.EvEnteredBossFloor:
		if player.CurrentFloor != e.Config.Floors {
			fmt.Printf("%s Player [%d] makes imposible move [%d]\n", event.TimeRaw, player.ID, event.ID)
			return
		}
		player.BossEnterTime = event.Time
		fmt.Printf("%s Player [%d] entered the boss's floor\n", event.TimeRaw, player.ID)

	case domain.EvKilledBoss:
		if player.CurrentFloor != e.Config.Floors {
			fmt.Printf("%s Player [%d] makes imposible move [%d]\n", event.TimeRaw, player.ID, event.ID)
			return
		}
		player.BossKillTime = event.Time.Sub(player.BossEnterTime)
		fmt.Printf("%s Player [%d] killed the boss\n", event.TimeRaw, player.ID)

		if player.ClearedFloors == e.Config.Floors-1 {
			player.State = domain.StateSuccess
		}

	case domain.EvLeftDungeon:
		player.HasLeft = true
		player.DungeonLeaveTime = event.Time
		fmt.Printf("%s Player [%d] left the dungeon\n", event.TimeRaw, player.ID)

		if player.State != domain.StateSuccess {
			player.State = domain.StateFail
		}

	case domain.EvCannotContinue:
		player.HasLeft = true
		player.State = domain.StateFail
		player.DungeonLeaveTime = event.Time
		fmt.Printf("%s Player [%d] cannot continue due to [%s]\n", event.TimeRaw, player.ID, event.ExtraParam)

	case domain.EvHealthRestored:
		heal, _ := strconv.Atoi(event.ExtraParam)
		player.HP += heal
		if player.HP > 100 {
			player.HP = 100
		}
		fmt.Printf("%s Player [%d] has restored [%d] of health\n", event.TimeRaw, player.ID, heal)

	case domain.EvDamageReceived:
		dmg, _ := strconv.Atoi(event.ExtraParam)
		player.HP -= dmg
		fmt.Printf("%s Player [%d] recieved [%d] of damage\n", event.TimeRaw, player.ID, dmg)

		if player.HP <= 0 {
			player.HP = 0
			player.IsDead = true
			player.State = domain.StateFail
			player.DungeonLeaveTime = event.Time
			fmt.Printf("%s Player [%d] is dead\n", event.TimeRaw, player.ID)
		}
	}
}

// PrintFinalReport выводит итоговую статистику по всем участникам.
func (e *Engine) PrintFinalReport() {
	fmt.Println("\nFinal report:")

	var keys []int
	for id := range e.Players {
		keys = append(keys, id)
	}
	sort.Ints(keys)

	for _, id := range keys {
		p := e.Players[id]

		if p.State == domain.StateInProcess {
			p.State = domain.StateFail
		}

		if p.HasEntered {
			if p.DungeonLeaveTime.IsZero() {
				p.TotalTimeSpent = e.CloseTime.Sub(p.DungeonEnterTime)
			} else {
				p.TotalTimeSpent = p.DungeonLeaveTime.Sub(p.DungeonEnterTime)
			}
		}

		var avgFloorTime time.Duration
		if p.ClearedFloors > 0 {
			avgFloorTime = time.Duration(int64(p.TotalFloorTime) / int64(p.ClearedFloors))
		}

		fmt.Printf("[%s] %d [%s, %s, %s] HP:%d\n",
			p.State,
			p.ID,
			formatDuration(p.TotalTimeSpent),
			formatDuration(avgFloorTime),
			formatDuration(p.BossKillTime),
			p.HP,
		)
	}
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
