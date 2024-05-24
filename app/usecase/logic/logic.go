package logic

import (
	"context"
	"math/rand"
	"time"
)

type UserLogic struct {
	quotes []string
}

func New() *UserLogic {
	return &UserLogic{
		quotes: []string{
			//nolint:lll // Just mocked logic
			"And all saints who remember to keep and do these sayings, walking in obedience to the commandments, shall receive health in their navel and marrow to their bones; and shall find wisdom and great treasures of knowledge, even hidden treasures.",
			"The only limit to our realization of tomorrow is our doubts of today.",
			"The future belongs to those who believe in the beauty of their dreams.",
			"Do not watch the clock. Do what it does. Keep going.",
			"Keep your face always toward the sunshine—and shadows will fall behind you.",
		},
	}
}

func (l *UserLogic) GetQuote(_ context.Context) (string, error) {
	//nolint:gosec // no any need to use strong secure generator here
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	quote := l.quotes[r.Intn(len(l.quotes))]

	return quote, nil
}
