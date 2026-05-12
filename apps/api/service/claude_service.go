package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type ClaudeService struct{}

func NewClaudeService() *ClaudeService {
	return &ClaudeService{}
}

type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// GenerateExplanation calls Claude API to generate a Bahasa Indonesia explanation
func (s *ClaudeService) GenerateExplanation(maqamName, maqamArabic string, candidates interface{}) string {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Warn("ANTHROPIC_API_KEY not set, returning default explanation")
		return s.defaultExplanation(maqamName, maqamArabic)
	}

	candidatesJSON, _ := json.Marshal(candidates)

	prompt := fmt.Sprintf(`Kamu adalah ahli musik Arab dan maqam yang menjelaskan dalam Bahasa Indonesia untuk komunitas banjari.

Hasil analisis mendeteksi maqam: %s (%s)
Kandidat lengkap: %s

Berikan penjelasan singkat (maks 3 paragraf) yang mencakup:
1. Karakteristik maqam ini (interval, nada khas)
2. Nuansa emosional yang dihasilkan
3. Tips praktis untuk musisi banjari

Gunakan bahasa yang mudah dipahami, hangat, dan informatif. Jangan gunakan format markdown, cukup paragraf biasa.`, maqamName, maqamArabic, string(candidatesJSON))

	reqBody := claudeRequest{
		Model:     "claude-sonnet-4-6",
		MaxTokens: 500,
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Errorf("Claude request build error: %v", err)
		return s.defaultExplanation(maqamName, maqamArabic)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Claude API error: %v", err)
		return s.defaultExplanation(maqamName, maqamArabic)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Claude API returned %d: %s", resp.StatusCode, string(body))
		return s.defaultExplanation(maqamName, maqamArabic)
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil || len(claudeResp.Content) == 0 {
		return s.defaultExplanation(maqamName, maqamArabic)
	}

	return strings.TrimSpace(claudeResp.Content[0].Text)
}

func (s *ClaudeService) defaultExplanation(name, arabic string) string {
	return fmt.Sprintf(
		"Audio yang dianalisis terdeteksi menggunakan maqam %s (%s). "+
			"Maqam ini merupakan salah satu tangga nada yang umum digunakan dalam musik Arab dan tradisi banjari. "+
			"Untuk penjelasan lebih detail, pastikan ANTHROPIC_API_KEY telah dikonfigurasi.",
		name, arabic,
	)
}
