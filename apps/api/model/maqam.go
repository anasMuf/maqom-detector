package model

// Maqam — 8 maqam core yang didukung sistem
type Maqam struct {
	ID                 string  `json:"id" gorm:"type:varchar(50);primaryKey"`
	NameLatin          string  `json:"name_latin" gorm:"type:varchar(100);not null"`
	NameArabic         string  `json:"name_arabic" gorm:"type:varchar(100);not null"`
	NameIndonesia      string  `json:"name_indonesia" gorm:"type:varchar(100)"`
	IntervalDescription string `json:"interval_description" gorm:"type:varchar(255)"`
	EmotionTags        string  `json:"emotion_tags" gorm:"type:text"`           // JSON array stored as text
	ExampleSongs       string  `json:"example_songs" gorm:"type:text"`          // JSON array stored as text
	TipsAransemen      string  `json:"tips_aransemen" gorm:"type:text"`
}

func (Maqam) TableName() string {
	return "maqamat"
}
