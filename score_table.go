package scoretable

import (
	"errors"
	"fmt"
)

type playerScoreMap map[int]*playerRoundScore

const (
	BLACK_WIN = 1
	DRAW      = 0
	WHITE_WIN = -1
)

var (
	ErrUnknownResult = errors.New("unknown result")
	ErrUnknownPlayer = errors.New("unknown player")
)

type OptionFunc func(*Option)

type ScoreTable struct {
	m playerScoreMap
	*Option
}

type Score struct {
	PlayerId int
	NBW      float32
	SOS      int
	SOSOS    int
}
type Scores []*Score

type Option struct {
	bSOS    bool
	bSOSOS  bool
	drawNBW float32
}

type playerRoundScore struct {
	Round int     // 轮次
	score float32 // 本轮得分

	prev *playerRoundScore // 前驱
	next *playerRoundScore // 后继
	op   *playerRoundScore // 对手
}

func NewScoreTable(options ...OptionFunc) *ScoreTable {
	option := &Option{drawNBW: 0.5}
	for _, o := range options {
		o(option)
	}
	return &ScoreTable{Option: option}
}

func WithSOS() OptionFunc {
	return func(option *Option) {
		option.bSOS = true
	}
}

func WithSOSOS() OptionFunc {
	return func(option *Option) {
		option.bSOSOS = true
	}
}

func WithDrawScore(drawNBW float32) OptionFunc {
	return func(option *Option) {
		option.drawNBW = drawNBW
	}
}

func newPlayerRoundScore(round int, NBW float32) *playerRoundScore {
	var playerRoundScore playerRoundScore
	playerRoundScore.score = NBW
	playerRoundScore.Round = round

	return &playerRoundScore
}

// recordResult
func (s *ScoreTable) RecordResult(round int, blackPlayerId int, whitePlayerId int, result int) error {
	// 判断map中是否存在，不存在则新增，如果存在，则以拉链的方式向后追加
	if result != BLACK_WIN && result != WHITE_WIN && result != DRAW {
		return ErrUnknownResult
	}
	if s.m == nil {
		s.m = make(playerScoreMap)
	}
	addMemberRoundScore := func(playerId int, isBlack bool) *playerRoundScore {
		score := calculateNBW(result, isBlack, s.drawNBW)
		if mrs, ok := s.m[playerId]; ok {
			return mrs.addPlayerRoundScore(round, score)
		} else {
			s.m[playerId] = &playerRoundScore{Round: round, score: score}
			return s.m[playerId]
		}
	}
	blackP := addMemberRoundScore(blackPlayerId, true)
	whiteP := addMemberRoundScore(whitePlayerId, false)
	blackP.AddOpponent(whiteP)

	return nil
}

func calculateNBW(result int, isBlack bool, drawScore float32) float32 {
	if result == BLACK_WIN && isBlack {
		return 1
	} else if result == WHITE_WIN && !isBlack {
		return 1
	} else if result == DRAW {
		return drawScore
	}

	return 0
}

// 根据轮次获得对手分
func (m *playerRoundScore) getSosByRound(round int) int {
	head := m
	// 回退到头节点
	for head.prev != nil {
		head = head.prev
	}

	sos := 0
	for head != nil && head.Round <= round {
		sos += head.getOpponentNBWByRound(round)
		head = head.next
	}

	return sos
}

func (m *playerRoundScore) getSososByRound(round int) int {
	head := m
	sosos := 0
	// 回到链表头部
	for head.prev != nil {
		head = head.prev
	}
	for head != nil && head.Round <= round {
		sosos += head.op.getSosByRound(round)
		head = head.next
	}

	return sosos
}

func (m *playerRoundScore) getOpponentNBWByRound(round int) int {
	op := m.op
	return int(op.getNBWByRound(round))
}

func (m *playerRoundScore) getNBWByRound(round int) float32 {
	head := m
	NBW := float32(0)
	// 回到链表头部
	for head.prev != nil {
		head = head.prev
	}

	for head != nil && head.Round <= round {
		NBW += head.score
		head = head.next
	}

	return NBW
}

func (s *playerRoundScore) addPlayerRoundScore(round int, score float32) *playerRoundScore {
	node := s
	newScore := newPlayerRoundScore(round, score)
	// 回退到头节点
	for node.prev != nil {
		node = node.prev
	}
	if round > node.Round {
		for node.next != nil && node.next.Round < round {
			node = node.next
		}
		// 插入节点
		tempNext := node.next
		node.next = newScore
		newScore.next = tempNext
		newScore.prev = node
		if tempNext != nil {
			tempNext.prev = newScore
		}
	} else {
		for node.prev != nil && node.prev.Round > round {
			node = node.prev
		}
		// 插入节点
		tempPrev := node.prev
		node.prev = newScore
		newScore.prev = tempPrev
		newScore.next = node
		if tempPrev != nil {
			tempPrev.next = newScore
		}
	}

	return newScore
}

func (m *playerRoundScore) AddOpponent(op *playerRoundScore) *playerRoundScore {
	m.op = op
	op.op = m
	return m
}

func (m *playerRoundScore) clone() *playerRoundScore {
	newM := *m
	return &newM
}

func (s *ScoreTable) GetScoreTableByRound(round int) Scores {
	scores := make(Scores, 0)

	// 遍历map，之后获取用户的memberScore
	for k := range s.m {
		score, _ := s.createMemberScore(k, round)
		scores = append(scores, score)
	}

	OrderedBy(NBW, SOS, SOSOS, PlayerId).Sort(scores)

	return scores
}

func (s *ScoreTable) GetPlayerScoreByRound(playerId int, round int) (*Score, error) {
	return s.createMemberScore(playerId, round)
}

func (s *ScoreTable) createMemberScore(playerId int, round int) (*Score, error) {
	playerRoundScore, ok := s.m[playerId]
	if !ok {
		return nil, ErrUnknownPlayer
	}
	var memberScore Score
	memberScore.PlayerId = playerId
	memberScore.NBW = playerRoundScore.getNBWByRound(round)
	if s.bSOS {
		memberScore.SOS = playerRoundScore.getSosByRound(round)
	}
	if s.bSOSOS {
		memberScore.SOSOS = playerRoundScore.getSososByRound(round)
	}

	return &memberScore, nil
}

func PrintNBW(m *playerRoundScore, id string, round int) {
	info := ""
	nbw := m.getNBWByRound(round)
	info += fmt.Sprintf("选手：%s  当前轮次为：%d，当前NBW为：%d\n", id, round, nbw)
	fmt.Println(info)
}

func PrintSOS(m *playerRoundScore, id string, round int) {
	sos := m.getSosByRound(round)
	info := fmt.Sprintf("选手：%s  当前轮次为：%d，当前Sos为：%d\n", id, round, sos)
	fmt.Println(info)
}

func PrintSOSOS(m *playerRoundScore, id string, round int) {
	sos := m.getSososByRound(round)
	info := fmt.Sprintf("选手：%s  当前轮次为：%d，当前Sosos为：%d\n", id, round, sos)
	fmt.Println(info)
}