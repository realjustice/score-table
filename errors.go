package scoretable

import "errors"

var (
	ErrUnknownResult = errors.New("unknown result")
	ErrUnknownPlayer = errors.New("unknown player")
)

type ErrMissScoreNode struct {
	Err  string // description of the error
	Node string
}

func NewErrMissScoreNode(node string) *ErrMissScoreNode {
	return &ErrMissScoreNode{Err: "cannot order by this method", Node: node}
}

func (e *ErrMissScoreNode) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Node + ":" + "missed" + e.Err
}
