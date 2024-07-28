package useragent

import (
	"embed"
	"encoding/json"
	"io"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type UserAgent struct {
	Device    string `json:"device"`
	UserAgent string `json:"useragent"`
}

// Встраиваем файл в бинарь
//go:embed user-agents.json
var userAgentsFS embed.FS

// RandomUserAgentAndDevice returns a random UserAgent struct
func Random(zl *zap.Logger) (*UserAgent, error) {
	rand.NewSource(time.Now().UnixNano())

	file, err := userAgentsFS.Open("user-agents.json")
	if err != nil {
		zl.Error("failed to open file", zap.String("err:", err.Error()))
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		zl.Error("failed to read file", zap.String("err:", err.Error()))
		return nil, err
	}

	var userAgents []UserAgent
	err = json.Unmarshal(bytes, &userAgents)
	if err != nil {
		zl.Error("failed to unmarshal JSON", zap.String("err:", err.Error()))
		return nil, err
	}

	randomIndex := rand.Intn(len(userAgents))
	randomUserAgent := userAgents[randomIndex]

	return &randomUserAgent, nil
}
