package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"

	"t2-fppd-server/shared"
)

type GameServer struct {
	gameState     *shared.GameState
	mutex         sync.RWMutex
	processedCmds map[string]map[int]bool
	cmdMutex      sync.RWMutex
}

func NewGameServer() *GameServer {
	return &GameServer{
		gameState: &shared.GameState{
			Players: make(map[int]*shared.Player),
			NextID:  1,
		},
		processedCmds: make(map[string]map[int]bool),
	}
}

func (gs *GameServer) UpdatePosition(req shared.UpdatePositionRequest, response *shared.UpdatePositionResponse) error {
	fmt.Printf("[SERVIDOR] Recebido UpdatePosition: PlayerID=%d, SeqNum=%d, Pos=(%.2f,%.2f)\n",
		req.PlayerID, req.SequenceNum, req.X, req.Y)

	gs.cmdMutex.Lock()
	key := fmt.Sprintf("%d", req.PlayerID)
	if gs.processedCmds[key] == nil {
		gs.processedCmds[key] = make(map[int]bool)
	}
	if gs.processedCmds[key][req.SequenceNum] {
		gs.cmdMutex.Unlock()
		*response = shared.UpdatePositionResponse{
			Success:     true,
			Message:     "Comando já processado",
			SequenceNum: req.SequenceNum,
		}
		fmt.Printf("[SERVIDOR] Comando já processado para PlayerID=%d, SeqNum=%d\n",
			req.PlayerID, req.SequenceNum)
		return nil
	}
	gs.processedCmds[key][req.SequenceNum] = true
	gs.cmdMutex.Unlock()

	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if _, exists := gs.gameState.Players[req.PlayerID]; !exists {
		gs.gameState.Players[req.PlayerID] = &shared.Player{
			ID: req.PlayerID,
			X:  req.X,
			Y:  req.Y,
		}
		fmt.Printf("[SERVIDOR] Jogador %d criado na posição (%.2f, %.2f)\n",
			req.PlayerID, req.X, req.Y)
	} else {
		gs.gameState.Players[req.PlayerID].X = req.X
		gs.gameState.Players[req.PlayerID].Y = req.Y
		fmt.Printf("[SERVIDOR] Jogador %d movido para (%.2f, %.2f)\n",
			req.PlayerID, req.X, req.Y)
	}

	*response = shared.UpdatePositionResponse{
		Success:     true,
		Message:     "Posição atualizada com sucesso",
		SequenceNum: req.SequenceNum,
	}

	fmt.Printf("[SERVIDOR] Resposta UpdatePosition: Success=%t, SeqNum=%d\n",
		response.Success, response.SequenceNum)
	return nil
}

func (gs *GameServer) GetPositions(req shared.GetPositionsRequest, response *shared.GetPositionsResponse) error {
	fmt.Printf("[SERVIDOR] Recebido GetPositions: PlayerID=%d, SeqNum=%d\n",
		req.PlayerID, req.SequenceNum)

	gs.cmdMutex.Lock()
	key := fmt.Sprintf("%d", req.PlayerID)
	if gs.processedCmds[key] == nil {
		gs.processedCmds[key] = make(map[int]bool)
	}
	if gs.processedCmds[key][req.SequenceNum] {
		gs.cmdMutex.Unlock()
		*response = shared.GetPositionsResponse{
			Success:     true,
			Message:     "Comando já processado",
			SequenceNum: req.SequenceNum,
			Players:     []*shared.Player{},
		}
		fmt.Printf("[SERVIDOR] Comando já processado para PlayerID=%d, SeqNum=%d\n",
			req.PlayerID, req.SequenceNum)
		return nil
	}
	gs.processedCmds[key][req.SequenceNum] = true
	gs.cmdMutex.Unlock()

	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	players := make([]*shared.Player, 0, len(gs.gameState.Players))
	for _, player := range gs.gameState.Players {
		if player.ID != req.PlayerID {
			players = append(players, player)
		}
	}

	*response = shared.GetPositionsResponse{
		Success:     true,
		Message:     "Posições obtidas com sucesso",
		SequenceNum: req.SequenceNum,
		Players:     players,
	}

	fmt.Printf("[SERVIDOR] Resposta GetPositions: Success=%t, SeqNum=%d, %d colegas\n",
		response.Success, response.SequenceNum, len(players))
	return nil
}

func (gs *GameServer) Ping(empty string, response *string) error {
	*response = "pong"
	fmt.Printf("[SERVIDOR] Ping recebido, enviando pong\n")
	return nil
}

func main() {
	gameServer := NewGameServer()

	rpc.Register(gameServer)
	rpc.HandleHTTP()

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	address := ":" + port

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}

	fmt.Printf("Servidor RPC iniciado em %s\n", address)
	fmt.Printf("Aguardando conexões...\n")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v", err)
			continue
		}

		fmt.Printf("[SERVIDOR] Nova conexão aceita de %s\n", conn.RemoteAddr())
		go rpc.ServeConn(conn)
	}
}
