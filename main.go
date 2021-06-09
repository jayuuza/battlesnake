package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type Game struct {
	ID      string `json:"id"`
	Timeout int32  `json:"timeout"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Battlesnake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int32   `json:"health"`
	Body   []Coord `json:"body"`
	Head   Coord   `json:"head"`
	Length int32   `json:"length"`
	Shout  string  `json:"shout"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`
}

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}

type GameRequest struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

// HandleIndex is called when your Battlesnake is created and refreshed
// by play.battlesnake.com. BattlesnakeInfoResponse contains information about
// your Battlesnake, including what it should look like on the game board.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "jayuuza",
		Color:      "#ff6600",
		Head:       "pixel",
		Tail:       "pixel",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleStart is called at the start of each game your Battlesnake is playing.
// The GameRequest object contains information about the game that's about to start.
// TODO: Use this function to decide how your Battlesnake is going to look on the board.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Printf("START GAME %+v\n", request)
}

// HandleMove is called for each turn of each game.
// Valid responses are "up", "down", "left", or "right".
func HandleMove(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	move := makeMove(request)

	fmt.Printf("MOVE: %s\n", move.Move)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(move)
	if err != nil {
		log.Fatal(err)
	}
}

// validMoves returns moves that won't result in death for a given position
func validMoves(pos Coord, board Board) []string {
	var moves []string

	left := Coord{
		X: pos.X - 1,
		Y: pos.Y,
	}
	if isValid(left, board) {
		moves = append(moves, "left")
	}

 	right := Coord{
		X: pos.X + 1,
		Y: pos.Y,
	}
	if isValid(right, board) {
		moves = append(moves, "right")
	}

	down := Coord{
		X: pos.X,
		Y: pos.Y + 1,
	}
	if isValid(down, board) {
		moves = append(moves, "down")
	}

	up := Coord{
		X: pos.X,
		Y: pos.Y - 1,
	}
	if isValid(up, board) {
		moves = append(moves, "up")
	}

	return moves
}

func isValid(pos Coord, board Board) bool {
	return !isEdge(pos, board) && !isSnake(pos, board)
}

func isEdge(pos Coord, board Board) bool {
	return pos.X > board.Width - 1 ||  pos.Y > board.Height + 1 || pos.X < 0 || pos.Y < 0
}

func isFood(pos Coord, board Board) bool {
	for _, coord := range board.Food {
		if coord.Y == pos.Y && coord.X == pos.X {
			return true
		}
	}
	return false
}

func isSnake(pos Coord, board Board) bool {
	for _, snake := range board.Snakes {
		for _, coord := range snake.Body {
			if coord.Y == pos.Y && coord.X == pos.X {
				return true
			}
		}
		if snake.Head.Y == pos.Y && snake.Head.X == pos.X {
			return true
		}
	}
	return false
}

func makeMove(game GameRequest) MoveResponse {
	possibleMoves := validMoves(game.You.Head, game.Board)
	if len(possibleMoves) == 0 {
		return randomMove()
	}

	return MoveResponse{
		Move: possibleMoves[rand.Intn(len(possibleMoves))],
	}
}

// randomMove Chooses a random direction to move in
func randomMove() MoveResponse {
	possibleMoves := []string{"up", "down", "left", "right"}
	move := possibleMoves[rand.Intn(len(possibleMoves))]

	return MoveResponse{
		Move: move,
	}
}

// HandleEnd is called when a game your Battlesnake was playing has ended.
// It's purely for informational purposes, no response required.
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("END\n")
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)

	fmt.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}