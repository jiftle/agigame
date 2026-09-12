package games

import (
	"fmt"
)

// ---------------------------------------------------------------------------
// WRAM / HRAM addresses for Super Mario Land.
//
// Sources: ROM Detectives wiki (Super_Mario_Land_(GB)_-_RAM) and the Super
// Mario Land disassembly on Data Crystal. Verified against a "SUPER MARIOLAND"
// (Rev.01, MBC1) dump.
//
// NOTE: addresses can differ between ROM revisions. If the rules misbehave on
// your dump, re-check these against your revision's RAM usage — everything
// below lives in this one file so it is trivial to adjust.
// ---------------------------------------------------------------------------

const (
	addrMarioY       = 0xC201 // Mario screen Y position
	addrMarioX       = 0xC202 // Mario screen X position
	addrMarioAnim    = 0xC203 // Mario animation frame
	addrMarioFacing  = 0xC205 // 0x00 right, 0x20 left
	addrMarioJump    = 0xC207 // mario jump state
	addrMarioJumpCtr = 0xC208 // jump state counter
	addrMarioGround  = 0xC20A // 0x00 air, 0x01 ground
	addrMarioMoveDir = 0xC20D // 01 turning, 10 right, 20 left
	addrMarioXSpeed  = 0xC20E // mario X speed

	addrTimeSplit = 0xDA00 // time split-seconds (00-28)
	addrTimeSec99 = 0xDA01 // time seconds low BCD digit
	addrTimeSec9  = 0xDA02 // time seconds high BCD digit
	addrLives     = 0xDA15 // lives (00-99, BCD)

	addrStatus     = 0xFF99 // mario status (00 small, 02 super)
	addrHardMode   = 0xFF9A // 00 normal, 01+ hard
	addrCameraX    = 0xFFA4 // camera X position
	addrDemoStatus = 0xFF9F // 00 in game, 28 demo mode
	addrCoins      = 0xFFFA // coins (00-99)
)

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

// EnemyInfo is the nearest threat detected ahead of Mario.
type EnemyInfo struct {
	ScreenX int // on-screen X
	ScreenY int // on-screen Y
	Tile    byte
	DX      int // horizontal distance ahead of Mario (negative = behind)
	DY      int // vertical distance (negative = above)
}

// MarioState is the extracted, per-frame game state.
type MarioState struct {
	CameraX   int
	MarioX    int
	MarioY    int
	OnGround  bool
	JumpState int
	Facing    string // "left" | "right"
	Lives     int
	Coins     int
	Time      int
	State     string // "small" | "super"
	HardMode  bool
	Demo      bool // attract mode / title screen
	Enemy     *EnemyInfo
	PitAhead  bool // NOTE: surface scan not implemented yet (TODO)
}

// Progress returns the run's absolute progress metric (camera X), used by the
// agent's safety net to detect stalls.
func (s *MarioState) Progress() int {
	if s == nil {
		return 0
	}
	return s.CameraX
}

// SMLPlugin is the game plugin for Super Mario Land.
type SMLPlugin struct{}

// ID implements games.GamePlugin.
func (p *SMLPlugin) ID() string { return "sml" }

// Name implements games.GamePlugin.
func (p *SMLPlugin) Name() string { return "Super Mario Land" }

// ExtractState implements games.GamePlugin.
func (p *SMLPlugin) ExtractState(r GameReader) (any, error) {
	s := &MarioState{
		CameraX:   int(r.ReadMemory(addrCameraX)),
		MarioX:    int(r.ReadMemory(addrMarioX)),
		MarioY:    int(r.ReadMemory(addrMarioY)),
		OnGround:  r.ReadMemory(addrMarioGround) != 0,
		JumpState: int(r.ReadMemory(addrMarioJump)),
		Lives:     bcd(r.ReadMemory(addrLives)),
		Coins:     bcd(r.ReadMemory(addrCoins)),
		Time:      bcd(r.ReadMemory(addrTimeSec9))*100 + bcd(r.ReadMemory(addrTimeSec99)),
		HardMode:  r.ReadMemory(addrHardMode) != 0,
		Demo:      r.ReadMemory(addrDemoStatus) != 0,
	}
	if f := r.ReadMemory(addrMarioFacing); f == 0x20 {
		s.Facing = "left"
	} else {
		s.Facing = "right"
	}
	switch st := r.ReadMemory(addrStatus); st {
	case 0x02:
		s.State = "super"
	default:
		s.State = "small"
	}
	s.Enemy = p.scanEnemies(r, s)
	return s, nil
}

// scanEnemies finds the nearest sprite ahead of Mario in OAM. This is a
// heuristic: it treats every non-Mario sprite ahead as a potential obstacle,
// without distinguishing enemy types yet (TODO: classify by tile id).
func (p *SMLPlugin) scanEnemies(r GameReader, s *MarioState) *EnemyInfo {
	var best *EnemyInfo
	for i := 0; i < 40; i++ {
		base := uint16(0xFE00 + i*4)
		oamX := int(r.ReadMemory(base + 1)) // stored as (screenX + 8)
		oamY := int(r.ReadMemory(base))     // stored as (screenY + 16)
		tile := r.ReadMemory(base + 2)
		attr := r.ReadMemory(base + 3) // unused for now
		_ = attr

		// Skip sprites whose on-screen origin matches Mario's own box:
		// those are Mario's sprite(s).
		sx, sy := oamX-8, oamY-16
		if abs(sx-s.MarioX) <= 2 && abs(sy-s.MarioY) <= 2 {
			continue
		}
		// Only consider sprites in front of Mario within the right radius.
		dx := sx - s.MarioX
		if s.Facing == "left" {
			dx = s.MarioX - sx
		}
		if dx < 0 || dx > 96 {
			continue
		}
		e := &EnemyInfo{
			ScreenX: sx,
			ScreenY: sy,
			Tile:    tile,
			DX:      dx,
			DY:      sy - s.MarioY,
		}
		if best == nil || e.DX < best.DX {
			best = e
		}
	}
	return best
}

// Decide implements games.GamePlugin. Rules strategy:
//   - title/demo screen → press Start
//   - in the air → keep holding right, do nothing else
//   - threat within 16px → jump
//   - pit ahead → jump (not yet detected, see TODO)
//   - otherwise hold right
func (p *SMLPlugin) Decide(state any) Buttons {
	s := state.(*MarioState)
	if s.Demo {
		return Buttons{Start: true}
	}
	if !s.OnGround {
		return Buttons{Right: true}
	}
	if s.Enemy != nil && s.Enemy.DX < 16 {
		return Buttons{A: true, Right: true}
	}
	if s.PitAhead {
		return Buttons{A: true}
	}
	return Buttons{Right: true}
}

// NeedsLLM implements games.GamePlugin. The rules are sufficient whenever
// something is clearly happening (jumping, dodging an enemy, demo screen);
// while casually scouting left to right the LLM may weigh in.
func (p *SMLPlugin) NeedsLLM(state any, _ Buttons) bool {
	s := state.(*MarioState)
	if s.Demo || !s.OnGround {
		return false
	}
	if s.Enemy != nil && s.Enemy.DX < 16 {
		return false
	}
	return true
}

// Prompt implements games.GamePlugin.
func (p *SMLPlugin) Prompt(state any) string {
	s := state.(*MarioState)
	return fmt.Sprintf(
		"You are an AI playing Super Mario Land."+
			"\nCurrent state: camera=%d mario=(x=%d,y=%d) on_ground=%v jump_state=%d"+
			" facing=%s lives=%d coins=%d time=%d state=%s hard_mode=%v demo=%v"+
			"\nThreat ahead: %v\nPit ahead: %v"+
			"\nChoose the next buttons as JSON: {\"commands\":[\"A\",\"Right\"],\"rationale\":\"...\"}."+
			" Available: A B Start Select Up Down Left Right (use A+Right style combos).",
		s.CameraX, s.MarioX, s.MarioY, s.OnGround, s.JumpState,
		s.Facing, s.Lives, s.Coins, s.Time, s.State, s.HardMode, s.Demo,
		s.Enemy, s.PitAhead)
}

// Reward implements games.GamePlugin. Reward = camera progress, penalising
// deaths and rewarding 1UPs.
func (p *SMLPlugin) Reward(prev, cur any) RewardInfo {
	before, ok1 := prev.(*MarioState)
	after, ok2 := cur.(*MarioState)
	info := RewardInfo{}
	if !ok1 || !ok2 || after.Demo || before == nil || after == nil {
		if after != nil {
			info.Progress = after.CameraX
		}
		return info
	}

	info.Progress = after.CameraX
	info.Delta = float64(after.CameraX-before.CameraX) * 1.0

	if after.Lives < before.Lives {
		info.Events = append(info.Events, "death")
		info.Delta -= 50
		if after.Lives == 0 {
			info.GameOver = true
			info.Events = append(info.Events, "game_over")
		}
	}
	if after.Lives > before.Lives {
		info.Events = append(info.Events, "one_up")
		info.Delta += 25
	}
	return info
}

// ExtraState implements games.GamePlugin.
func (p *SMLPlugin) ExtraState(cur any) map[string]any {
	s, ok := cur.(*MarioState)
	if !ok {
		return nil
	}
	m := map[string]any{
		"camera":   s.CameraX,
		"lives":    s.Lives,
		"coins":    s.Coins,
		"state":    s.State,
		"facing":   s.Facing,
		"ground":   s.OnGround,
		"time":     s.Time,
		"hard":     s.HardMode,
		"demo":     s.Demo,
		"pitAhead": s.PitAhead,
	}
	if s.Enemy != nil {
		m["enemyDX"] = s.Enemy.DX
	} else {
		m["enemyDX"] = nil
	}
	return m
}

// helpers ------------------------------------------------------------------

// bcd converts a packed-twos BCD byte to a decimal value.
func bcd(v byte) int {
	return int(v>>4)*10 + int(v&0x0F)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
