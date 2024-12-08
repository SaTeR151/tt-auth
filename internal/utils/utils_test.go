package utils

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCheckHost(t *testing.T) {
	tests := []struct {
		giveHost   string
		wantResult bool
		wantErr    error
	}{
		{
			giveHost:   "localhost:8080",
			wantResult: true,
			wantErr:    nil,
		},
		{
			giveHost:   "0.0.0.0:8080",
			wantResult: false,
			wantErr:    nil,
		},
	}
	err := godotenv.Load()
	assert.NotEqual(t, err, nil, err)

	for _, testTask := range tests {
		claims := &jwt.MapClaims{
			"Host": "localhost:8080",
		}
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		aToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			assert.NotEqual(t, err, nil, err)
		}

		check, err := CheckHost(aToken, testTask.giveHost)

		assert.Equal(t, testTask.wantResult, check, "хост проверился неверно")
		assert.Equal(t, testTask.wantErr, err, "вернулась ошибка, которой не должно быть")
	}
}
