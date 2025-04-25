// main_test.go
package main

import (
	"strings"
	"sync"
	"testing"
)

func TestBuildAlphabet(t *testing.T) {
	expected := " ABCDEFGHIJKLMNOPQRSTUVWXYZ-"
	result := buildAlphabet()
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestPercentProbability(t *testing.T) {
	for range 1000 {
		result := percentProbability()
		if result < 0 || result > 100 {
			t.Errorf("percentProbability returned out-of-range value: %d", result)
		}
	}
}

func TestRandomChar(t *testing.T) {
	alphabet := buildAlphabet()
	for range 1000 {
		char := randomChar()
		if !strings.ContainsRune(alphabet, char) {
			t.Errorf("randomChar returned invalid character: %q", char)
		}
	}
}

func TestCalcFitness(t *testing.T) {
	tests := []struct {
		guess   string
		answer  string
		fitness int
	}{
		{"HELLO", "HELLO", 5},
		{"HELLO", "WORLD", 1},
		{"HELLO", "HEART", 2},
	}

	for _, test := range tests {
		result := calcFitness(test.guess, test.answer)
		if result != test.fitness {
			t.Errorf("calcFitness(%q, %q) = %d; want %d", test.guess, test.answer, result, test.fitness)
		}
	}
}

func TestMutate(t *testing.T) {
	alphabet := buildAlphabet()
	guess := "HELLO"
	answer := "WORLD"
	ch := make(chan offspring, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go mutate(1, guess, answer, ch, &wg)
	wg.Wait()
	close(ch)

	mutated := <-ch
	if len(mutated.guess) != len(guess) {
		t.Errorf("Mutated guess length mismatch: got %d, want %d", len(mutated.guess), len(guess))
	}
	for _, char := range mutated.guess {
		if !strings.ContainsRune(alphabet, char) {
			t.Errorf("Mutated guess contains invalid character: %q", char)
		}
	}
}

func TestCrossover(t *testing.T) {
	parent1 := offspring{guess: "HELLO", fitness: 3}
	parent2 := offspring{guess: "WORLD", fitness: 2}
	answer := "HELLO"

	child := crossover(parent1, parent2, answer)
	if len(child.guess) != len(parent1.guess) {
		t.Errorf("Crossover child length mismatch: got %d, want %d", len(child.guess), len(parent1.guess))
	}
	if child.fitness != calcFitness(child.guess, answer) {
		t.Errorf("Crossover child fitness mismatch: got %d, want %d", child.fitness, calcFitness(child.guess, answer))
	}
}
