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
			"In the end, we will remember not the words of our enemies, but the silence of our friends.",
			"To handle yourself, use your head; to handle others, use your heart.",
			"Life is 10% what happens to us and 90% how we react to it.",
			"The best way to predict the future is to create it.",
			"Do not dwell in the past, do not dream of the future, concentrate the mind on the present moment.",
			"The only way to do great work is to love what you do.",
			"Success is not the key to happiness. Happiness is the key to success.",
		},
	}
}

func (l *UserLogic) GetQuote(_ context.Context) (string, error) {
	//nolint:gosec // no any need to use strong secure generator here
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	quote := l.quotes[r.Intn(len(l.quotes))]

	return quote, nil
}
