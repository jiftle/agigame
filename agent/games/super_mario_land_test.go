package games

import (
	"testing"

	"agigame/core/gb"
)

// fakeMem is a scripted GameReader for pure-logic tests.
type fakeMem map[uint16]byte

func (m fakeMem) ReadMemory(addr uint16) byte { return m[addr] }
func (m fakeMem) Snapshot() gb.State          { return gb.State{} }

func set(f fakeMem, addr uint16, v int) {
	if v < 0 {
		f[addr] = 0xFF
		return
	}
	f[addr] = byte(v)
}

func baseState() fakeMem {
	f := fakeMem{
		addrCameraX:   10,
		addrMarioX:    50,
		addrMarioY:    144,
		addrMarioGround: 1,
		addrMarioFacing: 0x00, // right
		addrLives:      0x05,
		addrCoins:      0x10,
		addrTimeSec9:   0x03,
		addrTimeSec99:  0x25,
		addrStatus:     0x00, // small
		addrDemoStatus: 0x00, // in-game
	}
	return f
}

func mustExtract(t *testing.T, p *SMLPlugin, r GameReader) *MarioState {
	t.Helper()
	v, err := p.ExtractState(r)
	if err != nil {
		t.Fatalf("ExtractState: %v", err)
	}
	return v.(*MarioState)
}

func TestExtractState(t *testing.T) {
	p := &SMLPlugin{}
	s := mustExtract(t, p, baseState())

	if s.CameraX != 10 || s.MarioX != 50 || s.MarioY != 144 {
		t.Fatalf("positions: %+v", s)
	}
	if !s.OnGround {
		t.Fatalf("expected grounded")
	}
	if s.Lives != 5 || s.Coins != 10 || s.Time != 325 {
		t.Fatalf("lives/coins/time: %+v", s)
	}
	if s.State != "small" || s.Facing != "right" || s.Demo {
		t.Fatalf("status/facing/demo: %+v", s)
	}
}

func TestExtractStateEnemy(t *testing.T) {
	f := baseState()
	// OAM sprite stored as (screenX+8, screenY+16).
	f[0xFE00] = 144 + 16
	f[0xFE01] = 60 + 8
	f[0xFE02] = 0x11 // tile
	f[0xFE03] = 0x00 // attr
	p := &SMLPlugin{}
	s := mustExtract(t, p, f)
	if s.Enemy == nil {
		t.Fatalf("expected an enemy ahead")
	}
	if s.Enemy.DX != 10 { // 60-50
		t.Fatalf("enemy DX=%d want 10", s.Enemy.DX)
	}
}

func TestDecide(t *testing.T) {
	p := &SMLPlugin{}

	demo := baseState()
	demo[addrDemoStatus] = 0x28
	if got := p.Decide(mustExtract(t, p, demo)); !got.Start {
		t.Fatalf("demo should press Start, got %+v", got)
	}

	// Enemy within 16px → jump + right.
	withEnemy := baseState()
	withEnemy[0xFE01] = 60 + 8
	withEnemy[0xFE02] = 0x10
	s := mustExtract(t, p, withEnemy)
	b := p.Decide(s)
	if !b.A || !b.Right {
		t.Fatalf("near-enemy want A+Right, got %+v", b)
	}

	// Scouting, nothing near → just right.
	plain := baseState()
	b = p.Decide(mustExtract(t, p, plain))
	if b.Right && !b.A && !b.Start {
		return
	}
	t.Fatalf("scouting want Right only, got %+v", b)
}

func TestNeedsLLM(t *testing.T) {
	p := &SMLPlugin{}

	withEnemy := baseState()
	withEnemy[0xFE01] = 60 + 8
	if p.NeedsLLM(mustExtract(t, p, withEnemy), Buttons{}) {
		t.Fatalf("rules should handle nearby enemy")
	}

	air := baseState()
	air[addrMarioGround] = 0
	if p.NeedsLLM(mustExtract(t, p, air), Buttons{}) {
		t.Fatalf("rules should handle in-air")
	}

	if !p.NeedsLLM(mustExtract(t, p, baseState()), Buttons{}) {
		t.Fatalf("scouting state should consult the LLM")
	}
}

func TestReward(t *testing.T) {
	p := &SMLPlugin{}
	prev := mustExtract(t, p, baseState())
	cur := mustExtract(t, p, baseState()) // same camera

	// No change → skips (progress only).
	if ri := p.Reward(prev, cur); ri.Progress != 10 || ri.Delta != 0 {
		t.Fatalf("no-change reward: %+v", ri)
	}

	// Move forward.
	cur.CameraX = 30
	if ri := p.Reward(prev, cur); ri.Delta != 20 || len(ri.Events) != 0 {
		t.Fatalf("progress reward: %+v", ri)
	}

	// Death.
	cur2 := mustExtract(t, p, baseState())
	cur2.CameraX = 40
	cur2.Lives = 4
	ri := p.Reward(cur, cur2)
	if len(ri.Events) != 1 || ri.Events[0] != "death" || ri.Delta != -40 {
		t.Fatalf("death reward: %+v", ri)
	}

	// Game over.
	cur3 := mustExtract(t, p, baseState())
	cur3.CameraX = 50
	cur3.Lives = 0
	ri = p.Reward(cur2, cur3)
	if !ri.GameOver {
		t.Fatalf("want game over, got %+v", ri)
	}
}