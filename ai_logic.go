package main

import (
	"math/rand"
	"time"
)

// DecideCpuAction : CPUのアクションを決定するメイン関数
func DecideCpuAction(cpuHand []Card, board []Card, pot int, toCall int) (string, int) {
	// 1. 勝率を計算 (モンテカルロ・シミュレーション)
	winRate := calculateWinRate(cpuHand, board, 1000) // 1000回シミュレーション

	// 2. ポットオッズの計算
	// (コールした場合、最終的なポットに対して自分が払う割合)
	// 例: Pot=100, Call=50 -> TotalPot=150. 必要勝率は 50/150 = 33%
	var requiredWinRate float64
	if toCall > 0 {
		requiredWinRate = float64(toCall) / float64(pot+toCall)
	} else {
		requiredWinRate = 0 // チェックならタダなので必要勝率は0
	}

	// 3. アクションの決定（期待値に基づく）
	
	// 圧倒的に強い場合 (勝率80%以上) -> レイズしてバリューを取りに行く
	if winRate > 0.8 {
		return "RAISE", 50
	}

	// 勝率が必要勝率を上回っている -> コール（またはチェック）
	if winRate >= requiredWinRate {
		// かなり強い (勝率60%以上) ならレイズ頻度を上げる
		if winRate > 0.6 && rand.Float64() < 0.5 {
			return "RAISE", 50
		}
		if toCall == 0 {
			return "CHECK", 0
		}
		return "CALL", toCall
	}

	// 勝率が低い -> 基本フォールドだが、たまにブラフ
	// (チェックできるならチェック)
	if toCall == 0 {
		return "CHECK", 0
	}
	// 10%の確率でブラフのレイズ
	if rand.Float64() < 0.1 {
		return "RAISE", 50
	}

	return "FOLD", 0
}

// calculateWinRate : シミュレーションを行って勝率(0.0~1.0)を返す
func calculateWinRate(myHand []Card, currentBoard []Card, trials int) float64 {
	wins := 0
	
	// 乱数の種をセット
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < trials; i++ {
		// 1. デッキを再現 (既に出ているカードを除外)
		deck := createSimDeck(myHand, currentBoard)

		// 2. 相手の手札をランダムに配る (2枚)
		oppHand := deal(&deck, 2)

		// 3. ボードの残り枚数をランダムに配る (最大5枚になるまで)
		simBoard := make([]Card, len(currentBoard))
		copy(simBoard, currentBoard)
		
		needed := 5 - len(simBoard)
		if needed > 0 {
			simBoard = append(simBoard, deal(&deck, needed)...)
		}

		// 4. 勝敗判定
		_, myRank, myScore := EvaluateHand(myHand, simBoard)
		_, oppRank, oppScore := EvaluateHand(oppHand, simBoard)

		if myRank > oppRank {
			wins++
		} else if myRank == oppRank && myScore > oppScore {
			wins++
		} else if myRank == oppRank && myScore == oppScore {
			wins++ // 引き分けも勝ちカウント(簡易) または0.5加算が正確だが今回は簡易化
		}
	}

	return float64(wins) / float64(trials)
}

// createSimDeck : シミュレーション用に、すでに使われたカードを除いたデッキを作る
func createSimDeck(exclude1 []Card, exclude2 []Card) []Card {
	fullDeck := createDeck() // 通常のデッキ生成関数を利用
	simDeck := []Card{}

	// 使われているカードをマップに記録
	used := map[string]bool{}
	for _, c := range exclude1 { used[c.Suit+c.Rank] = true }
	for _, c := range exclude2 { used[c.Suit+c.Rank] = true }

	for _, c := range fullDeck {
		if !used[c.Suit+c.Rank] {
			simDeck = append(simDeck, c)
		}
	}
	// シャッフル
	rand.Shuffle(len(simDeck), func(i, j int) {
		simDeck[i], simDeck[j] = simDeck[j], simDeck[i]
	})
	
	return simDeck
}