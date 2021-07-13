/*
Package mchain implements Markov text generation. It maps prefixes to suffixes found in texts files.
It string together prefixes and suffixes randomly to generate logical sentences.
*/

package mchain

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

type MarkovChain struct {
	File      []byte
	Words     []string
	WordsSize int
	Chain     map[[2]string][]string
}

type triplet struct{ w1, w2, w3 string }

func NewMarkovChain(txtFile []byte) (m *MarkovChain) {
	m = &MarkovChain{File: txtFile}
	m.readFile()
	m.train()
	return
}

func (m *MarkovChain) readFile() {
	text := strings.TrimSpace(string(m.File))
	text = strings.Replace(text, "\n", " ", -1)
	words := strings.Fields(text)
	m.Words = words
	m.WordsSize = len(words)
}

func (m *MarkovChain) train() {
	m.Chain = make(map[[2]string][]string)
	for t := range m.triples() {
		key := [2]string{t.w1, t.w2}
		if values, ok := m.Chain[key]; ok {
			m.Chain[key] = append(values, t.w3)
		} else {
			m.Chain[key] = []string{t.w3}
		}
	}
}

func (m *MarkovChain) triples() <-chan triplet {
	if m.WordsSize < 3 {
		err := errors.New("source text too short")
		panic(err)
	}
	ch := make(chan triplet)
	go func() {
		for i := 0; i < m.WordsSize-2; i++ {
			t := triplet{
				w1: m.Words[i],
				w2: m.Words[i+1],
				w3: m.Words[i+2],
			}
			ch <- t
		}
		close(ch)
	}()
	return ch
}

func (m *MarkovChain) ShowChain() {
	for k, v := range m.Chain {
		fmt.Println(k, "=>", v)
	}
}

func (m *MarkovChain) Generate(size int) string {
	seed := rand.Intn(m.WordsSize - 3)
	seedWord, nextWord := m.Words[seed], m.Words[seed+1]
	w1, w2 := seedWord, nextWord
	current := []string{}
	sentences := []string{}
	counter := 0
	for counter < size {
		current = append(current, w1)
		if strings.HasSuffix(w1, ".") {
			sentence := strings.Join(current, " ")
			sentences = append(sentences, sentence)
			current = []string{}
		}
		v := m.Chain[[2]string{w1, w2}]
		if len(v) == 0 { 
			newSeed := rand.Intn(m.WordsSize - 3)
			newWord := m.Words[newSeed]
			w1, w2 = w2, newWord
		} else {
			next := rand.Intn(len(v))
			w1, w2 = w2, v[next]
		}
		counter ++
	}
	current = append(current, w2)
	sentence := strings.Join(current, " ")
	sentence += "."
	sentences = append(sentences, sentence)
	return strings.Join(sentences, " ")
}
