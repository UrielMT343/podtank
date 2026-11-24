package gameWS

import (
	"context"
	"sync"

	"github.com/UrielMT343/podtank/internal/domain/arena"
	"github.com/UrielMT343/podtank/internal/domain/tank"
	"github.com/coder/websocket"
)

type PlayerInput struct {
	Player tank.Player
	Cmd    tank.CmdState
}

type PlayerSession struct {
	Role   tank.Player
	Conn   *websocket.Conn
	Send   chan []byte
	Ctx    context.Context
	Cancel context.CancelFunc
}

type GameRoom struct {
	ID          string
	Battlefield *arena.Arena
	Tanks       map[tank.Player]*tank.Tank
	BulletPools map[tank.Player]*tank.BulletPool
	Players     map[tank.Player]*PlayerSession
	Inputs      chan PlayerInput
	Ctx         context.Context
	Cancel      context.CancelFunc
}

type RoomManager struct {
	mu      sync.Mutex
	waiting *PlayerSession
	rooms   map[string]*GameRoom
}

func NewRoomManager() *RoomManager {
	r := make(map[string]*GameRoom)
	return &RoomManager{
		rooms: r,
	}
}

func (rm *RoomManager) JoinOrCreate(p *PlayerSession) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.waiting == nil {
		p.Role = tank.Player1
		rm.waiting = p
		return
	} else {
		p.Role = tank.Player2
	}
}
