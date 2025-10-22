package main

import (
	"github.com/nsf/termbox-go"
)

type Cor = termbox.Attribute

const (
	CorPadrao     Cor = termbox.ColorDefault
	CorCinzaEscuro    = termbox.ColorDarkGray
	CorVermelho       = termbox.ColorRed
	CorVerde          = termbox.ColorGreen
	CorParede         = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede    = termbox.ColorDarkGray
	CorTexto          = termbox.ColorDarkGray
)

type EventoTeclado struct {
	Tipo  string
	Tecla rune
}

func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

func interfaceFinalizar() {
	termbox.Close()
}

func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()

	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	for _, inimigo := range jogo.InimigosTipo1 {
  		interfaceDesenharElemento(inimigo.X, inimigo.Y, Inimigo)
	}

	for _, inimigo := range jogo.InimigosTipo2 {
    	interfaceDesenharElemento(inimigo.X, inimigo.Y, Inimigo2)
	}

	for _, inimigo := range jogo.InimigosTipo3 {
    	interfaceDesenharElemento(inimigo.X, inimigo.Y, Inimigo3)
	}

	interfaceDesenharElemento(jogo.Jog[0].PosX, jogo.Jog[0].PosY, Personagem)
	interfaceDesenharElemento(jogo.Jog[1].PosX, jogo.Jog[1].PosY, Personagem2)

	interfaceDesenharBarraDeStatus(jogo)

	interfaceAtualizarTela()
}

func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

func interfaceAtualizarTela() {
	termbox.Flush()
}

func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
	}

	msg := "Use WASD para mover e E para interagir. ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorPadrao)
	}
}
