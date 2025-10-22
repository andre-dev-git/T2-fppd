package main

func personagemMover(tecla rune, jogo *Jogo, jogador int) {
    if jogo.fimJogo{ 
        return;
    }
	dx, dy := 0, 0
	switch tecla {
	case 'w': dy = -1
	case 'a': dx = -1
	case 's': dy = 1
	case 'd': dx = 1
	}

    nx, ny := jogo.Jog[jogador].PosX+dx, jogo.Jog[jogador].PosY+dy
	if jogoPodeMoverPara(jogo, nx, ny, true) {
		jogoMoverElemento(jogo, jogo.Jog[jogador].PosX, jogo.Jog[jogador].PosY, dx, dy)
		jogo.Jog[jogador].PosX, jogo.Jog[jogador].PosY = nx, ny
    }
}

func personagemInteragir(jogo *Jogo, jogador int) {
    direcoes := []struct {
        dx, dy int
    }{
        {0, -1},
        {0, 1},
        {-1, 0}, 
        {1, 0},
    }

    interagiu := false

    for _, dir := range direcoes {
        nx, ny := jogo.Jog[jogador].PosX+dir.dx, jogo.Jog[jogador].PosY+dir.dy
        if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
            elem := jogo.Mapa[ny][nx]

            if elem == Porta {
                if jogo.inventarioChave {
                    jogo.StatusMsg = "Você usou a chave para tentar abrir a porta!"
                    jogo.PortaCH <- true
                } else {
                    jogo.StatusMsg = "Você precisa de uma chave para abrir a porta."
                }
                interagiu = true
            } else if elem == Chave {
                jogo.StatusMsg = "Você pegou uma chave!"
                jogo.Mapa[ny][nx] = Vazio
                jogo.ChaveCH <- true
                jogo.inventarioChave = true
                interagiu = true
            } else if elem == Alavanca {
                jogo.StatusMsg = "Você acionou a alavanca!"
                jogo.AlavancaCH <- true 
                interagiu = true
            }
        }
    }

    if !interagiu {
        jogo.StatusMsg = "Nada para interagir ao seu redor."
    }
}

func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo, jogador int) bool {
	switch ev.Tipo {
	case "sair":
		return false
	case "interagir":
	case "mover":
		personagemMover(ev.Tecla, jogo, jogador)
	}
	return true
}
