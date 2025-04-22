// personagem.go - Funções para movimentação e ações do personagem
package main

import "fmt"

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w': dy = -1
	case 'a': dx = -1
	case 's': dy = 1
	case 'd': dx = 1
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy

	// Verifica se há um inimigo no destino
	if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
		enemigo := jogo.Mapa[ny][nx]
		if enemigo.simbolo == Inimigo.simbolo || enemigo.simbolo == Alien.simbolo {
			jogo.Vida--
			jogo.StatusMsg = fmt.Sprintf("Você foi atingido! Vida restante: %d", jogo.Vida)
			if jogo.Vida <= 0 {
				jogo.StatusMsg = "💀 GAME OVER"
				return
			}
		}
	}

	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}


// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo) {
	// Atualmente apenas exibe uma mensagem de status
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo)
	}
	return true // Continua o jogo
}
