package main

type inimigo struct {
    X, Y       int
    Simbolo    rune
    Direcao    string
    Movimento  func(*Jogo, *inimigo)
}

func novoInimigo(x, y int, simbolo rune, direcao string, movimento func(*Jogo, *inimigo)) inimigo {
    return inimigo{
        X:         x,
        Y:         y,
        Simbolo:   simbolo,
        Direcao:   direcao,
        Movimento: movimento,
    }
}

func inimigoMover(jogo *Jogo, inimigo *inimigo) {
    
    if inimigo.Movimento != nil {
        inimigo.Movimento(jogo, inimigo)
    }
}

func movimentoQuadrado(jogo *Jogo, inimigo *inimigo) {
    switch inimigo.Direcao {
    case "baixo":
        if jogoPodeMoverPara(jogo, inimigo.X, inimigo.Y+1, false) {
            jogoMoverElemento(jogo, inimigo.X, inimigo.Y, 0, 1)
            inimigo.Y++
        } else {
            inimigo.Direcao = "direita"
        }
    case "direita":
        if jogoPodeMoverPara(jogo, inimigo.X+1, inimigo.Y, false) {
            jogoMoverElemento(jogo, inimigo.X, inimigo.Y, 1, 0)
            inimigo.X++
        } else {
            inimigo.Direcao = "cima"
        }
    case "cima":
        if jogoPodeMoverPara(jogo, inimigo.X, inimigo.Y-1, false) {
            jogoMoverElemento(jogo, inimigo.X, inimigo.Y, 0, -1)
            inimigo.Y--
        } else {
            inimigo.Direcao = "esquerda"
        }
    case "esquerda":
        if jogoPodeMoverPara(jogo, inimigo.X-1, inimigo.Y, false) {
            jogoMoverElemento(jogo, inimigo.X, inimigo.Y, -1, 0)
            inimigo.X--
        } else {
            inimigo.Direcao = "baixo"
        }
    }
}

func movimentoCimaBaixo(jogo *Jogo, inimigo *inimigo) {
    if inimigo.Direcao == "cima" {
        if jogoPodeMoverPara(jogo, inimigo.X, inimigo.Y-1, false) {
            jogoMoverElemento(jogo, inimigo.X, inimigo.Y, 0, -1)
            inimigo.Y--
        } else {
            inimigo.Direcao = "baixo"
        }
    } else if inimigo.Direcao == "baixo" {
        if jogoPodeMoverPara(jogo, inimigo.X, inimigo.Y+1, false) {
            jogoMoverElemento(jogo, inimigo.X, inimigo.Y, 0, 1)
            inimigo.Y++
        } else {
            inimigo.Direcao = "cima"
        }
    }
}

func movimentoEsquerdaDireita(jogo *Jogo, inimigo *inimigo) {
    if inimigo.Direcao == "direita" {
        if jogoPodeMoverPara(jogo, inimigo.X+1, inimigo.Y, false) {
            jogoMoverElemento(jogo, inimigo.X, inimigo.Y, 1, 0)
            inimigo.X++
        } else {
            inimigo.Direcao = "esquerda"
        }
    } else if inimigo.Direcao == "esquerda" {
        if jogoPodeMoverPara(jogo, inimigo.X-1, inimigo.Y, false) {
            jogoMoverElemento(jogo, inimigo.X, inimigo.Y, -1, 0)
            inimigo.X--
        } else {
            inimigo.Direcao = "direita"
        }
    }
}
