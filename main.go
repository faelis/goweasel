package main

import (
	"cmp"
	"fmt"
	"math/big"
	"os"
	"slices"
	"sync"
	"time"

	"crypto/rand"
)

func buildAlphabet() (asciiChars string) { // build alphabet
	// asciiChars += string(rune(10)) // CR
	// asciiChars += string(rune(13)) // LF
	asciiChars += string(rune(32)) // space
	// for i := 33; i < 48; i++ { // special chars
	// 	asciiChars += string(rune(i))
	// }
	// for i := 48; i < 58; i++ { // digits chars
	// 	asciiChars += string(rune(i))
	// }
	// for i := 58; i < 65; i++ { // special chars
	// 	asciiChars += string(rune(i))
	// }
	for i := 65; i < 91; i++ { // upper letter chars
		asciiChars += string(rune(i))
	}
	// for i := 91; i < 97; i++ { // special chars
	// 	asciiChars += string(rune(i))
	// }
	// for i := 97; i < 123; i++ { // lower letter chars
	// 	asciiChars += string(rune(i))
	// }
	// for i := 123; i < 127; i++ { // special chars
	// 	asciiChars += string(rune(i))
	// }
	asciiChars += "-" // custom chars
	return
}

var alphabet = buildAlphabet()

const defaultGuess = "METHINKS IT IS LIKE A WEASEL" // default string to guess if none is provided
const mutationRate = 5                              // mutation rate in percent used to mutate each character of the guess
const offspringMaxCount = 100                       // number of concurent offspring to generate

type organism struct { // parent organism struct
	guess      string
	answer     string
	generation int
	fitness    int
}

type offspring struct { // offspring struct
	id      int
	guess   string
	fitness int
}

func percentProbability() int { // return a random number between 0 and 100
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(100)))
	return int(n.Int64())
}

func randomChar() rune { // return a random character from the alphabet
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
	return rune(alphabet[int(n.Int64())])
}

func (o *organism) init(answer string) { // initialize the organism with a random guess based on the answer
	for i := 0; i < len(answer); i++ {
		o.guess = string(append([]rune(o.guess), randomChar())) // append a random character to the guess casted as a rune sclice, then cast to a string
	}
	o.answer = answer // set the answer in the parent organism
}

func calcFitness(guess, answer string) int { // calculate the fitness of the guess based on the answer
	var fitness int
	for i := range []rune(guess) {
		if []rune(guess)[i] == []rune(answer)[i] { // increase fitness when 2 identicals characters are at the same position between the guess and the answer
			fitness++
		}
	}
	return fitness
}

func mutate(id int, guess string, answer string, ch chan offspring, wg *sync.WaitGroup) {
	defer wg.Done()
	mutatedGuess := []rune(guess)
	for i := range []rune(guess) {
		if percentProbability() < mutationRate { // mutate the character if the random number is below the mutation rate
			mutatedGuess[i] = randomChar()
		}
	}
	newOffspring := offspring{id: id, guess: string(mutatedGuess), fitness: calcFitness(string(mutatedGuess), answer)} // initiate the offspring struct with the mutated guess and its fitness
	ch <- newOffspring                                                                                                 // send the offspring to the channel
}

func crossover(i, j offspring, answer string) offspring { // crossover 2 offsprings
	var k, l offspring
	k.guess = i.guess[:len(i.guess)/2] + j.guess[len(j.guess)/2:] // crossover the first half of the first offspring with the second half of the second offspring
	l.guess = j.guess[:len(j.guess)/2] + i.guess[len(i.guess)/2:] // crossover the first half of the second offspring with the second half of the first offspring
	k.fitness = calcFitness(k.guess, answer)                      // calculate the fitness of the first offspring
	l.fitness = calcFitness(l.guess, answer)                      // calculate the fitness of the second offspring
	if k.fitness > l.fitness {                                    // return the offspring with the best fitness
		return k
	} else {
		return l
	}
}

func (o *organism) evolve() { // evolve the organism
	var wg sync.WaitGroup                         // wait group to wait for all the offspring to be generated
	ch := make(chan offspring, offspringMaxCount) // channel to send the offspring to the main thread

	for i := 1; i <= offspringMaxCount; i++ { // generate the offspring with a goroutine for each
		wg.Add(1)
		go mutate(i, o.guess, o.answer, ch, &wg) //	launch the goroutine
	}
	go func(ch chan offspring, wg *sync.WaitGroup) { // wait for all the offspring to be generated
		wg.Wait()
		close(ch) // close the channel
	}(ch, &wg)

	offsprings := []offspring{}
	for i := range ch { // iterate over the offspring channel
		offsprings = append(offsprings, i) // append the offspring to the offsprings slice
	}

	slices.SortFunc(offsprings, func(i, j offspring) int { // sort the offsprings by fitness
		return cmp.Compare(i.fitness, j.fitness)
	})

	offsprings = slices.Compact(offsprings) // then remove the duplicates

	bestOffspring := crossover(offsprings[len(offsprings)-1], offsprings[len(offsprings)-2], o.answer) // crossover the 2 best offsprings

	// keep the best offspring
	o.guess = bestOffspring.guess
	o.fitness = bestOffspring.fitness

	o.generation++ // increment the generation counter
}

func main() {
	var o organism
	if len(os.Args) < 2 { // if no argument is provided, use the default guess
		o.init(defaultGuess)
	} else if len(os.Args) > 2 { // if too many arguments are provided, panic
		panic("Too many arguments")
	} else { // if one argument is provided, use it as the answer
		o.init(os.Args[1])
	}
	for {
		fmt.Print("\033[H\033[2J") // refresh the screen
		o.evolve()                 // evolve the organism
		fmt.Println(o.guess)
		time.Sleep(10 * time.Millisecond) // sleep because go is too fast :D
		// bufio.NewReader(os.Stdin).ReadBytes('\n') // wait for the user to press enter
		if o.guess == o.answer { // stop when the guessing game is ended
			fmt.Printf("Generation %d\n", o.generation)
			break
		}
	}
}
