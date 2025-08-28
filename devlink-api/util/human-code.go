package util

import (
    "fmt"
    "math/rand"
    "time"
)

var (
    adjectives = []string{"swift", "silent", "brave", "bright", "cool", "eager"}
    nouns      = []string{"badger", "falcon", "river", "stone", "moon", "wave"}
    seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GenerateHumanCode creates a memorable code like "3-swift-badger".
func GenerateHumanCode() string {
    num := seededRand.Intn(10)
    adj := adjectives[seededRand.Intn(len(adjectives))]
    noun := nouns[seededRand.Intn(len(nouns))]
    return fmt.Sprintf("%d-%s-%s", num, adj, noun)
}