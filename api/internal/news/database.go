package news

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
)

type News struct {
	Id      int       `json:"id"`
	Date    time.Time `json:"time"`
	Tags    []string  `json:"tags"`
	Message string    `json:"message"`
}

var JsonNews []News

func init() {
	var err error
	JsonNews, err = FetchJsonNews()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load news.json")
	}
}

func FetchJsonNews() ([]News, error) {
	// get path to file of this function, news json is in the same dir
	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("cannot get caller info")
	}
	dir := filepath.Dir(callerFile)
	filePath := filepath.Join(dir, "news.json")

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var newsList []News
	if err := json.NewDecoder(file).Decode(&newsList); err != nil {
		return nil, err
	}

	return newsList, nil
}

func GetNewsById(id int) (*News, bool) {
	for i := range JsonNews {
		if JsonNews[i].Id == id {
			return &JsonNews[i], true
		}
	}
	return nil, false
}
