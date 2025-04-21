// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
	"sync"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo   rune
	cor       Cor
	corFundo  Cor
	tangivel  bool // Indica se o elemento bloqueia passagem
}

type InimigoMovel struct {
	X, Y     int
	Direita  bool
}
type AlienMovel struct {
	X, Y     int
	Subindo  bool
}
// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa            [][]Elemento // grade 2D representando o mapa
	PosX, PosY      int          // posição atual do personagem
	UltimoVisitado  Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg       string       // mensagem para a barra de status
	Inimigos       []InimigoMovel // inimigos móveis
	Aliens			[]AlienMovel // aliens móveis
	Mutex          sync.Mutex     // para proteger o mapa
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
	Tiro	  = Elemento{'*', CorAmarelo, CorPadrao, true}
	Boss	  = Elemento{'♡', CorVermelho, CorPadrao, true}
	Explosao  = Elemento{'*', CorVermelho, CorPadrao, true}
	Radiativo  = Elemento{'☢', CorVerde, CorPadrao, true}
	Alien = Elemento{'Ψ', CorCiano, CorPadrao, true}

)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{UltimoVisitado: Vazio}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
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
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				jogo.Inimigos = append(jogo.Inimigos, InimigoMovel{
					X: x, Y: y, Direita: true,
				})
				e = Vazio // Remove do mapa estático, para a goroutine cuidar do desenho			
			case Alien.simbolo:
				jogo.Aliens = append(jogo.Aliens, AlienMovel{
					X: x, Y: y, Subindo: true,
				})
				e = Vazio
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
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

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado     // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx]   // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento              // move o elemento
}

func moverInimigo(inimigo *InimigoMovel, jogo *Jogo) {
	jogo.Mutex.Lock() // Garantir que o mapa não será alterado por outra goroutine enquanto esse inimigo se move
	defer jogo.Mutex.Unlock()

	dx := 1
	if !inimigo.Direita {
		dx = -1
	}
	nx := inimigo.X + dx
	ny := inimigo.Y

	if nx < 0 || nx >= len(jogo.Mapa[0]) {
		inimigo.Direita = !inimigo.Direita
		return
	}

	destino := jogo.Mapa[ny][nx]
	if destino.tangivel || (jogo.PosX == nx && jogo.PosY == ny) {
		inimigo.Direita = !inimigo.Direita
		return
	}

	// Move o inimigo
	jogo.Mapa[inimigo.Y][inimigo.X] = Vazio
	jogo.Mapa[ny][nx] = Inimigo
	inimigo.X = nx
	inimigo.Y = ny
}

func moverAlien(alien *AlienMovel, jogo *Jogo) {
	jogo.Mutex.Lock() // Garantir que o mapa não será alterado por outra goroutine enquanto esse inimigo se move
	defer jogo.Mutex.Unlock()

	dy := 1
	if !alien.Subindo {
		dy = -1
	}
	nx := alien.X
	ny := alien.Y + dy // Mudando a posição no eixo Y, para o movimento vertical

	// Verifica se a nova posição 'ny' está dentro dos limites do mapa
	if ny < 0 || ny >= len(jogo.Mapa) {
		alien.Subindo = !alien.Subindo // Inverte a direção se atingir o limite do mapa
		return
	}

	// Verifica se a posição é válida antes de acessar o mapa
	if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
		destino := jogo.Mapa[ny][nx]
		if destino.tangivel || (jogo.PosX == nx && jogo.PosY == ny) {
			alien.Subindo = !alien.Subindo // Inverte a direção se o alien bater em algo
			return
		}

		// Move o alien
		jogo.Mapa[alien.Y][alien.X] = Vazio    // Limpa a posição anterior do alien
		jogo.Mapa[ny][nx] = Alien              // Coloca o alien na nova posição
		alien.X = nx                            // Atualiza a posição X
		alien.Y = ny                            // Atualiza a posição Y
	}
}










