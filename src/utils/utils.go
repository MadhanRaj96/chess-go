package utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// RandomString of A-Z chars with len = l
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

// GenerateGameID -  new Game ID on each request
func GenerateGameID() string {
	out := uuid.New()
	uuid := strings.Replace(out.String(), "-", "", -1)
	fmt.Println(uuid)

	return base64.RawURLEncoding.EncodeToString([]byte(uuid))
}

// GetColor of player
func GetColor() string {
	r := randomInt(1, 9)
	if r%2 == 0 {
		return "white"
	}
	return "black"
}
