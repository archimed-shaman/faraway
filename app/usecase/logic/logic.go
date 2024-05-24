package logic

import "context"

type UserLogic struct{}

func New() *UserLogic {
	return &UserLogic{}
}

func (l *UserLogic) GetQuote(_ context.Context) (string, error) {
	//nolint:lll // Just mocked logic
	quote := "And all saints who remember to keep and do these sayings, walking in obedience to the commandments, shall receive health in their navel and marrow to their bones; and shall find wisdom and great treasures of knowledge, even hidden treasures."

	return quote, nil
}
