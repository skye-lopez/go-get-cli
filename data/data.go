package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type Store struct {
	GoIndex       map[string]GoIndexEntry
	LastWriteTime string
}

type GoIndexEntry struct {
	Path      string `json:"Path"`
	Version   string `json:"Version"`
	Timestamp string `json:"Timestamp"`
}

func Init(db *leveldb.DB) {
	// Collect data
	ParseGoIndex(db)

	// Update lastWriteTime
	newTime := time.Now().Format(time.RFC3339Nano)
	db.Put([]byte("lastWriteTime"), []byte(newTime), nil)

	fmt.Println("Done fetching index data! Enjoy.")
}

func ParseGoIndex(db *leveldb.DB) {
	startTime := time.Now()
	var endTimeString string
	storedEndTime, err := db.Get([]byte("lastWriteTime"), nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		fmt.Println("Fetching index data... Its your first time so this could take awhile, please be patient! :)")
		endTimeString = "2019-04-10T19:08:52.997264Z"
	} else {
		endTimeString = string(storedEndTime)
		fmt.Println("Fetching new index data...")
	}

	endTime, err := time.Parse(time.RFC3339Nano, endTimeString)
	if err != nil {
		panic(err)
	}

	// NOTE: 12hrs is a relatively arbitrary time... could be less or more need to do some testing on that.
	step := time.Duration(12) * time.Hour

	urls := []string{}
	for startTime.Unix() > endTime.Unix() {
		url := "https://index.golang.org/index?since=" + startTime.Format(time.RFC3339Nano)
		urls = append(urls, url)

		startTime = startTime.Add(-step)
	}

	maxChanSize := len(urls) * 2000
	entries := make(chan GoIndexEntry, maxChanSize)
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Interesting race condition with TCP io.ReadAll...
			time.Sleep(time.Millisecond * 50)

			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}

			rawEntries := strings.Split(string(body), "\n")
			for _, re := range rawEntries {
				e := GoIndexEntry{}

				json.Unmarshal([]byte(re), &e)
				entries <- e
			}
		}()
	}

	wg.Wait()
	close(entries)

	for e := range entries {
		_, err := db.Get([]byte(e.Path), nil)
		if errors.Is(err, leveldb.ErrNotFound) {
			// NOTE: For now we arent going to store the actual version of the package, this would be an easy change with some kind of delimiter
			db.Put([]byte(e.Path), []byte(e.Path), nil)
		}
	}
}
