package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// --- 定数・構造体 ---

type Card struct {
	Suit string
	Rank string
	Val  int // 強さ比較用 (2=2, ... A=14)
}

type Player struct {
	Name  string
	Hand  []Card
	Chips int
	IsCPU bool
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
	human := &Player{Name: "You", Chips: 1000, IsCPU: false}
	cpu := &Player{Name: "CPU", Chips: 1000, IsCPU: true}

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

		// ベットラウンド（降りたら即終了）
		if !bettingRound(p1, p2, &pot, phase) {
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
	fmt.Println("\n--------------------------------")
	fmt.Printf("PHASE: %s  |  POT: $%d\n", phase, pot)
	fmt.Printf("BOARD: %v\n", board)
	fmt.Printf("YOU  : %v  (Chips: $%d)\n", p1.Hand, p1.Chips)
	fmt.Printf("CPU  : [??] [??]  (Chips: $%d)\n", p2.Chips) // CPUは見せない
	fmt.Println("--------------------------------")
}

// ベット処理（戻り値 false = 誰かが降りてラウンド終了）
func bettingRound(p1, p2 *Player, pot *int, phase string) bool {
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
	cpuAction, cpuAmount := getCpuAction(p2, *pot, amount, phase)
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

// 簡易AI: 強いカードや後半のフェーズなら強気にでる
func getCpuAction(p *Player, pot int, playerBet int, phase string) (string, int) {
	// 1. 手札の強さを簡易評価 (数字の合計)
	score := p.Hand[0].Val + p.Hand[1].Val
	
	// 2. 確率でアクション決定
	randVal := rand.Intn(100)

	// 相手がベットしてきた場合
	if playerBet > 0 {
		if score < 10 && randVal < 80 { // 手札が弱ければ80%降りる
			return "FOLD", 0
		}
		// コールする（相手の額に合わせる処理は省略し、とりあえず同額払うとする）
		bet := playerBet
		if p.Chips < bet { bet = p.Chips }
		return "CALL", bet
	}

	// 相手がチェックの場合
	if score > 20 || randVal > 70 { // 手札が強いか、30%の確率でレイズ
		bet := 50
		if p.Chips < bet { bet = p.Chips }
		return "RAISE", bet
	}

	return "CHECK", 0
}

func showdown(p1, p2 *Player, board []Card, pot int) {
	fmt.Println("\n=== SHOW DOWN ===")
	fmt.Printf("YOU: %v  (Score: %d)\n", p1.Hand, p1.Hand[0].Val+p1.Hand[1].Val)
	fmt.Printf("CPU: %v  (Score: %d)\n", p2.Hand, p2.Hand[0].Val+p2.Hand[1].Val)

	// 簡易勝敗判定（カードの数字の合計が大きい方が勝ち）
	// 本来はフラッシュやストレートの判定ロジックが必要
	score1 := p1.Hand[0].Val + p1.Hand[1].Val
	score2 := p2.Hand[0].Val + p2.Hand[1].Val

	if score1 > score2 {
		fmt.Printf("あなたの勝ち！ Pot($%d)を獲得！\n", pot)
		p1.Chips += pot
	} else if score2 > score1 {
		fmt.Printf("CPUの勝ち... Pot($%d)を奪われました。\n", pot)
		p2.Chips += pot
	} else {
		fmt.Printf("引き分け！ Pot($%d)を分け合います。\n", pot)
		p1.Chips += pot / 2
		p2.Chips += pot / 2
	}
}