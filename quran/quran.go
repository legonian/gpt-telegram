package quran

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Quran struct {
	surahs    map[string]Surah
	positions map[string]Position
	words     map[string]string
	ayahs     []string
}

const (
	FileSurah     = "./surah.json"
	FilePositions = "./positions.json"
	FileWord      = "./word.json"
)

type Surah struct {
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

type Position struct {
	Surah    int64 `json:"surah"`
	Ayah     int64 `json:"ayah"`
	Position int64 `json:"position"`
}

type WordPosition struct {
	Word     string
	Position int64
}

func NewQuran() (*Quran, error) {
	surah, err := getSurah()
	if err != nil {
		return nil, fmt.Errorf("getSurah: %w", err)
	}
	positions, err := getPositions()
	if err != nil {
		return nil, fmt.Errorf("getPosition: %w", err)
	}
	wordValues, err := getWord()
	if err != nil {
		return nil, fmt.Errorf("getWord: %w", err)
	}
	if len(positions) != len(wordValues) {
		return nil, fmt.Errorf("positions not matched with words number")
	}

	q := &Quran{
		surahs:    surah,
		positions: positions,
		words:     wordValues,
	}

	ayahs := make([]string, 0)
	buff := make([]string, 0)
	wordsNumber := len(q.positions)
	var lastS, lastA int64
	lastS, lastA = 1, 1
	for i := 0; i < wordsNumber; i++ {
		wordID := strconv.Itoa(i)
		p := q.positions[wordID]
		if p.Surah != lastS || p.Ayah != lastA {
			lastS, lastA = p.Surah, p.Ayah
			ayahs = append(ayahs, strings.Join(buff, " "))
			buff = make([]string, 0)
		}

		buff = append(buff, q.words[wordID])
	}
	q.ayahs = ayahs
	return q, nil
}

func getSurah() (map[string]Surah, error) {
	f, err := os.Open(FileSurah)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	var s map[string]Surah
	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder(f).Decode: %w", err)
	}
	err = f.Close()
	if err != nil {
		return nil, fmt.Errorf("f.Close: %w", err)
	}

	return s, nil
}

func getPositions() (map[string]Position, error) {
	f, err := os.Open(FilePositions)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	var p map[string]Position
	err = json.NewDecoder(f).Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder(f).Decode: %w", err)
	}
	err = f.Close()
	if err != nil {
		return nil, fmt.Errorf("f.Close: %w", err)
	}

	return p, nil
}

func getWord() (map[string]string, error) {
	f, err := os.Open(FileWord)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	var w map[string]string
	err = json.NewDecoder(f).Decode(&w)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder(f).Decode: %w", err)
	}
	err = f.Close()
	if err != nil {
		return nil, fmt.Errorf("f.Close: %w", err)
	}

	return w, nil
}

func (q *Quran) GetAyah(surah, ayah int64) string {
	quotePositions := []WordPosition{}
	for wordID, p := range q.positions {
		if p.Surah != surah {
			continue
		}
		if p.Ayah != ayah {
			continue
		}

		quotePositions = append(quotePositions, WordPosition{
			Word:     q.words[wordID],
			Position: p.Position,
		})
	}

	quote := make([]string, len(quotePositions))
	for _, qp := range quotePositions {
		quote[qp.Position] = qp.Word
	}

	return strings.Join(quote, " ")
}

func (q *Quran) MaxSurah() int64 {
	return q.positions[strconv.Itoa(len(q.positions))].Surah
}

func (q *Quran) MaxAyah(surah int64) int64 {
	if q.MaxSurah() < surah {
		return -1
	}

	var maxAyah int64
	for _, p := range q.positions {
		if surah != p.Surah {
			continue
		}

		if maxAyah < p.Ayah {
			maxAyah = p.Ayah
		}
	}

	return maxAyah
}

func (q *Quran) RandomAyah(count int) string {
	if count < 1 || len(q.ayahs) < count {
		return ""
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	pos := int64(r1.Intn(len(q.ayahs)-count)) + 1
	quote := make([]string, count)
	for i := 0; i < count; i++ {
		quote[i] = q.ayahs[pos+int64(i)]
	}

	return strings.Join(quote, "\n\n")
}
