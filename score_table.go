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
	BOTH_LOSE = 2
	BOTH_WIN  = 3
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
	SOS      float32
	SOSOS    float32
	Rank     int
}
type Scores []*Score

type Option struct {
	bSOS    bool
	bSOSOS  bool
	bRank   bool
	drawNBW float32
}

type playerRoundScore struct {
	Round       int     // 轮次
	score       float32 // 本轮得分
	isByePlayer bool    // 是否是轮空选手

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

func WithRank() OptionFunc {
	return func(option *Option) {
		option.bRank = true
	}
}

func WithDrawScore(drawNBW float32) OptionFunc {
	return func(option *Option) {
		option.drawNBW = drawNBW
	}
}

func newPlayerRoundScore(round int, NBW float32, isByePlayer bool) *playerRoundScore {
	var playerRoundScore playerRoundScore
	playerRoundScore.score = NBW
	playerRoundScore.Round = round
	playerRoundScore.isByePlayer = isByePlayer

	return &playerRoundScore
}

// recordResult
func (s *ScoreTable) RecordResult(round int, blackPlayerId int, whitePlayerId int, result int) error {
	// 判断map中是否存在，不存在则新增，如果存在，则以拉链的方式向后追加
	if result != BLACK_WIN && result != WHITE_WIN && result != DRAW && result != BOTH_LOSE && result != BOTH_WIN {
		return ErrUnknownResult
	}
	if s.m == nil {
		s.m = make(playerScoreMap)
	}
	addMemberRoundScore := func(playerId int, isBlack bool) *playerRoundScore {
		score := calculateNBW(result, isBlack, s.drawNBW)
		if mrs, ok := s.m[playerId]; ok {
			return mrs.addPlayerRoundScore(round, score, playerId == 0)
		} else {
			s.m[playerId] = newPlayerRoundScore(round, score, playerId == 0)
			return s.m[playerId]
		}
	}
	blackP := addMemberRoundScore(blackPlayerId, true)
	whiteP := addMemberRoundScore(whitePlayerId, false)
	blackP.addOpponent(whiteP)

	return nil
}

func calculateNBW(result int, isBlack bool, drawScore float32) (score float32) {
	if result == BLACK_WIN && isBlack {
		score = 1
	} else if result == WHITE_WIN && !isBlack {
		score = 1
	} else if result == DRAW {
		score = drawScore
	} else if result == BOTH_LOSE {
		score = 0
	} else if result == BOTH_WIN {
		score = 1
	}

	return
}

// 根据轮次获得对手分
func (m *playerRoundScore) getSosByRound(round int) float32 {
	head := m
	// 回退到头节点
	for head.prev != nil {
		head = head.prev
	}

	var sos float32
	for head != nil && head.Round <= round {
		if !head.isByePlayer {
			sos += head.getOpponentNBWByRound(round)
		}

		head = head.next
	}

	return sos
}

func (m *playerRoundScore) getSososByRound(round int) float32 {
	head := m
	var sosos float32
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

func (m *playerRoundScore) getOpponentNBWByRound(round int) float32 {
	op := m.op

	return op.getNBWByRound(round)
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

func (s *playerRoundScore) addPlayerRoundScore(round int, score float32, isByePlayer bool) *playerRoundScore {
	node := s
	newScore := newPlayerRoundScore(round, score, isByePlayer)
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

func (m *playerRoundScore) addOpponent(op *playerRoundScore) *playerRoundScore {
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
		if !s.m[k].isByePlayer {
			scores = append(scores, score)
		}
	}

	OrderedBy(NBW, SOS, SOSOS, PlayerId).Sort(scores)
	// 设置排名
	if s.bRank {
		setRank(scores)
	}

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

func setRank(scores Scores) {
	rank := 0
	isSame := func(s1 *Score, s2 *Score) bool {
		return s1.NBW == s2.NBW && s1.SOS == s2.SOS && s1.SOSOS == s2.SOSOS
	}
	var lastScore Score
	sameCount := 0
	for _, score := range scores {
		if isSame(&lastScore, score) {
			score.Rank = rank
			sameCount++
			continue
		}
		rank = rank + sameCount + 1
		sameCount = 0
		score.Rank = rank

		lastScore = *score
	}
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
