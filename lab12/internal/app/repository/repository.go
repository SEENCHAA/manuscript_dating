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

type Repository struct {
	Manuscript []Feature // вместо Order теперь "Рукопись"
}

// Список признаков (буквы, особенности орфографии)
var Features = []Feature{
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

// Создание нового репозитория с предустановленной рукописью
func NewRepository() (*Repository, error) {
	initialManuscript := []Feature{}
	if len(Features) >= 2 {
		// Первый признак с количеством 2
		f1 := Features[0]
		f1.Count = 2
		initialManuscript = append(initialManuscript, f1)

		// Второй признак с количеством 1
		f2 := Features[1]
		f2.Count = 1
		initialManuscript = append(initialManuscript, f2)
	}

	return &Repository{
		Manuscript: initialManuscript,
	}, nil
}

// Получить все признаки
func (r *Repository) GetFeatures() ([]Feature, error) {
	if len(Features) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}
	return Features, nil
}

// Получить один признак
func (r *Repository) GetFeature(id int) (Feature, error) {
	features, err := r.GetFeatures()
	if err != nil {
		return Feature{}, err
	}
	for _, f := range features {
		if f.ID == id {
			return f, nil
		}
	}
	return Feature{}, errors.New("не найдено")
}

// Поиск по признакам
func (r *Repository) SearchFeatures(query string) ([]Feature, error) {
	features, err := r.GetFeatures()
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

// Добавить признак в рукопись
func (r *Repository) AddToManuscript(id int) error {
	feature, err := r.GetFeature(id)
	if err != nil {
		return err
	}

	for i := range r.Manuscript {
		if r.Manuscript[i].ID == id {
			r.Manuscript[i].Count++ // увеличиваем количество
			return nil
		}
	}

	feature.Count = 1 // первый раз
	r.Manuscript = append(r.Manuscript, feature)
	return nil
}

// Обновить количество признака в рукописи
func (r *Repository) UpdateManuscriptCount(id int, count int) {
	for i := range r.Manuscript {
		if r.Manuscript[i].ID == id {
			r.Manuscript[i].Count = count
			return
		}
	}
}

// Получить всю рукопись
func (r *Repository) GetManuscript() []Feature {
	return r.Manuscript
}

// Получить общее количество признаков в рукописи
func (r *Repository) GetTotalManuscriptCount() int {
	total := 0
	for _, f := range r.Manuscript {
		total += f.Count
	}
	return total
}
