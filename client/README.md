# Client — Jogo de Terminal em Go (modo multiplayer via RPC)

Este client renderiza e executa **toda a lógica do jogo** (mapa, movimentos, inimigos, HUD) no terminal com `termbox-go`.  
A **sincronização multiplayer** é feita por RPC com um servidor que mantém **apenas o estado compartilhado** (posições dos jogadores).

> Controles do jogo (inalterados): **W/A/S/D** movem, **E** interage (local), **ESC** sai. :contentReference[oaicite:0]{index=0}

## Requisitos
- Go 1.21+ (ou superior)
- Windows, macOS ou Linux
- Terminal “puro” recomendado (no Windows, PowerShell ou Windows Terminal)  
- Arquivo `mapa.txt` presente na pasta do client (qualquer mapa válido)

## Como funciona (arquitetura)
- **Cliente (este projeto):** desenha interface, lê teclado, move jogador local, move inimigos, e **envia/recebe posições** por RPC.
- **Servidor:** não tem GUI nem mapa; apenas guarda as **posições dos jogadores** e responde a RPC.
- **Polling:** o client faz uma consulta periódica (~80ms) para obter **outros** jogadores conectados.
- **Exactly-once (escrita):** cada envio de posição usa um `SequenceNum` único; em falha de rede, o client reenvia **com o mesmo** número.

## Instalação
```bash
# dentro da pasta client
go mod tidy    # baixa termbox-go e dependências
