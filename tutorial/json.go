package tutorial

import (
	"encoding/json"
	"fmt"
)

type Movie struct {
	Title     string `json:"title"`
	Content   string `json:"content, omitempty"`
	Year      int    `json:"year"`
	ReadCount int    `json:"-"`
}

func jsonMain() {
	myMovie := Movie{Title: "수상한 그녀", Content: "코미디", Year: 2014, ReadCount: 1}

	jsonData, _ := json.Marshal(myMovie)

	fmt.Println(string(jsonData))

	jsonString := `{"title": "라디오 스타", "year": 2005}`
	var movieObject Movie
	json.Unmarshal([]byte(jsonString), &movieObject)
	fmt.Printf("제목: %s, 개봉: %d\n", movieObject.Title, movieObject.Year)
}
