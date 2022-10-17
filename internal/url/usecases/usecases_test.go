package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortLink(t *testing.T) {
	testsPositive := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Positive test",
			url:  "https://google.com",
			want: "aHR0cHM6Ly9nb29nbGUuY29t",
		},
		{
			name: "Positive test 2",
			url:  "https://amazon.com",
			want: "aHR0cHM6Ly9hbWF6b24uY29t",
		},
	}

	testsNegative := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Negative Test 1",
			url:  "https://google.com",
			want: "aHR0cHM6Ly9nb29fsdUuY29t",
		},
		{
			name: "Negative test 2",
			url:  "https://amazon.com",
			want: "aHR0cHM6Ly9hbWFddfd4uY29t",
		},
	}

	for _, test := range testsPositive {
		t.Run(test.name, func(t *testing.T) {
			res := GenerateShortLink([]byte(test.url))
			assert.Equal(t, test.want, res)
		})
	}

	for _, test := range testsNegative {
		t.Run(test.name, func(t *testing.T) {
			res := GenerateShortLink([]byte(test.url))
			assert.NotEqual(t, test.want, res)
		})
	}
}
