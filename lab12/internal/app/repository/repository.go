package repository

import (
	"errors"
	"fmt"
	"strings"
)

// Feature — структура для одной буквы/признака
type Feature struct {
	ID          int
	Name        string
	Description string
	Period      string
	Details     string
	ImageURL    string
	Count       int
}

// Manuscript — структура для рукописи
type Manuscript struct {
	ID    int
	Signs []Feature
}

// Repository — хранилище данных
type Repository struct {
	Manuscripts []Manuscript
}

// ====== "База данных" признаков ======
func getSeedFeatures() ([]Feature, error) {
	features := []Feature{
		{
			ID:          1,
			Name:        "Буква ѣ (ять)",
			Description: "Активно используется с XI по XVIII вв.",
			Period:      "1000–1750",
			Details:     "Буква пришла из старославянской письменности. В живой речи к XVII веку различие между ѣ и е исчезло, но в письменности буква сохранялась до реформы 1918 года. После реформы её полностью заменили на е.",
			ImageURL:    "http://127.0.0.1:9000/manuscripts/yat.jpg",
		},
		{
			ID:          2,
			Name:        "Буква ѵ (ижица)",
			Description: "Окончательно исчезает в XVIII веке.",
			Period:      "1000–1700",
			Details:     "Ижица использовалась для передачи греческой буквы 'υ'. Уже в XVII веке из употребления почти исчезла, встречалась только в церковных книгах. В 1735 году ижица была исключена из гражданского шрифта указом Академии наук.",
			ImageURL:    "http://127.0.0.1:9000/manuscripts/izhitsa.jpg",
		},
		{
			ID:          3,
			Name:        "Буква Ѳ (фита)",
			Description: "Исчезает в начале XX века после реформы.",
			Period:      "1000–1918",
			Details:     "Фита использовалась в заимствованных греческих словах для передачи звука 'ф'. В русском языке звука [θ] никогда не существовало, поэтому фита звучала так же, как обычная ф. В результате она стала избыточной и была исключена в 1918 году, заменившись на ф.",
			ImageURL:    "http://127.0.0.1:9000/manuscripts/fita.jpg",
		},
		{
			ID:          4,
			Name:        "Твёрдый знак на конце слов",
			Description: "Обязателен до реформы 1918 года.",
			Period:      "1000–1918",
			Details:     "До реформы 1918 года твёрдый знак ставился на конце каждого слова после согласной. Эта норма сильно утяжеляла тексты: считается, что около 4–5% печатной площади в книгах занимали ненужные твёрдые знаки. В 1918 году обязательное написание на конце слов отменили.",
			ImageURL:    "http://127.0.0.1:9000/manuscripts/hardsing.jpg",
		},
		{
			ID:          5,
			Name:        "Буква i",
			Description: "Используется до реформы 1918 года.",
			Period:      "1000–1918",
			Details:     "Эта буква была введена для более точного соответствия греческому письму и для различения омонимов. Но к XIX веку различие стало формальным и мешало обучению. В 1918 году её отменили, и все случаи написания заменили на и.",
			ImageURL:    "http://127.0.0.1:9000/manuscripts/decimalI.jpg",
		},
	}

	if len(features) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return features, nil
}

// ====== Конструктор ======
func NewRepository() (*Repository, error) {
	features, err := getSeedFeatures()
	if err != nil {
		return nil, err
	}

	// создаём стартовую рукопись с ID = 1
	initial := Manuscript{
		ID: 1,
		Signs: []Feature{
			features[0], // добавляем первый признак
			features[1], // добавляем второй признак
		},
	}
	initial.Signs[0].Count = 2
	initial.Signs[1].Count = 1

	return &Repository{
		Manuscripts: []Manuscript{initial},
	}, nil
}

// ====== Методы работы ======

// Получить все признаки
func (r *Repository) GetFeatures() ([]Feature, error) {
	return getSeedFeatures()
}

// Получить один признак
func (r *Repository) GetFeature(id int) (Feature, error) {
	features, err := getSeedFeatures()
	if err != nil {
		return Feature{}, err
	}
	for _, f := range features {
		if f.ID == id {
			return f, nil
		}
	}
	return Feature{}, errors.New("признак не найден")
}

// Поиск по признакам
func (r *Repository) SearchFeatures(query string) ([]Feature, error) {
	features, err := getSeedFeatures()
	if err != nil {
		return []Feature{}, err
	}
	var result []Feature
	queryLower := strings.ToLower(query)
	for _, f := range features {
		if strings.Contains(strings.ToLower(f.Name), queryLower) {
			result = append(result, f)
		}
	}
	return result, nil
}

// Получить рукопись по ID
func (r *Repository) GetManuscript(id int) (Manuscript, error) {
	for _, m := range r.Manuscripts {
		if m.ID == id {
			return m, nil
		}
	}
	return Manuscript{}, fmt.Errorf("рукопись с id %d не найдена", id)
}

// Получить общее количество признаков в рукописи
func (r *Repository) GetTotalManuscriptCount(id int) int {
	for _, m := range r.Manuscripts {
		if m.ID == id {
			total := 0
			for _, f := range m.Signs {
				total += f.Count
			}
			return total
		}
	}
	return 0
}
