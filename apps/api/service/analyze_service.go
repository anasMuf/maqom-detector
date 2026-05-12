package service

import (
	"api/model"
	"api/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AnalyzeService struct {
	analysisRepo *repository.AnalysisRepository
	sessionRepo  *repository.SessionRepository
	claudeService *ClaudeService
}

func NewAnalyzeService(
	analysisRepo *repository.AnalysisRepository,
	sessionRepo *repository.SessionRepository,
	claudeService *ClaudeService,
) *AnalyzeService {
	return &AnalyzeService{
		analysisRepo:  analysisRepo,
		sessionRepo:   sessionRepo,
		claudeService: claudeService,
	}
}

// AnalyzerResponse from Python analyzer
type AnalyzerResponse struct {
	Top3Candidates []struct {
		MaqamID         string  `json:"maqam_id"`
		NameLatin       string  `json:"name_latin"`
		NameArabic      string  `json:"name_arabic"`
		ConfidenceScore float64 `json:"confidence_score"`
		Rank            int     `json:"rank"`
	} `json:"top3_candidates"`
	ConfidenceLabel string    `json:"confidence_label"`
	AudioQuality    string    `json:"audio_quality"`
	ProcessingMs    int       `json:"processing_ms"`
	PCP             []float64 `json:"pcp"`
}

// StartYoutubeAnalysis creates a new analysis and starts async processing
func (s *AnalyzeService) StartYoutubeAnalysis(sessionID uuid.UUID, url string, segStart, segDuration int) (*model.Analysis, error) {
	// Ensure session exists
	if _, err := s.sessionRepo.FindOrCreate(sessionID); err != nil {
		return nil, fmt.Errorf("gagal buat session: %w", err)
	}

	analysis := &model.Analysis{
		SessionID:   sessionID,
		InputType:   "youtube",
		InputSource: url,
		Status:      "pending",
	}

	if err := s.analysisRepo.Create(analysis); err != nil {
		return nil, fmt.Errorf("gagal buat record analisis: %w", err)
	}

	// Async processing
	go s.processYoutubeAnalysis(analysis.ID, url, segStart, segDuration)

	return analysis, nil
}

// StartUploadAnalysis creates analysis from uploaded file
func (s *AnalyzeService) StartUploadAnalysis(sessionID uuid.UUID, filename string, fileContent []byte, sourceType string) (*model.Analysis, error) {
	if _, err := s.sessionRepo.FindOrCreate(sessionID); err != nil {
		return nil, fmt.Errorf("gagal buat session: %w", err)
	}

	analysis := &model.Analysis{
		SessionID:   sessionID,
		InputType:   sourceType,
		InputSource: filename,
		Status:      "pending",
	}

	if err := s.analysisRepo.Create(analysis); err != nil {
		return nil, fmt.Errorf("gagal buat record analisis: %w", err)
	}

	go s.processUploadAnalysis(analysis.ID, fileContent, sourceType)

	return analysis, nil
}

func (s *AnalyzeService) processYoutubeAnalysis(analysisID uuid.UUID, url string, segStart, segDuration int) {
	s.analysisRepo.UpdateStatus(analysisID, "processing")

	analyzerURL := os.Getenv("ANALYZER_BASE_URL")
	if analyzerURL == "" {
		analyzerURL = "http://localhost:8000"
	}

	// Build multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("source_type", "youtube")
	writer.WriteField("url", url)
	writer.WriteField("segment_start", fmt.Sprintf("%d", segStart))
	writer.WriteField("segment_duration", fmt.Sprintf("%d", segDuration))
	writer.WriteField("mode", "normal")
	writer.Close()

	s.callAnalyzerAndFinalize(analysisID, analyzerURL, writer.FormDataContentType(), body)
}

func (s *AnalyzeService) processUploadAnalysis(analysisID uuid.UUID, fileContent []byte, sourceType string) {
	s.analysisRepo.UpdateStatus(analysisID, "processing")

	analyzerURL := os.Getenv("ANALYZER_BASE_URL")
	if analyzerURL == "" {
		analyzerURL = "http://localhost:8000"
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("source_type", sourceType)

	mode := "normal"
	if sourceType == "humming" {
		mode = "humming"
	} else if sourceType == "microphone" {
		mode = "microphone"
	}
	writer.WriteField("mode", mode)
	writer.WriteField("segment_start", "0")
	writer.WriteField("segment_duration", "120")

	// Add file part
	part, _ := writer.CreateFormFile("file", "audio.wav")
	part.Write(fileContent)
	writer.Close()

	s.callAnalyzerAndFinalize(analysisID, analyzerURL, writer.FormDataContentType(), body)
}

func (s *AnalyzeService) callAnalyzerAndFinalize(analysisID uuid.UUID, analyzerURL, contentType string, body *bytes.Buffer) {
	client := &http.Client{Timeout: 5 * time.Minute}

	req, err := http.NewRequest("POST", analyzerURL+"/internal/analyze", body)
	if err != nil {
		s.failAnalysis(analysisID, "ANALYSIS_FAILED", "Gagal membuat request ke analyzer")
		return
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Analyzer request failed for %s: %v", analysisID, err)
		s.failAnalysis(analysisID, "ANALYSIS_FAILED", "Gagal menghubungi analyzer service")
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Detail struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"detail"`
		}
		json.Unmarshal(respBody, &errResp)

		code := errResp.Detail.Code
		msg := errResp.Detail.Message
		if code == "" {
			code = "ANALYSIS_FAILED"
		}
		if msg == "" {
			msg = fmt.Sprintf("Analyzer returned status %d", resp.StatusCode)
		}
		s.failAnalysis(analysisID, code, msg)
		return
	}

	var analyzerResp AnalyzerResponse
	if err := json.Unmarshal(respBody, &analyzerResp); err != nil {
		s.failAnalysis(analysisID, "ANALYSIS_FAILED", "Gagal parsing response analyzer")
		return
	}

	if len(analyzerResp.Top3Candidates) == 0 {
		s.failAnalysis(analysisID, "ANALYSIS_FAILED", "Tidak ada kandidat maqam ditemukan")
		return
	}

	// Generate Claude explanation
	topMaqam := analyzerResp.Top3Candidates[0]
	explanation := s.claudeService.GenerateExplanation(topMaqam.NameLatin, topMaqam.NameArabic, analyzerResp.Top3Candidates)

	// Build candidates
	var candidates []model.AnalysisCandidate
	for _, c := range analyzerResp.Top3Candidates {
		candidates = append(candidates, model.AnalysisCandidate{
			MaqamID:         c.MaqamID,
			ConfidenceScore: c.ConfidenceScore,
			Rank:            c.Rank,
		})
	}

	result := &repository.AnalysisResult{
		DetectedMaqamID: topMaqam.MaqamID,
		ConfidenceScore: topMaqam.ConfidenceScore,
		ConfidenceLabel: analyzerResp.ConfidenceLabel,
		ExplanationText: explanation,
		AudioQuality:    analyzerResp.AudioQuality,
		ProcessingMs:    analyzerResp.ProcessingMs,
		Candidates:      candidates,
	}

	if err := s.analysisRepo.UpdateCompleted(analysisID, result); err != nil {
		log.Errorf("Failed to save analysis result %s: %v", analysisID, err)
		s.failAnalysis(analysisID, "ANALYSIS_FAILED", "Gagal menyimpan hasil analisis")
	}
}

func (s *AnalyzeService) failAnalysis(id uuid.UUID, code, message string) {
	if err := s.analysisRepo.UpdateFailed(id, code, message); err != nil {
		log.Errorf("Failed to update failed status for %s: %v", id, err)
	}
}

// GetAnalysis returns analysis by ID and session
func (s *AnalyzeService) GetAnalysis(analysisID, sessionID uuid.UUID) (*model.Analysis, error) {
	return s.analysisRepo.FindByIDAndSession(analysisID, sessionID)
}
