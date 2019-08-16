package utils

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func initRand() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestMain(m *testing.M) {
	initRand()
	os.Exit(m.Run())
}

func TestGen_String(t *testing.T) {

	t.Run("Should return same hash for same string", func(t *testing.T) {
		gen := NewGen()
		for i := 0; i < 10; i++ {
			testStr := RandStringRunes(16)
			hash1 := gen.String(testStr)
			hash2 := gen.String(testStr)

			if hash1 != hash2 {
				t.Errorf("String() doesn't return same hash for one string.\nHash1: %d\tHash2: %d", hash1, hash2)
			}
		}
	})

	t.Run("Should return different hash when two separate generators are used", func(t *testing.T) {
		testStr := RandStringRunes(16)
		gen := NewGen()
		hash1 := gen.String(testStr)
		gen = NewGen()
		hash2 := gen.String(testStr)
		if hash1 == hash2 {
			t.Errorf("String() doesn't return diffrent hash.\nHash1: %d\tHash2: %d", hash1, hash2)
		}
	})
}
