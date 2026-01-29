package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Card : トランプ1枚を表す構造体
type Card struct {
	Suit string
	Rank string
}

var suits = []string{"♠", "♥", "♦", "♣"}
var ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

func createDeck() []Card {
	var deck []Card
	for _, s := range suits {
		for _, r := range ranks {
			deck = append(deck, Card{Suit: s, Rank: r})
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

func deal(deck *[]Card, n int) []Card {
	hand := (*deck)[:n]
	*deck = (*deck)[n:]
	return hand
}

func main() {
	deck := createDeck()
	fmt.Printf("デッキ初期枚数: %d枚\n", len(deck))

	// 1. プリフロップ：プレイヤーとCPUに2枚ずつ配る
	playerHand := deal(&deck, 2)
	cpuHand := deal(&deck, 2)

	// 2. フロップ：場に3枚出す
	flop := deal(&deck, 3)

	fmt.Println("--- テキサスホールデム (Go版) ---")
	fmt.Printf("あなたの手札: %v\n", playerHand)
	fmt.Printf("CPUの手札:   %v\n", cpuHand)
	
	fmt.Println("\n--- 場のカード (Flop) ---")
	fmt.Printf("[%v] [%v] [%v]\n", flop[0], flop[1], flop[2])

	fmt.Printf("\n残りの山札: %d枚\n", len(deck))
}