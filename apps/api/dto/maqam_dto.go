package dto

// ── Analyze Requests ─────────────────────────

// AnalyzeYoutubeRequest represents the request body for YouTube analysis
type AnalyzeYoutubeRequest struct {
	URL             string `json:"url" validate:"required,url" example:"https://www.youtube.com/watch?v=dQw4w9WgXcQ"`
	SegmentStart    int    `json:"segment_start" validate:"gte=0" example:"0"`
	SegmentDuration int    `json:"segment_duration" validate:"gte=5,lte=120" example:"60"`
}

// ── Analyze Response (202) ───────────────────

// AnalyzeAcceptedResponse returned when analysis is queued
type AnalyzeAcceptedResponse struct {
	AnalysisID       string `json:"analysis_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status           string `json:"status" example:"pending"`
	EstimatedSeconds int    `json:"estimated_seconds" example:"30"`
	PollingURL       string `json:"polling_url" example:"/api/v1/analyses/550e8400-e29b-41d4-a716-446655440000"`
}

// ── Analysis Detail (polling result) ─────────

// CandidateResponse represents a maqam candidate in the analysis result
type CandidateResponse struct {
	MaqamID         string  `json:"maqam_id" example:"hijaz"`
	NameLatin       string  `json:"name_latin" example:"Hijaz"`
	NameArabic      string  `json:"name_arabic" example:"حجاز"`
	ConfidenceScore float64 `json:"confidence_score" example:"0.87"`
	Rank            int     `json:"rank" example:"1"`
}

// AnalysisDetailResponse returned when polling for analysis status/result
type AnalysisDetailResponse struct {
	ID              string              `json:"id"`
	Status          string              `json:"status" example:"completed"`
	InputType       string              `json:"input_type" example:"youtube"`
	InputSource     string              `json:"input_source" example:"https://www.youtube.com/watch?v=..."`
	DetectedMaqamID *string             `json:"detected_maqam_id,omitempty" example:"hijaz"`
	ConfidenceScore *float64            `json:"confidence_score,omitempty" example:"0.87"`
	ConfidenceLabel string              `json:"confidence_label,omitempty" example:"tinggi"`
	ExplanationText string              `json:"explanation_text,omitempty"`
	AudioQuality    string              `json:"audio_quality,omitempty" example:"good"`
	ProcessingMs    *int                `json:"processing_ms,omitempty" example:"12500"`
	Candidates      []CandidateResponse `json:"candidates,omitempty"`
	ErrorCode       string              `json:"error_code,omitempty"`
	ErrorMessage    string              `json:"error_message,omitempty"`
	CreatedAt       string              `json:"created_at"`
	CompletedAt     *string             `json:"completed_at,omitempty"`
}

// ── History ──────────────────────────────────

// HistoryListResponse paginated history list
type HistoryListResponse struct {
	Items []HistoryItemResponse `json:"items"`
	Meta  PaginationMeta        `json:"meta"`
}

// HistoryItemResponse single item in history list
type HistoryItemResponse struct {
	ID              string   `json:"id"`
	InputType       string   `json:"input_type" example:"youtube"`
	InputSource     string   `json:"input_source"`
	DetectedMaqamID *string  `json:"detected_maqam_id,omitempty"`
	MaqamNameLatin  string   `json:"maqam_name_latin,omitempty"`
	ConfidenceScore *float64 `json:"confidence_score,omitempty"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"created_at"`
}

// PaginationMeta pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	TotalItems int64 `json:"total_items" example:"42"`
	TotalPages int   `json:"total_pages" example:"5"`
}

// ── Maqam ────────────────────────────────────

// MaqamListResponse list of all supported maqams
type MaqamListResponse struct {
	ID                  string   `json:"id" example:"hijaz"`
	NameLatin           string   `json:"name_latin" example:"Hijaz"`
	NameArabic          string   `json:"name_arabic" example:"حجاز"`
	NameIndonesia       string   `json:"name_indonesia" example:"Hijaz"`
	IntervalDescription string   `json:"interval_description" example:"D – E♭ – F♯ – G – A – B♭ – C"`
	EmotionTags         []string `json:"emotion_tags"`
}

// MaqamDetailResponse full detail of a single maqam
type MaqamDetailResponse struct {
	ID                  string   `json:"id"`
	NameLatin           string   `json:"name_latin"`
	NameArabic          string   `json:"name_arabic"`
	NameIndonesia       string   `json:"name_indonesia"`
	IntervalDescription string   `json:"interval_description"`
	EmotionTags         []string `json:"emotion_tags"`
	ExampleSongs        []string `json:"example_songs"`
	TipsAransemen       string   `json:"tips_aransemen"`
}
