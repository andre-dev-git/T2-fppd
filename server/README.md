# Game Server RPC

Servidor de jogo com RPC para sincronização de posições entre jogadores.

## 🚀 Executar
```bash
go run main.go [porta]
# Padrão: porta 8080
```

## 📡 Conectar Cliente
```go
client, err := rpc.Dial("tcp", "IP:8080")
// Ex: "192.168.1.100:8080" ou "localhost:8080"
```

## 🔧 Métodos RPC

### UpdatePosition
**Atualiza posição do jogador**
```go
type UpdatePositionRequest struct {
    PlayerID    int     `json:"player_id"`
    SequenceNum int     `json:"sequence_num"`
    X           float64 `json:"x"`
    Y           float64 `json:"y"`
}

type UpdatePositionResponse struct {
    Success     bool   `json:"success"`
    Message     string `json:"message"`
    SequenceNum int    `json:"sequence_num"`
}

// Chamada
client.Call("GameServer.UpdatePosition", req, &resp)
```

### GetPositions
**Retorna posições dos outros jogadores**
```go
type GetPositionsRequest struct {
    PlayerID    int `json:"player_id"`
    SequenceNum int `json:"sequence_num"`
}

type GetPositionsResponse struct {
    Success     bool      `json:"success"`
    Message     string    `json:"message"`
    SequenceNum int       `json:"sequence_num"`
    Players     []*Player `json:"players"`
}

type Player struct {
    ID int     `json:"id"`
    X  float64 `json:"x"`
    Y  float64 `json:"y"`
}

// Chamada
client.Call("GameServer.GetPositions", req, &resp)
```

## ⚠️ Importante
- **SequenceNum**: Previne comandos duplicados
- **PlayerID único**: Cada cliente deve usar ID diferente
- **Thread-safe**: Suporta múltiplos clientes simultâneos
