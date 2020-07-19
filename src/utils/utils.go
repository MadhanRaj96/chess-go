package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
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
	out, err := exec.Command("uuidgen").Output()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	return base64.RawURLEncoding.EncodeToString(out)
}

// GetColor of player
func GetColor() string {
	r := randomInt(1, 9)
	if r%2 == 0 {
		return "white"
	}
	return "black"
}
