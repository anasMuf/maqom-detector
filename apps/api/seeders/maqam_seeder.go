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
			ExampleSongs:        []string{"Adzan — banyak muadzin menggunakan Hijaz", "Sholawat bertema kerinduan dan cinta Rasulullah", "Lagu pembukaan yang dramatis dan khidmat"},
			TipsAransemen:       "Maqam Hijaz sangat cocok untuk bagian pembuka yang dramatis. Pada vokal banjari, perhatikan ornamentasi pada nada E♭ dan interval augmented second (E♭ ke F♯) yang menjadi ciri khas.",
		},
		{
			ID: "rast", NameLatin: "Rast", NameArabic: "راست", NameIndonesia: "Rast",
			IntervalDescription: "C – D – E♭↑ – F – G – A – B♭↑",
			EmotionTags:         []string{"tenang", "gembira", "natural", "megah"},
			ExampleSongs:        []string{"Pembukaan acara maulid dan perayaan", "Tilawah Al-Quran — salah satu maqam paling fundamental", "Lagu pujian yang megah dan berwibawa"},
			TipsAransemen:       "Maqam Rast adalah maqam paling fundamental dan sering digunakan untuk pembuka acara. Nuansa ceria dan megah cocok untuk lagu pujian.",
		},
		{
			ID: "bayati", NameLatin: "Bayati", NameArabic: "بياتي", NameIndonesia: "Bayati",
			IntervalDescription: "D – E♭↑ – F – G – A – B♭ – C",
			EmotionTags:         []string{"melankolis", "khusyuk", "mendalam", "reflektif"},
			ExampleSongs:        []string{"Tilawah Al-Quran — maqam paling sering digunakan qari", "Sholawat bernuansa mendalam dan reflektif", "Doa dan munajat yang penuh kekhusyukan"},
			TipsAransemen:       "Maqam Bayati adalah salah satu yang paling sering digunakan dalam musik Arab. Cocok untuk lagu-lagu yang mendalam dan penuh perasaan.",
		},
		{
			ID: "nahawand", NameLatin: "Nahawand", NameArabic: "نهاوند", NameIndonesia: "Nahawand",
			IntervalDescription: "C – D – E♭ – F – G – A♭ – B♭",
			EmotionTags:         []string{"sedih", "romantis", "lembut", "penuh perasaan"},
			ExampleSongs:        []string{"Sholawat bertema romantis dan lembut", "Lagu doa yang penuh perasaan dan kerinduan", "Nasheed modern bergaya minor — familiar bagi pendengar musik Barat"},
			TipsAransemen:       "Nahawand mirip dengan tangga nada minor natural Barat sehingga terasa familiar. Cocok untuk lagu-lagu bertema cinta, kerinduan, atau doa yang lembut.",
		},
		{
			ID: "kurd", NameLatin: "Kurd", NameArabic: "كرد", NameIndonesia: "Kurd",
			IntervalDescription: "D – E♭ – F – G – A – B♭ – C",
			EmotionTags:         []string{"sedih", "syahdu", "haru", "tenang"},
			ExampleSongs:        []string{"Sholawat syahdu yang tenang dan mengalir", "Tilawah Al-Quran bernuansa menenangkan jiwa", "Lagu bertema patriotik dan nasionalisme Arab"},
			TipsAransemen:       "Kurd mirip Bayati tetapi tanpa quarter-tone. Lebih mudah dinyanyikan oleh vokalis yang belum terbiasa dengan quarter-tone.",
		},
		{
			ID: "saba", NameLatin: "Saba", NameArabic: "صبا", NameIndonesia: "Saba",
			IntervalDescription: "D – E♭↑ – F – G♭ – A – B♭ – C",
			EmotionTags:         []string{"sangat sedih", "ratapan", "duka", "pilu"},
			ExampleSongs:        []string{"Bagian klimaks yang mengharukan dalam pertunjukan", "Doa taubat dan permohonan yang penuh emosi", "Ratapan dan ungkapan duka yang mendalam"},
			TipsAransemen:       "Maqam Saba dikenal sebagai maqam paling emosional dan sedih. Ciri khasnya adalah interval diminished fourth (G♭). Pola rebana sebaiknya minimal agar vokal mendominasi.",
		},
		{
			ID: "ajam", NameLatin: "Ajam", NameArabic: "عجم", NameIndonesia: "Ajam",
			IntervalDescription: "C – D – E – F – G – A – B",
			EmotionTags:         []string{"ceria", "optimis", "megah", "bersemangat"},
			ExampleSongs:        []string{"Lagu perayaan, syukur, dan kegembiraan", "Pembukaan atau penutup acara yang enerjik dan megah", "Nasheed ceria bertema optimisme dan semangat"},
			TipsAransemen:       "Ajam identik dengan tangga nada mayor Barat. Nuansanya ceria dan bersemangat — cocok untuk lagu pembuka yang enerjik atau penutup yang megah.",
		},
		{
			ID: "jiharkah", NameLatin: "Jiharkah", NameArabic: "جهاركاه", NameIndonesia: "Jiharkah",
			IntervalDescription: "F – G – A – B♭ – C – D – E♭↑",
			EmotionTags:         []string{"gembira", "cerah", "merayakan", "ringan"},
			ExampleSongs:        []string{"Lagu maulid yang riang dan merayakan", "Qasidah bertema pujian yang cerah dan ringan", "Pembacaan Al-Quran bernuansa gembira dan ekspresif"},
			TipsAransemen:       "Jiharkah memiliki nuansa ceria dan merayakan. Sering digunakan dalam lagu-lagu maulid yang riang.",
		},
	}
}
