// (C) Copyright 2023-2024 Hewlett Packard Enterprise Development LP

package dbtesting

import (
	"math/rand"
	"time"
)

const (
	timeOffsetWindowMin int = -300
	timeOffsetWindowMax int = 300
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		//nolint:gosec // crytpgraphy strength random generator not required for test data
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomInt(max int) int {
	//nolint:gosec // crytpgraphy strength random generator not required for test data
	return rand.Intn(max)
}

func RandomInt64(max int64) int64 {
	//nolint:gosec // crytpgraphy strength random generator not required for test data
	return rand.Int63n(max)
}

func RandomTime() time.Time {
	offset := RandomInt(timeOffsetWindowMax-timeOffsetWindowMin) - timeOffsetWindowMax
	return time.Now().UTC().Add(time.Duration(offset) * time.Minute)
}
