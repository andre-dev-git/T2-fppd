package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

type rpcPlayer struct {
	ID int
	X  float64
	Y  float64
}

type rpcUpdatePosReq struct {
	PlayerID    int
	SequenceNum int
	X           float64
	Y           float64
}

type rpcUpdatePosResp struct {
	Success     bool
	Message     string
	SequenceNum int
}

type rpcGetPosReq struct {
	PlayerID    int
	SequenceNum int
}

type rpcGetPosResp struct {
	Success     bool
	Message     string
	SequenceNum int
	Players     []*rpcPlayer
}

func init() {
	gob.RegisterName("shared.Player", rpcPlayer{})
	gob.RegisterName("shared.UpdatePositionRequest", rpcUpdatePosReq{})
	gob.RegisterName("shared.UpdatePositionResponse", rpcUpdatePosResp{})
	gob.RegisterName("shared.GetPositionsRequest", rpcGetPosReq{})
	gob.RegisterName("shared.GetPositionsResponse", rpcGetPosResp{})
}

type rpcShim struct {
	c        *rpc.Client
	playerID int32
	seq      int32
}

func newRPCShim(addr string, playerID int) (*rpcShim, error) {
	c, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &rpcShim{c: c, playerID: int32(playerID), seq: 0}, nil
}

func (s *rpcShim) nextSeq() int {
	return int(atomic.AddInt32(&s.seq, 1))
}

func retry(attempts int, backoff time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(backoff)
	}
	return err
}

func (s *rpcShim) SendPosition(px, py int) error {
	seq := s.nextSeq()
	req := rpcUpdatePosReq{
		PlayerID:    int(s.playerID),
		SequenceNum: seq,
		X:           float64(px),
		Y:           float64(py),
	}
	var resp rpcUpdatePosResp
	return retry(3, 150*time.Millisecond, func() error {
		return s.c.Call("GameServer.UpdatePosition", req, &resp)
	})
}

func (s *rpcShim) PollOthers() ([]*rpcPlayer, error) {
	seq := s.nextSeq()
	req := rpcGetPosReq{PlayerID: int(s.playerID), SequenceNum: seq}
	var resp rpcGetPosResp
	if err := retry(3, 150*time.Millisecond, func() error {
		return s.c.Call("GameServer.GetPositions", req, &resp)
	}); err != nil {
		return nil, err
	}
	return resp.Players, nil
}

func (s *rpcShim) Close() error { return s.c.Close() }

func defaultAddr(host string, port int) string { return fmt.Sprintf("%s:%d", host, port) }

func generateClientPlayerID() int {
	return int(time.Now().UnixNano() & 0x3fffffff)
}

func main() {

	host := "localhost"
	if len(os.Args) > 1 && os.Args[1] != "" {
		host = os.Args[1]
	}
	port := 8080
	if len(os.Args) > 2 && os.Args[2] != "" {
		if p, err := strconv.Atoi(os.Args[2]); err == nil {
			port = p
		}
	}
	mapaFile := "mapa.txt"
	if len(os.Args) > 3 && os.Args[3] != "" {
		mapaFile = os.Args[3]
	}

	interfaceIniciar()
	defer interfaceFinalizar()

	j := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &j); err != nil {
		log.Fatalf("Erro ao carregar mapa %q: %v", mapaFile, err)
	}

	if len(j.Jog) < 1 {
		j.Jog = append(j.Jog, Posicao{})
	}

	playerID := generateClientPlayerID()
	addr := defaultAddr(host, port)

	shim, err := newRPCShim(addr, playerID)
	if err != nil {
		j.StatusMsg = fmt.Sprintf("Erro ao conectar ao servidor (%s): %v", addr, err)
		interfaceDesenharJogo(&j)
		return
	}
	defer shim.Close()

	j.StatusMsg = fmt.Sprintf("Conectado a %s | PlayerID=%d | Mapa=%s", addr, playerID, mapaFile)
	interfaceDesenharJogo(&j)

	go jogoMoverInimigosTipo1(&j)
	go jogoMoverInimigosTipo2(&j)
	go jogoMoverInimigosTipo3(&j)

	go gerenciarImpedimentos(&j)
	go gerenciarPorta(&j)
	go fazerLava(&j)

	go jogoAtualizarTela(&j)

	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			others, err := shim.PollOthers()
			if err != nil {
				continue
			}

			j.Lock()
			if len(j.Jog) < 1 {
				j.Jog = append(j.Jog, Posicao{})
			}
			j.Jog = j.Jog[:1]

			for _, op := range others {
				if op == nil || op.ID == playerID {
					continue
				}
				j.Jog = append(j.Jog, Posicao{PosX: int(op.X), PosY: int(op.Y)})
			}
			j.Unlock()
		}
	}()

	for {
		ev := interfaceLerEventoTeclado()
		switch ev.Tipo {
		case "sair":
			return

		case "interagir":
			personagemInteragir(&j, 0)

		case "mover":
			antes := j.Jog[0]
			personagemExecutarAcao(ev, &j, 0)
			depois := j.Jog[0]
			if antes != depois {
				_ = shim.SendPosition(depois.PosX, depois.PosY)
			}
		default:
		}
		interfaceDesenharJogo(&j)
	}
}
