package domain

import (
	"errors"
	"math/rand"
	"time"
)

type GameStatus string

const (
	GameStatusWaiting  GameStatus = "waiting"
	GameStatusPlaying  GameStatus = "playing"
	GameStatusFinished GameStatus = "finished"
)

type Game struct {
	room         Room
	currentTheme string
	status       GameStatus
	round        int
	voteState    *VoteState
}

type VoteState struct {
	targetPlayerID uint
	approvals      map[uint]bool
	rejections     map[uint]bool
	timeoutAt      time.Time
}

func NewGame(room Room) *Game {
	return &Game{
		room:   room,
		status: GameStatusWaiting,
		round:  1,
	}
}

// Getters
func (g *Game) Room() Room {
	return g.room
}

func (g *Game) CurrentTheme() string {
	return g.currentTheme
}

func (g *Game) Status() GameStatus {
	return g.status
}

func (g *Game) Round() int {
	return g.round
}

func (g *Game) VoteState() *VoteState {
	return g.voteState
}

// VoteStateのゲッターを追加
func (v *VoteState) TargetPlayerID() uint {
	return v.targetPlayerID
}

func (v *VoteState) Approvals() map[uint]bool {
	return v.approvals
}

func (v *VoteState) Rejections() map[uint]bool {
	return v.rejections
}

func (v *VoteState) TimeoutAt() time.Time {
	return v.timeoutAt
}

// Game Methods
func (g *Game) SetTheme(theme string) {
	g.currentTheme = theme
}

func (g *Game) StartPlaying() {
	g.status = GameStatusPlaying
}

func (g *Game) Finish() {
	g.status = GameStatusFinished
}

func (g *Game) InitiateVote(targetPlayerID uint) error {
	if g.status != GameStatusPlaying {
		return errors.New("ゲームがプレイ中ではありません")
	}

	if g.voteState != nil {
		return errors.New("他の投票が進行中です")
	}

	g.voteState = &VoteState{
		targetPlayerID: targetPlayerID,
		approvals:     make(map[uint]bool),
		rejections:    make(map[uint]bool),
		timeoutAt:     time.Now().Add(30 * time.Second),
	}

	return nil
}

func (g *Game) ProcessVote(voterID uint, approve bool) error {
	if g.voteState == nil {
		return errors.New("投票が開始されていません")
	}

	if time.Now().After(g.voteState.timeoutAt) {
		return errors.New("投票時間が終了しています")
	}

	if approve {
		g.voteState.approvals[voterID] = true
	} else {
		g.voteState.rejections[voterID] = true
	}

	return nil
}

func (g *Game) ClearVoteState() {
	g.voteState = nil
}

func (g *Game) IsVoteSuccessful() bool {
	if g.voteState == nil {
		return false
	}
	return len(g.voteState.rejections) == 0
}

func (g *Game) HasAllVoted() bool {
	if g.voteState == nil {
		return false
	}

	players := g.room.GetPlayers()
	votedCount := len(g.voteState.approvals) + len(g.voteState.rejections)
	return votedCount >= len(players)-1 // 投票対象者を除く
}

func (g *Game) IsVoteTimeout() bool {
	if g.voteState == nil {
		return false
	}
	return time.Now().After(g.voteState.timeoutAt)
}

func (g *Game) Theme() string {
	return g.currentTheme
}

func (g *Game) SetPlayerReady(userID uint) error {
	for _, player := range g.room.GetPlayers() {
		user := player.User()
		if user.ID() == userID {
			player.SetIsReady(true)
			return nil
		}
	}
	return errors.New("プレイヤーが見つかりません")
}

func (g *Game) DealCards() error {
	players := g.room.GetPlayers()
	numbers := make([]int, len(players))
	for i := range numbers {
		numbers[i] = i + 1
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	for i, player := range players {
		player.SetCardValue(numbers[i])
	}
	return nil
}

func (g *Game) Start() {
	g.status = GameStatusPlaying
} 