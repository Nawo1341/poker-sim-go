package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Card : トランプ1枚を表す構造体
type Card struct {
	Suit string // マーク (♠, ♥, ♦, ♣)
	Rank string // 数字 (2, 3, ... A)
}

// 定数定義
var suits = []string{"♠", "♥", "♦", "♣"}
var ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

// createDeck : 52枚のデッキを作成してシャッフルする
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

// deal : デッキから指定枚数を引いて返す (引いた分、デッキは減る)
func deal(deck *[]Card, n int) []Card {
	// 配るカードを取り出す
	hand := (*deck)[:n]
	// デッキを更新（配った分を削除する）
	*deck = (*deck)[n:]
	return hand
}

func main() {
	// ゲームのセットアップ
	deck := createDeck()
	fmt.Printf("デッキ初期枚数: %d枚\n", len(deck))

	// プレイヤーとCPUに2枚ずつ配る
	// (&deck と書くことで、deckそのものを書き換えられるようにしています)
	playerHand := deal(&deck, 2)
	cpuHand := deal(&deck, 2)

	fmt.Println("--- テキサスホールデム (Go版) ---")
	fmt.Printf("あなたの手札: %v\n", playerHand)
	fmt.Printf("CPUの手札:   %v\n", cpuHand)
	fmt.Printf("残りの山札: %d枚\n", len(deck))
}