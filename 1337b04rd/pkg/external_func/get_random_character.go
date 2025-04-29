package externalfunc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type RMCharacter struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func GetRandomCharacter() (*RMCharacter, error) {
	const maxCharacters = 826 // обновить при необходимости
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(maxCharacters) + 1

	resp, err := http.Get(fmt.Sprintf("https://rickandmortyapi.com/api/character/%d", id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var character RMCharacter
	if err := json.NewDecoder(resp.Body).Decode(&character); err != nil {
		return nil, err
	}
	return &character, nil
}
