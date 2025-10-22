# T2-FPPD — Jogo Multiplayer em Go (Cliente/Servidor via RPC)

Jogo de terminal (renderização e lógica no client) com sincronização de posições de jogadores via RPC para um server simples (sem GUI, sem mapa). O cliente inicia todas as comunicações; o servidor apenas responde e mantém o estado compartilhado.

# Controles no client: W/A/S/D movem · E interage (local) · ESC sai.

# Requisitos
Go 1.21+ (Windows/macOS/Linux)
Terminal “puro” recomendado (PowerShell/Windows Terminal no Windows)
mapa.txt válido dentro de client/

# Estrutura do projeto
client/ — jogo em termbox-go (render, entrada, inimigos, lógica) + integração RPC
server/ — RPC em TCP com métodos GameServer.UpdatePosition e GameServer.GetPositions (estado de jogadores + logs)

```
├── 📁 client
│   ├── 📄 Makefile
│   ├── 📝 README.md
│   ├── 📄 build.bat
│   ├── 📄 go.mod
│   ├── 📄 go.sum
│   ├── 🐹 inimigo.go
│   ├── 🐹 interface.go
│   ├── 📄 jogo
│   ├── ⚙️ jogo.exe
│   ├── 🐹 jogo.go
│   ├── 🐹 main.go
│   ├── 📄 mapa.txt
│   ├── 🐹 personagem.go
│   └── 🐹 posicao.go
├── 📁 server
│   ├── 📁 shared
│   │   └── 🐹 types.go
│   ├── 📝 README.md
│   ├── 📄 build.bat
│   ├── 📄 go.mod
│   ├── 🐹 main.go
│   └── ⚙️ server.exe
└── 📝 README.md
```

# Como rodar
1) Iniciar o servidor

```bash
go run . 8080
# Padrão: porta 8080
```

2) Iniciar dois client (dois terminais/abas)
## 🚀 Executar

# Terminal A
```bash
go run .
```

# Terminal A
```bash
go run .
```

# Reiniciar a partida
Reset completo (recomendado): feche os dois clients (ESC), pare o server (Ctrl+C), suba o server e abra os dois clients novamente.
Reset rápido: apenas feche e reabra os clients (gera novos PlayerIDs). O server pode manter jogadores antigos até ser reiniciado.
