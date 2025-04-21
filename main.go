// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"
)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// No main.go, dentro do loop de inicialização
for i := range jogo.Inimigos {
	inimigo := &jogo.Inimigos[i]
	go func(inimigo *InimigoMovel) {
		for {
			moverInimigo(inimigo, &jogo)
			interfaceDesenharJogo(&jogo) // Atualiza a tela após o movimento
			time.Sleep(300 * time.Millisecond)
		}
	}(inimigo)
}

// Lança uma goroutine para cada alien (movimento vertical)
for i := range jogo.Aliens {
	alien := &jogo.Aliens[i]
	go func(alien *AlienMovel) {
		for {
			moverAlien(alien, &jogo)
			interfaceDesenharJogo(&jogo) // Atualiza a tela após o movimento
			time.Sleep(300 * time.Millisecond)
		}
	}(alien)
}


	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo) // Atualiza a tela após a ação do jogador
	}
}
