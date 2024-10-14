package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth   = 320
	screenHeight  = 240
	gridSize      = 10
	snakeSpeed    = 5  // Controla a velocidade da cobra (maior valor = mais devagar)
	placarHeight  = 20 // Altura reservada para o placar no topo
	borderPadding = 10 // Largura da borda ao redor do campo de jogo
)

type Game struct {
	snake      []image.Point
	direction  image.Point
	food       image.Point
	score      int
	gameOver   bool
	frameCount int // Contador de frames para controlar a velocidade
}

func (g *Game) init() {
	g.snake = []image.Point{{X: screenWidth / 2, Y: screenHeight / 2}}
	g.direction = image.Point{X: gridSize, Y: 0}
	g.spawnFood()
	g.score = 0
	g.gameOver = false
	g.frameCount = 0
}

func (g *Game) spawnFood() {
	rand.Seed(time.Now().UnixNano())
	g.food = image.Point{
		X: rand.Intn((screenWidth-borderPadding*2)/gridSize)*gridSize + borderPadding,
		Y: rand.Intn((screenHeight-placarHeight-borderPadding*2)/gridSize)*gridSize + placarHeight + borderPadding,
	}
}

func (g *Game) Update() error {
	// Verifica se o jogador pressionou 'Q' ou 'ESC' para sair
	if ebiten.IsKeyPressed(ebiten.KeyQ) || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	// Verifica se o jogador pressionou 'R' para reiniciar
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.init()
		return nil
	}

	if g.gameOver {
		return nil
	}

	g.frameCount++
	if g.frameCount < snakeSpeed {
		// Aguarda até atingir a contagem de frames para mover a cobra
		return nil
	}
	g.frameCount = 0 // Reseta o contador após o movimento

	// Controle da direção com o teclado
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.direction.Y == 0 {
		g.direction = image.Point{X: 0, Y: -gridSize}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.direction.Y == 0 {
		g.direction = image.Point{X: 0, Y: gridSize}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.direction.X == 0 {
		g.direction = image.Point{X: -gridSize, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.direction.X == 0 {
		g.direction = image.Point{X: gridSize, Y: 0}
	}

	// Movimento da cobra
	newHead := image.Point{
		X: g.snake[0].X + g.direction.X,
		Y: g.snake[0].Y + g.direction.Y,
	}

	// Verifica colisão com as bordas da área de jogo
	if newHead.X < borderPadding || newHead.X >= screenWidth-borderPadding ||
		newHead.Y < placarHeight+borderPadding || newHead.Y >= screenHeight-borderPadding {
		g.gameOver = true
		return nil
	}

	// Verifica colisão com a própria cobra
	for _, segment := range g.snake {
		if newHead == segment {
			g.gameOver = true
			return nil
		}
	}

	// Adiciona a nova cabeça e remove a cauda
	g.snake = append([]image.Point{newHead}, g.snake...)

	// Verifica se comeu a comida
	if newHead == g.food {
		g.score++
		g.spawnFood()
	} else {
		// Remove a cauda se não comeu
		g.snake = g.snake[:len(g.snake)-1]
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Preenche o fundo de preto
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Desenha a borda
	borderColor := color.RGBA{255, 255, 255, 255}
	ebitenutil.DrawRect(screen, 0, placarHeight, screenWidth, borderPadding, borderColor)                                       // Topo
	ebitenutil.DrawRect(screen, 0, screenHeight-borderPadding, screenWidth, borderPadding, borderColor)                         // Base
	ebitenutil.DrawRect(screen, 0, placarHeight, borderPadding, screenHeight-placarHeight, borderColor)                         // Esquerda
	ebitenutil.DrawRect(screen, screenWidth-borderPadding, placarHeight, borderPadding, screenHeight-placarHeight, borderColor) // Direita

	// Desenha a cobra
	for _, segment := range g.snake {
		ebitenutil.DrawRect(screen, float64(segment.X), float64(segment.Y), gridSize, gridSize, color.RGBA{0, 255, 0, 255})
	}

	// Desenha a comida
	ebitenutil.DrawRect(screen, float64(g.food.X), float64(g.food.Y), gridSize, gridSize, color.RGBA{255, 0, 0, 255})

	// Desenha o placar no topo à direita
	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, screenWidth-80, 5)

	// Mensagem de game over
	if g.gameOver {
		ebitenutil.DebugPrint(screen, "Game Over! Press R to Restart")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")

	game := &Game{}
	game.init()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
