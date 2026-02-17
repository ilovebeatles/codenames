package game

import (
	"math/rand"

	"codenames/internal/model"
)

// GenerateBoard creates 25 cards for a new game.
// firstTeam gets 9 cards, the other gets 8, 7 neutral, 1 assassin.
func GenerateBoard(gameID string, firstTeam model.Team) []model.Card {
	words := pickRandomWords(25)

	second := firstTeam.Opposite()
	types := make([]model.CardType, 25)

	// 9 for first team
	idx := 0
	for i := 0; i < 9; i++ {
		types[idx] = model.CardType(firstTeam)
		idx++
	}
	// 8 for second team
	for i := 0; i < 8; i++ {
		types[idx] = model.CardType(second)
		idx++
	}
	// 7 neutral
	for i := 0; i < 7; i++ {
		types[idx] = model.CardTypeNeutral
		idx++
	}
	// 1 assassin
	types[idx] = model.CardTypeAssassin

	// shuffle types
	rand.Shuffle(len(types), func(i, j int) {
		types[i], types[j] = types[j], types[i]
	})

	cards := make([]model.Card, 25)
	for i := 0; i < 25; i++ {
		cards[i] = model.Card{
			GameID:   gameID,
			Word:     words[i],
			CardType: types[i],
			Position: i,
		}
	}
	return cards
}

func pickRandomWords(n int) []string {
	perm := rand.Perm(len(RussianWords))
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = RussianWords[perm[i]]
	}
	return words
}
