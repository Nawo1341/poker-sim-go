package main

import (
	"sort"
)

// 役の強さを定数で定義
const (
	HighCard = iota // 0
	OnePair         // 1
	TwoPair         // 2
	ThreeOfAKind    // 3
	Straight        // 4
	Flush           // 5
	FullHouse       // 6
	FourOfAKind     // 7
	StraightFlush   // 8
)

var handNames = []string{
	"ハイカード", "ワンペア", "ツーペア", "スリーカード",
	"ストレート", "フラッシュ", "フルハウス", "フォーカード", "ストレートフラッシュ",
}

// EvaluateHand : 7枚のカードから最強の役とスコアを計算する
// 戻り値: (役の名前, 役の強さ数値, 勝敗判定用の詳細スコア)
func EvaluateHand(hole []Card, board []Card) (string, int, int) {
	// 全7枚をまとめる
	allCards := append([]Card{}, hole...)
	allCards = append(allCards, board...)

	// 数字の強い順に並べ替える (A=14, ..., 2=2)
	sort.Slice(allCards, func(i, j int) bool {
		return allCards[i].Val > allCards[j].Val
	})

	// 1. フラッシュのチェック (同じマークが5枚以上)
	isFlush, flushSuit := checkFlush(allCards)
	
	// 2. ストレートのチェック (連続した数字が5枚以上)
	// (フラッシュの場合は、そのマークだけでストレート判定が必要)
	isStraight, highRank := checkStraight(allCards, isFlush, flushSuit)

	// --- 判定ロジック ---

	// ストレートフラッシュ
	if isFlush && isStraight {
		return handNames[StraightFlush], StraightFlush, highRank
	}
	// フォーカード、フルハウス、スリーカード、ペア系
	// (同じ数字の重なりをチェック)
	rank, score := checkPairs(allCards)
	
	// 役の強さ比較
	// すでにペア系判定で役が見つかっている場合、フラッシュやストレートと比較
	// (ペア系のロジックでFullHouseなどが出ているならそれが優先)
	if rank >= FourOfAKind {
		return handNames[rank], rank, score
	} else if rank == FullHouse {
		return handNames[rank], rank, score
	} 

	// フラッシュ
	if isFlush {
		// フラッシュの強さは一番強いカードで決まる
		return handNames[Flush], Flush, getHighCardScore(allCards, flushSuit)
	}

	// ストレート
	if isStraight {
		return handNames[Straight], Straight, highRank
	}

	// 残りのペア系 (3カード、2ペア、1ペア) またはハイカード
	return handNames[rank], rank, score
}

// checkFlush : フラッシュがあるか判定
func checkFlush(cards []Card) (bool, string) {
	suitsCount := map[string]int{"♠": 0, "♥": 0, "♦": 0, "♣": 0}
	for _, c := range cards {
		suitsCount[c.Suit]++
	}
	for s, count := range suitsCount {
		if count >= 5 {
			return true, s
		}
	}
	return false, ""
}

// checkStraight : ストレートがあるか判定
func checkStraight(cards []Card, isFlush bool, flushSuit string) (bool, int) {
	// 重複を除いてユニークな数字リストを作る
	// (フラッシュの場合はそのマークのカードだけで判定)
	uniqueVals := []int{}
	seen := map[int]bool{}

	for _, c := range cards {
		if isFlush && c.Suit != flushSuit {
			continue // フラッシュなら違うマークは無視
		}
		if !seen[c.Val] {
			uniqueVals = append(uniqueVals, c.Val)
			seen[c.Val] = true
		}
	}
	// A(14)がある場合、A-2-3-4-5のストレートのために1としても扱う
	if seen[14] {
		uniqueVals = append(uniqueVals, 1)
	}
	
	sort.Sort(sort.Reverse(sort.IntSlice(uniqueVals)))

	// 連続チェック
	count := 0
	for i := 0; i < len(uniqueVals)-1; i++ {
		if uniqueVals[i]-1 == uniqueVals[i+1] {
			count++
			if count >= 4 { // 4回連続成功＝5枚連続
				return true, uniqueVals[i-3] // 一番強い数字を返す
			}
		} else {
			count = 0
		}
	}
	return false, 0
}

// checkPairs : ペア、3カード、4カード、フルハウスを判定
func checkPairs(cards []Card) (int, int) {
	// 数字ごとの枚数をカウント
	counts := map[int]int{}
	for _, c := range cards {
		counts[c.Val]++
	}

	four, three, pair := 0, 0, 0
	// 強い数字から順にチェック
	for i := 14; i >= 2; i-- {
		if counts[i] == 4 { four = i }
		if counts[i] == 3 { 
			if three == 0 { three = i } else if pair == 0 { pair = i } // すでに3カードがあればそれはペア扱い（フルハウス用）
		}
		if counts[i] == 2 {
			if pair == 0 { pair = i } else if three != 0 { /* フルハウス成立待ち */ }
		}
	}

	// 判定順序
	if four > 0 {
		return FourOfAKind, four*100 // 簡易スコア
	}
	// フルハウス (3枚組 + 2枚組以上)
	// countsをもう一度見て判定（簡易実装）
	has3 := 0
	has2 := 0
	for i := 14; i >= 2; i-- {
		if counts[i] >= 3 && has3 == 0 { has3 = i; continue }
		if counts[i] >= 2 && has2 == 0 { has2 = i }
	}
	if has3 > 0 && has2 > 0 {
		return FullHouse, has3*100 + has2
	}

	if has3 > 0 {
		return ThreeOfAKind, has3
	}

	// ツーペア・ワンペア判定
	pairs := []int{}
	for i := 14; i >= 2; i-- {
		if counts[i] == 2 {
			pairs = append(pairs, i)
		}
	}
	if len(pairs) >= 2 {
		return TwoPair, pairs[0]*100 + pairs[1]
	}
	if len(pairs) == 1 {
		return OnePair, pairs[0]
	}

	// ハイカード (一番強いカード)
	return HighCard, cards[0].Val
}

func getHighCardScore(cards []Card, suit string) int {
	for _, c := range cards {
		if c.Suit == suit {
			return c.Val
		}
	}
	return 0
}