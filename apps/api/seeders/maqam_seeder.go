package seeders

import (
	"api/model"
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

// MaqamSeedData represents the JSON structure from maqam_templates.json
type MaqamSeedData struct {
	ID                  string   `json:"id"`
	NameLatin           string   `json:"name_latin"`
	NameArabic          string   `json:"name_arabic"`
	NameIndonesia       string   `json:"name_indonesia"`
	IntervalDescription string   `json:"interval_description"`
	EmotionTags         []string `json:"emotion_tags"`
	ExampleSongs        []string `json:"example_songs"`
	TipsAransemen       string   `json:"tips_aransemen_banjari"`
}

// SeedMaqamat seeds the maqamat table with 8 core maqams.
// Only runs if the table is empty.
func SeedMaqamat(db *gorm.DB) {
	var count int64
	db.Model(&model.Maqam{}).Count(&count)
	if count > 0 {
		log.Println("Maqamat sudah ada di database, skip seeding")
		return
	}

	maqams := getMaqamData()
	for _, m := range maqams {
		emotionJSON, _ := json.Marshal(m.EmotionTags)
		songsJSON, _ := json.Marshal(m.ExampleSongs)

		maqam := model.Maqam{
			ID:                  m.ID,
			NameLatin:           m.NameLatin,
			NameArabic:          m.NameArabic,
			NameIndonesia:       m.NameIndonesia,
			IntervalDescription: m.IntervalDescription,
			EmotionTags:         string(emotionJSON),
			ExampleSongs:        string(songsJSON),
			TipsAransemen:       m.TipsAransemen,
		}

		if err := db.Create(&maqam).Error; err != nil {
			log.Printf("Gagal seed maqam %s: %v\n", m.ID, err)
		}
	}

	log.Printf("Berhasil seed %d maqam ke database\n", len(maqams))
}

func getMaqamData() []MaqamSeedData {
	return []MaqamSeedData{
		{
			ID: "hijaz", NameLatin: "Hijaz", NameArabic: "حجاز", NameIndonesia: "Hijaz",
			IntervalDescription: "D – E♭ – F♯ – G – A – B♭ – C",
			EmotionTags:         []string{"dramatis", "kerinduan", "agung", "spiritual"},
			ExampleSongs:        []string{"Ya Hanana", "Qasidah Burda (pembuka)", "Tala'al Badru Alayna (versi Hijaz)"},
			TipsAransemen:       "Maqam Hijaz sangat cocok untuk bagian pembuka yang dramatis. Pada vokal banjari, perhatikan ornamentasi pada nada E♭ dan interval augmented second (E♭ ke F♯) yang menjadi ciri khas.",
		},
		{
			ID: "rast", NameLatin: "Rast", NameArabic: "راست", NameIndonesia: "Rast",
			IntervalDescription: "C – D – E♭↑ – F – G – A – B♭↑",
			EmotionTags:         []string{"tenang", "gembira", "natural", "megah"},
			ExampleSongs:        []string{"Mawtini", "Ya Taiba", "Assalamualaika Ya Rasulallah"},
			TipsAransemen:       "Maqam Rast adalah maqam paling fundamental dan sering digunakan untuk pembuka acara. Nuansa ceria dan megah cocok untuk lagu pujian.",
		},
		{
			ID: "bayati", NameLatin: "Bayati", NameArabic: "بياتي", NameIndonesia: "Bayati",
			IntervalDescription: "D – E♭↑ – F – G – A – B♭ – C",
			EmotionTags:         []string{"melankolis", "khusyuk", "mendalam", "reflektif"},
			ExampleSongs:        []string{"Deen Assalam", "Law Kana Bainana", "Ya Nabi Salam Alayka"},
			TipsAransemen:       "Maqam Bayati adalah salah satu yang paling sering digunakan dalam musik Arab. Cocok untuk lagu-lagu yang mendalam dan penuh perasaan.",
		},
		{
			ID: "nahawand", NameLatin: "Nahawand", NameArabic: "نهاوند", NameIndonesia: "Nahawand",
			IntervalDescription: "C – D – E♭ – F – G – A♭ – B♭",
			EmotionTags:         []string{"sedih", "romantis", "lembut", "penuh perasaan"},
			ExampleSongs:        []string{"Kun Anta", "Maher Zain - Insha Allah", "Sholawat Nariyah"},
			TipsAransemen:       "Nahawand mirip dengan tangga nada minor natural Barat sehingga terasa familiar. Cocok untuk lagu-lagu bertema cinta, kerinduan, atau doa yang lembut.",
		},
		{
			ID: "kurd", NameLatin: "Kurd", NameArabic: "كرد", NameIndonesia: "Kurd",
			IntervalDescription: "D – E♭ – F – G – A – B♭ – C",
			EmotionTags:         []string{"sedih", "syahdu", "haru", "tenang"},
			ExampleSongs:        []string{"Isyfa Lana", "Ya Rasulallah (versi Kurd)", "Antal Hayat"},
			TipsAransemen:       "Kurd mirip Bayati tetapi tanpa quarter-tone. Lebih mudah dinyanyikan oleh vokalis yang belum terbiasa dengan quarter-tone.",
		},
		{
			ID: "saba", NameLatin: "Saba", NameArabic: "صبا", NameIndonesia: "Saba",
			IntervalDescription: "D – E♭↑ – F – G♭ – A – B♭ – C",
			EmotionTags:         []string{"sangat sedih", "ratapan", "duka", "pilu"},
			ExampleSongs:        []string{"Adfaita (bagian sedih)", "Ya Arhamar Rahimin", "Tawasul - Abu Ratib"},
			TipsAransemen:       "Maqam Saba dikenal sebagai maqam paling emosional dan sedih. Ciri khasnya adalah interval diminished fourth (G♭). Pola rebana sebaiknya minimal agar vokal mendominasi.",
		},
		{
			ID: "ajam", NameLatin: "Ajam", NameArabic: "عجم", NameIndonesia: "Ajam",
			IntervalDescription: "C – D – E – F – G – A – B",
			EmotionTags:         []string{"ceria", "optimis", "megah", "bersemangat"},
			ExampleSongs:        []string{"Hasbi Rabbi", "Ahmad Ya Habibi", "Ya Badratim"},
			TipsAransemen:       "Ajam identik dengan tangga nada mayor Barat. Nuansanya ceria dan bersemangat — cocok untuk lagu pembuka yang enerjik atau penutup yang megah.",
		},
		{
			ID: "jiharkah", NameLatin: "Jiharkah", NameArabic: "جهاركاه", NameIndonesia: "Jiharkah",
			IntervalDescription: "F – G – A – B♭ – C – D – E♭↑",
			EmotionTags:         []string{"gembira", "cerah", "merayakan", "ringan"},
			ExampleSongs:        []string{"Sidnan Nabi", "Ya Hanana (versi ceria)", "Busyra Lana"},
			TipsAransemen:       "Jiharkah memiliki nuansa ceria dan merayakan. Sering digunakan dalam lagu-lagu maulid yang riang.",
		},
	}
}
