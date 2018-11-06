package types

import "fmt"

type OperatorMessage struct {
	MessageType string
	Matches     *Matches
	ErrorType   string
}

func (m *OperatorMessage) String() string {
	if m.ErrorType != "" {
		return fmt.Sprintf("%v: %v (%v)", m.MessageType, m.Matches.String(), m.ErrorType)
	}

	return fmt.Sprintf("%v: %v", m.MessageType, m.Matches.String())
}

type PendingTradeBatch struct {
	Matches *Matches
}
