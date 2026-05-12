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
	explanations := map[string]string{
		"Bayati": "Maqam Bayati (بياتي) adalah salah satu maqam paling populer dalam musik Arab dan sering dijumpai dalam seni banjari. " +
			"Ciri khasnya terletak pada interval tiga perempat nada (3/4 tone) antara nada pertama dan kedua, yang menghasilkan nuansa emosional yang mendalam, penuh kerinduan, dan spiritual. " +
			"Bayati sering digunakan dalam pembacaan Al-Quran, sholawat, dan lagu-lagu bertema cinta ilahi karena karakternya yang lembut namun mengena di hati.",
		"Rast": "Maqam Rast (راست) dianggap sebagai 'raja' dari seluruh maqam dalam tradisi musik Arab. " +
			"Rast memiliki karakter yang megah, cerah, dan penuh keagungan. Interval khasnya menggunakan tiga perempat nada yang memberikan warna unik berbeda dari tangga nada mayor Barat. " +
			"Dalam kesenian banjari, Rast sering digunakan untuk pembukaan acara atau lagu-lagu yang bersifat pujian dan kegembiraan.",
		"Hijaz": "Maqam Hijaz (حجاز) mudah dikenali dari interval augmented second (satu setengah nada) antara nada kedua dan ketiga, yang menghasilkan nuansa eksotis dan dramatis. " +
			"Maqam ini sangat identik dengan suasana Timur Tengah dan sering diasosiasikan dengan kerinduan, kesedihan yang indah, dan spiritualitas mendalam. " +
			"Dalam banjari, Hijaz adalah pilihan favorit untuk sholawat yang bertema rindu kepada Rasulullah SAW.",
		"Saba": "Maqam Saba (صبا) dikenal sebagai maqam yang paling emosional dan menyentuh perasaan. " +
			"Karakteristiknya yang melankolis dan penuh keharuan menjadikannya pilihan utama untuk mengekspresikan kesedihan, taubat, dan doa yang tulus. " +
			"Interval khasnya menciptakan ketegangan harmonis yang membuat pendengar merasa terhanyut dalam suasana reflektif dan khusyuk.",
		"Nahawand": "Maqam Nahawand (نهاوند) memiliki kesamaan dengan tangga nada minor natural dalam musik Barat, namun dengan ornamentasi khas Arab yang membuatnya lebih kaya. " +
			"Karakternya romantis, lembut, dan sedikit melankolis — cocok untuk lagu-lagu bertema cinta dan kerinduan. " +
			"Dalam banjari, Nahawand sering digunakan untuk sholawat yang bernuansa syahdu dan penuh perasaan.",
		"Ajam": "Maqam Ajam (عجم) memiliki struktur yang mirip dengan tangga nada mayor dalam musik Barat, sehingga terdengar familiar dan ceria. " +
			"Karakter Ajam yang cerah, optimis, dan penuh semangat menjadikannya cocok untuk lagu-lagu perayaan dan kegembiraan. " +
			"Dalam konteks banjari, Ajam sering dipakai untuk marawis dan lagu-lagu bertema syukur serta kebahagiaan.",
		"Kurd": "Maqam Kurd (كرد) memiliki karakter yang sederhana namun sangat ekspresif, mirip dengan tangga nada minor Barat pada derajat tertentu. " +
			"Nuansa yang dihasilkan cenderung tenang, meditatif, dan penuh kekhusyukan. " +
			"Kurd sering digunakan dalam tilawah Al-Quran dan dzikir karena sifatnya yang menenangkan jiwa dan mudah diterima oleh telinga.",
		"Jiharkah": "Maqam Jiharkah (جهاركاه) memiliki karakter yang kuat, tegas, dan berwibawa. " +
			"Maqam ini sering dikaitkan dengan keagungan dan kemegahan, menjadikannya pilihan yang tepat untuk lagu-lagu pujian dan sholawat yang bersifat agung. " +
			"Dalam tradisi banjari, Jiharkah memberikan nuansa yang membangkitkan semangat dan rasa percaya diri.",
	}

	if desc, ok := explanations[name]; ok {
		return desc
	}

	return fmt.Sprintf(
		"Audio yang dianalisis terdeteksi menggunakan maqam %s (%s). "+
			"Maqam ini merupakan salah satu tangga nada yang umum digunakan dalam musik Arab dan tradisi banjari. "+
			"Setiap maqam memiliki karakteristik interval dan nuansa emosional yang unik, "+
			"menciptakan warna musik yang khas dan bermakna dalam setiap penampilan.",
		name, arabic,
	)
}
