package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings" // ← これを復活させてください！
	"time"
)

// --- 定数・構造体 ---

type Card struct {
	Suit string
	Rank string
	Val  int // 強さ比較用 (2=2, ... A=14)
}

type Player struct {
	Name   string
	Hand   []Card
	Chips  int
	IsCPU  bool
	Folded bool
}

var suits = []string{"♠", "♥", "♦", "♣"}
var ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

// --- カード操作系 ---

func (c Card) String() string {
	s := fmt.Sprintf("%s%s", c.Suit, c.Rank)
	if c.Suit == "♥" || c.Suit == "♦" {
		return fmt.Sprintf("\033[31m%s\033[0m", s) // 赤字
	}
	return s
}

func createDeck() []Card {
	deck := []Card{}
	for _, s := range suits {
		for i, r := range ranks {
			deck = append(deck, Card{Suit: s, Rank: r, Val: i + 2})
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

// --- ゲーム進行系 ---

func main() {
	fmt.Println("=== TEXAS HOLD'EM POKER SIMULATOR ===")

	// プレイヤー初期設定
	human := &Player{Name: "You", Chips: 500, IsCPU: false}
	cpu := &Player{Name: "CPU", Chips: 500, IsCPU: true}

	// ゲームループ（どちらかのチップが尽きるまで）
	round := 1
	for human.Chips > 0 && cpu.Chips > 0 {
		playRound(round, human, cpu)
		round++
		if human.Chips > 0 && cpu.Chips > 0 {
			fmt.Println("\n[Enter]キーで次のゲームへ...")
			bufio.NewScanner(os.Stdin).Scan()
		}
	}

	fmt.Println("\n=== GAME OVER ===")
	if human.Chips > 0 {
		fmt.Println("おめでとうございます！あなたの勝利です！")
	} else {
		fmt.Println("残念...破産しました。")
	}
}

func playRound(roundNum int, p1 *Player, p2 *Player) {
	deck := createDeck()
	pot := 0
	p1.Folded = false
	p2.Folded = false
	p1.Hand = deal(&deck, 2)
	p2.Hand = deal(&deck, 2)
	board := []Card{}

	fmt.Printf("\n--- ROUND %d ---\n", roundNum)
	fmt.Printf("あなたのチップ: $%d  |  CPUのチップ: $%d\n", p1.Chips, p2.Chips)

	// ブラインド（参加費）
	blind := 10
	p1.Chips -= blind
	p2.Chips -= blind
	pot += blind * 2
	fmt.Printf("アンティ(参加費)として $%d 支払いました。Pot: $%d\n", blind, pot)

	// 各フェーズの実行
	phases := []string{"Pre-Flop", "Flop", "Turn", "River"}

	for i, phase := range phases {
		// カードをめくる処理
		if phase == "Flop" {
			board = append(board, deal(&deck, 3)...)
		} else if phase == "Turn" || phase == "River" {
			board = append(board, deal(&deck, 1)...)
		}

		// 状況表示
		printTable(phase, p1, p2, board, pot)

		// ... (playRound関数の中の bettingRound 判定部分) ...

		if !bettingRound(p1, p2, &pot, phase, board) {
			fmt.Println("\n--- 手札公開 (Fold) ---")
			fmt.Println("YOU:")
			printCardsAA(p1.Hand)
			fmt.Println("CPU:")
			printCardsAA(p2.Hand)
			return
		}

		// どちらかがAll-inしてたら以降のベットはスキップなどの処理は今回は省略
		if i < 3 {
			fmt.Println("\n[Enter]キーで次へ...")
			bufio.NewScanner(os.Stdin).Scan()
		}
	}

	// ショーダウン
	showdown(p1, p2, board, pot)
}

// 状況を表示する
func printTable(phase string, p1, p2 *Player, board []Card, pot int) {
	fmt.Println("\n========================================")
	fmt.Printf("PHASE: %-10s |  POT: $%d\n", phase, pot)
	fmt.Println("========================================")

	// CPUの手札（裏向き表示）
	fmt.Printf("CPU  (Chips: $%d)\n", p2.Chips)
	printHiddenHandAA(2)

	fmt.Println("\n------------- BOARD -------------")
	// 場のカード
	if len(board) > 0 {
		printCardsAA(board)
	} else {
		fmt.Println("(まだカードはありません)")
	}
	fmt.Println("---------------------------------")

	// 自分の手札
	fmt.Printf("YOU  (Chips: $%d)\n", p1.Chips)
	printCardsAA(p1.Hand)

	fmt.Println("========================================")
}

// ベット処理（戻り値 false = 誰かが降りてラウンド終了）
// ★引数に board []Card を追加しました
func bettingRound(p1, p2 *Player, pot *int, phase string, board []Card) bool {
	// 簡易的なベットロジック（1ターンずつ行動）
	// 人間の行動
	action, amount := getPlayerAction(p1, *pot)
	if action == "FOLD" {
		fmt.Println("あなたは降りました。CPUの勝ちです。")
		p2.Chips += *pot
		return false
	}

	if amount > 0 {
		p1.Chips -= amount
		*pot += amount
		fmt.Printf("あなたは $%d ベットしました。\n", amount)
	} else {
		fmt.Println("あなたはチェックしました。")
	}

	// CPUの行動（AI）
	// ★ここでAIに board を渡しています
	cpuAction, cpuAmount := getCpuAction(p2, *pot, amount, phase, board)

	if cpuAction == "FOLD" {
		fmt.Println("CPUが降りました！あなたの勝ちです！")
		p1.Chips += *pot
		return false
	}

	if cpuAmount > 0 {
		p2.Chips -= cpuAmount
		*pot += cpuAmount
		fmt.Printf("CPUはコール(またはレイズ)しました ($%d)。\n", cpuAmount)
	} else {
		fmt.Println("CPUはチェックしました。")
	}

	return true
}

func getPlayerAction(p *Player, currentPot int) (string, int) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("アクションを選択 [c:Check/Call, r:Raise $50, f:Fold] > ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "c":
			return "CALL", 0 // 今回は簡単のためコール額計算は省略し0とする
		case "r":
			if p.Chips >= 50 {
				return "RAISE", 50
			} else {
				fmt.Println("チップが足りません！")
			}
		case "f":
			return "FOLD", 0
		default:
			fmt.Println("c, r, f のいずれかを入力してください")
		}
	}
}

// getCpuAction : ここで ai_logic.go の頭脳を使います！
// ★引数に board []Card を追加しました
func getCpuAction(p *Player, pot int, playerBet int, phase string, board []Card) (string, int) {
	fmt.Print("CPUが思考中...")

	// プレイヤーのベット額に合わせて「コールに必要な額」を計算
	toCall := playerBet

	// ★ここでAIロジック（DecideCpuAction）を呼び出すときに board を渡します
	action, amount := DecideCpuAction(p.Hand, board, pot, toCall)

	// ちょっとウェイトを入れて「考えている感」を出す
	time.Sleep(1 * time.Second)
	fmt.Println(" 決定！")

	return action, amount
}

func showdown(p1, p2 *Player, board []Card, pot int) {
	fmt.Println("\n\n################ SHOW DOWN ################")

	p1Name, p1Rank, p1Score := EvaluateHand(p1.Hand, board)
	p2Name, p2Rank, p2Score := EvaluateHand(p2.Hand, board)

	fmt.Printf("YOU: 【%s】\n", p1Name)
	printCardsAA(p1.Hand)

	fmt.Printf("CPU: 【%s】\n", p2Name)
	printCardsAA(p2.Hand)

	fmt.Println("###########################################")

	// ... (以下の勝敗判定ロジックはそのまま) ...
	// 勝敗判定
	win := false
	draw := false

	if p1Rank > p2Rank {
		win = true
	} else if p1Rank < p2Rank {
		win = false
	} else {
		// 役が同じ場合、数字で比較
		if p1Score > p2Score {
			win = true
		} else if p1Score < p2Score {
			win = false
		} else {
			draw = true
		}
	}

	if draw {
		fmt.Printf("引き分け！ Pot($%d)を分け合います。\n", pot)
		p1.Chips += pot / 2
		p2.Chips += pot / 2
	} else if win {
		fmt.Printf("あなたの勝ち！ Pot($%d)を獲得！\n", pot)
		p1.Chips += pot
	} else {
		fmt.Printf("CPUの勝ち... Pot($%d)を奪われました。\n", pot)
		p2.Chips += pot
	}
}
