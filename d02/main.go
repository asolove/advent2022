package main

import (
  "bufio"
	"github.com/pkg/errors"
  "fmt"
  "log"
  "os"
)

type Move int;

const (
  Rock Move = iota;
  Paper =     iota;
  Scissors =  iota;
)

type Outcome int;
const (
  LOSE Outcome = -1;
  DRAW = 0;
  WIN = 1;
)

const (
  WIN_SCORE = 6;
  DRAW_SCORE = 3;
  LOSE_SCORE = 0;
)

func scoreGame(opponentMove, myMove Move) int {
  return shapeScore(myMove) + outcomeScore(opponentMove, myMove);
}

func shapeScore(move Move) int {
  return int(move) + 1;
}

func outcome(opponentMove, myMove Move) Outcome {
  switch true {
  case opponentMove == myMove:
    return DRAW;
  case opponentMove == Rock && myMove == Scissors:
    return LOSE;
  case opponentMove == Rock && myMove == Paper:
    return WIN;
  case opponentMove == Scissors && myMove == Paper:
    return LOSE;
  }
  return -outcome(myMove, opponentMove);
}

func outcomeScore(opponentMove, myMove Move) int {
  // Replace with direct math?
  switch outcome(opponentMove, myMove) {
  case WIN:
    return WIN_SCORE;
  case DRAW:
    return DRAW_SCORE;
  case LOSE:
    return LOSE_SCORE;
  }
  return 0;
}

func readOpponentMove(r rune) (Move, error) {
  switch r {
  case 'A':
    return Rock, nil;
  case 'B':
    return Paper, nil;
  case 'C':
    return Scissors, nil;
  }
  return Rock, errors.Errorf("Invalid opponent move: %r", r);
}


func readMyMove(r rune) (Move, error) {
  switch r {
  case 'X':
    return Rock, nil;
  case 'Y':
    return Paper, nil;
  case 'Z':
    return Scissors, nil;
  }
  return Rock, errors.Errorf("Invalid self move: %r", r);
}

func readOutcome(r rune) (Outcome, error) {
  switch r {
  case 'X':
    return LOSE, nil;
  case 'Y':
    return DRAW, nil;
  case 'Z':
    return WIN, nil;
  }
  return LOSE, errors.Errorf("Invalid self move: %r", r);
}

func findMyMove(opponentMove Move, desiredOutcome Outcome) (Move, error) {
  moves := []Move{Rock, Paper, Scissors};
  for _, move := range moves {
    if outcome(opponentMove, move) == desiredOutcome {
      return move, nil;
    }
  }
  return Rock, errors.Errorf("Cannot findMyMove");
}

func main() {
  totalScore := 0;
	s := bufio.NewScanner(os.Stdin);
	for s.Scan() {
    line := s.Text();
    var opponentMove, myMove Move;
    var desiredOutcome Outcome;
    var err error;
    if opponentMove, err = readOpponentMove(rune(line[0])); err != nil {
      log.Fatal(err);
    }
    if desiredOutcome, err = readOutcome(rune(line[2])); err != nil {
      log.Fatal(err);
    }
    if myMove, err = findMyMove(opponentMove, desiredOutcome); err != nil {
      log.Fatal(err);
    }
    totalScore += scoreGame(opponentMove, myMove);
  }
  fmt.Printf("Total: %d\n", totalScore);
}
