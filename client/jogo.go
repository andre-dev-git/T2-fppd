package main

import (
	"bufio"
	"os"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool
}

type Jogo struct {
	Mapa           [][]Elemento
	Jog            []Posicao
	UltimoVisitado Elemento
	StatusMsg      string
	InimigosTipo1  []inimigo
	InimigosTipo2  []inimigo
	InimigosTipo3  []inimigo
	AlavancaCH     chan bool
	ChaveCH        chan bool
	PortaCH        chan bool
	sync.Mutex
	inventarioChave bool
	lavaEmExecucao  bool
	fimJogo         bool
}

var (
	Personagem  = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Personagem2 = Elemento{'☻', CorCinzaEscuro, CorPadrao, true}
	Inimigo     = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede      = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao   = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio       = Elemento{' ', CorPadrao, CorPadrao, false}

	Inimigo2       = Elemento{'☢', CorVermelho, CorPadrao, true}
	Inimigo3       = Elemento{'☣', CorVermelho, CorPadrao, true}
	Alavanca       = Elemento{'♠', CorPadrao, CorPadrao, true}
	Chave          = Elemento{'⚿', CorPadrao, CorPadrao, true}
	Porta          = Elemento{'⍈', CorPadrao, CorPadrao, true}
	Lava           = Elemento{'■', CorVermelho, CorVermelho, false}
	Impedimento    = Elemento{'-', CorVermelho, CorPadrao, true}
	SemImpedimento = Elemento{'~', CorVermelho, CorPadrao, false}
)

func jogoNovo() Jogo {
	return Jogo{
		UltimoVisitado:  Vazio,
		AlavancaCH:      make(chan bool, 1),
		ChaveCH:         make(chan bool, 1),
		PortaCH:         make(chan bool, 1),
		inventarioChave: false,
		lavaEmExecucao:  false,
		fimJogo:         false,
	}
}

func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Inimigo.simbolo:
				inimigo := novoInimigo(x, y, Inimigo.simbolo, "esquerda", movimentoQuadrado)
				jogo.InimigosTipo1 = append(jogo.InimigosTipo1, inimigo)
			case Inimigo2.simbolo:
				inimigo := novoInimigo(x, y, Inimigo2.simbolo, "cima", movimentoCimaBaixo)
				jogo.InimigosTipo2 = append(jogo.InimigosTipo2, inimigo)
			case Inimigo3.simbolo:
				inimigo := novoInimigo(x, y, Inimigo3.simbolo, "esquerda", movimentoEsquerdaDireita)
				jogo.InimigosTipo3 = append(jogo.InimigosTipo3, inimigo)
			case Parede.simbolo:
				e = Parede
			case Vegetacao.simbolo:
				e = Vegetacao
			case Alavanca.simbolo:
				e = Alavanca
			case Impedimento.simbolo:
				e = Impedimento
			case SemImpedimento.simbolo:
				e = SemImpedimento
			case Chave.simbolo:
				e = Chave
			case Porta.simbolo:
				e = Porta
			case Lava.simbolo:
				e = Lava
			case Personagem.simbolo:
				jogo.Jog = append(jogo.Jog, Posicao{PosX: x, PosY: y})
			case Personagem2.simbolo:
				jogo.Jog = append(jogo.Jog, Posicao{PosX: x, PosY: y})
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func jogoPodeMoverPara(jogo *Jogo, x, y int, isPlayer bool) bool {
	Destino := Posicao{PosX: x, PosY: y} //adaptado do COPILOT

	if y < 0 || y >= len(jogo.Mapa) || x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	if isPlayer {
		elem := jogo.Mapa[Destino.PosY][Destino.PosX]
		if elem == Inimigo || elem == Inimigo2 || elem == Inimigo3 || elem == Lava {
			jogoFim(jogo)
			return false
		}
	} else {
		if Destino == jogo.Jog[0] || Destino == jogo.Jog[1] {
			jogoFim(jogo)
			return false
		}
	}

	if jogo.Mapa[Destino.PosY][Destino.PosX].tangivel {
		return false
	}

	return true
}

func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	elemento := jogo.Mapa[y][x]
	jogo.Mapa[y][x] = jogo.UltimoVisitado
	jogo.UltimoVisitado = jogo.Mapa[ny][nx]
	jogo.Mapa[ny][nx] = elemento
}

func jogoMoverInimigosTipo1(jogo *Jogo) {
	for {
		for i := range jogo.InimigosTipo1 {
			inimigoMover(jogo, &jogo.InimigosTipo1[i])
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func jogoMoverInimigosTipo2(jogo *Jogo) {
	for {
		for i := range jogo.InimigosTipo2 {
			inimigoMover(jogo, &jogo.InimigosTipo2[i])
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func jogoMoverInimigosTipo3(jogo *Jogo) {
	for {
		for i := range jogo.InimigosTipo3 {
			inimigoMover(jogo, &jogo.InimigosTipo3[i])
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func fazerLava(jogo *Jogo) {
	var x []int
	var y []int

	for linha := range jogo.Mapa {
		for coluna, elem := range jogo.Mapa[linha] {
			if elem == Impedimento {
				x = append(x, coluna)
				y = append(y, linha)
			}
		}
	}

	for {
		select {
		case <-jogo.AlavancaCH:
			{
				jogo.lavaEmExecucao = true
				for i := 0; i < len(x); i++ {
					for linha := y[i] + 1; linha < len(jogo.Mapa); linha++ {
						time.Sleep(1 * time.Second)
						posLava := Posicao{PosX: x[i], PosY: linha}
						if (jogo.Jog[0] == posLava) || (jogo.Jog[1] == posLava) {
							jogoFim(jogo)
						}
						jogo.Mapa[linha][x[i]] = Lava
						interfaceDesenharJogo(jogo)
					}
				}
			}
		case <-jogo.ChaveCH:
			{
				jogo.lavaEmExecucao = true
				for i := 0; i < len(x); i++ {
					for linha := y[i] + 1; linha < len(jogo.Mapa); linha++ {
						time.Sleep(3 * time.Second)
						posLava := Posicao{PosX: x[i], PosY: linha}
						if (jogo.Jog[0] == posLava) || (jogo.Jog[1] == posLava) {
							jogoFim(jogo)
						}
						jogo.Mapa[linha][x[i]] = Lava
						interfaceDesenharJogo(jogo)
					}
				}
			}
		default:
			jogo.lavaEmExecucao = false
		}
	}
}

func gerenciarPorta(jogo *Jogo) {
	for {
		select {
		case <-jogo.PortaCH:
			jogo.StatusMsg = "Você abriu a porta e venceu o jogo!"
			interfaceDesenharJogo(jogo)
			time.Sleep(3 * time.Second)
			jogoFim(jogo)
		case <-time.After(180 * time.Second):
			jogo.StatusMsg = "Tempo esgotado! Você não conseguiu abrir a porta."
			time.Sleep(3 * time.Second)
			jogoFim(jogo)
			interfaceDesenharJogo(jogo)
		}
	}
}

func gerenciarImpedimentos(jogo *Jogo) {
	for {
		<-jogo.AlavancaCH
		jogo.Lock()
		for y, linha := range jogo.Mapa {
			for x, elem := range linha {
				if elem == Impedimento {
					jogo.Mapa[y][x] = SemImpedimento
				}
			}
		}
		interfaceDesenharJogo(jogo)
		jogo.Unlock()

		time.Sleep(15 * time.Second)

		jogo.Lock()
		for y, linha := range jogo.Mapa {
			for x, elem := range linha {
				if elem == SemImpedimento && jogo.Mapa[y][x] == SemImpedimento {
					jogo.Mapa[y][x] = Impedimento
				}
			}
		}
		interfaceDesenharJogo(jogo)
		jogo.Unlock()
		if !jogo.lavaEmExecucao {
			jogo.AlavancaCH <- true
		}
		time.Sleep(1 * time.Second)
	}
}

func jogoAtualizarTela(jogo *Jogo) {
	for {
		interfaceDesenharJogo(jogo)
		time.Sleep(100 * time.Millisecond)
	}
}

func jogoFim(jogo *Jogo) {
	jogo.fimJogo = true
	jogo.StatusMsg = "Fim de Jogo."
	time.Sleep(5 * time.Second)
	termbox.Close()
	os.Exit(0)
}
