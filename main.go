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

	// カードを生成
	for _, s := range suits {
		for _, r := range ranks {
			deck = append(deck, Card{Suit: s, Rank: r})
		}
	}

	// 現在時刻を種にしてシャッフル (Goの乱数は固定されがちなので)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}

func main() {
	deck := createDeck()
	fmt.Printf("デッキ作成完了: %d枚\n", len(deck))
	
	// 最初の1枚を確認
	topCard := deck[0]
	fmt.Printf("一番上のカード: %s%s\n", topCard.Suit, topCard.Rank)
}