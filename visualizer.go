package main

import (
	"fmt"
)

// カードの見た目（幅7文字 x 高さ5行）
// ┌─────┐
// │A    │
// │  ♠  │
// │    A│
// └─────┘

func getCardAA(c Card) []string {
	rankLeft := fmt.Sprintf("%-2s", c.Rank)  // 左上の数字 (例: "10", "A ")
	rankRight := fmt.Sprintf("%2s", c.Rank) // 右下の数字

	// 色の設定
	colorCode := ""
	resetCode := ""
	if c.Suit == "♥" || c.Suit == "♦" {
		colorCode = "\033[31m" // 赤
		resetCode = "\033[0m"  // リセット
	}

	// 各行を作成
	lines := []string{
		"┌─────┐",
		fmt.Sprintf("│%s   │", rankLeft),
		fmt.Sprintf("│  %s%s%s  │", colorCode, c.Suit, resetCode), // 真ん中のマークだけ色をつける
		fmt.Sprintf("│   %s│", rankRight),
		"└─────┘",
	}
	return lines
}

// 裏向きのカード
func getHiddenCardAA() []string {
	return []string{
		"┌─────┐",
		"│/////│",
		"│/////│",
		"│/////│",
		"└─────┘",
	}
}

// 複数のカードを横並びで表示する関数
func printCardsAA(cards []Card) {
	if len(cards) == 0 {
		return
	}

	// 描画用のバッファ（5行分）
	outputLines := make([]string, 5)

	for _, c := range cards {
		aa := getCardAA(c)
		for i := 0; i < 5; i++ {
			outputLines[i] += aa[i] + " " // カードの間にスペースを入れる
		}
	}

	// 出力
	for _, line := range outputLines {
		fmt.Println(line)
	}
}

// 裏向きカードをn枚表示する関数
func printHiddenHandAA(n int) {
	outputLines := make([]string, 5)
	for k := 0; k < n; k++ {
		aa := getHiddenCardAA()
		for i := 0; i < 5; i++ {
			outputLines[i] += aa[i] + " "
		}
	}
	for _, line := range outputLines {
		fmt.Println(line)
	}
}