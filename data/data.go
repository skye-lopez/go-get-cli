/*
* @authors: Aria Lopez github(skye-lopez) email(aria.lopez.dev@gmail.com)
*
* <Data>
*   Data manager for both the global index and the context based awesome-go packages.
*   Global index source: << https://index.golang.org/ >>
*   Awesome go source: << https://github.com/avelino/awesome-go >>
 */

// TODO: Errors should eventaully bubble up to the UI in a nice way.

package data

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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

func Init() {
	var existing Store

	// If we dont have the file just make a base state
	if _, err := os.Stat("data.json"); errors.Is(err, os.ErrNotExist) {
		baseState := Store{
			GoIndex:       make(map[string]GoIndexEntry, 0),
			LastWriteTime: "2019-04-10T19:08:52.997264Z",
		}
		existing = baseState
	} else {
		// Otherwise read in saved state.
		file, err := os.ReadFile("data.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(file, &existing)
	}

	// Collect data
	ParseGoIndex(existing.GoIndex, existing.LastWriteTime)

	// Store data with new LastWriteTime
	lastWrite := time.Now().Format(time.RFC3339Nano)
	existing.LastWriteTime = lastWrite

	jsonData, err := json.Marshal(existing)
	if err != nil {
		panic(err)
	}
	os.WriteFile("data.json", jsonData, os.ModePerm)
}

func ParseGoIndex(existingMap map[string]GoIndexEntry, lastTime string) {
	startTime := time.Now()
	endTime, err := time.Parse(time.RFC3339Nano, lastTime)
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

	// NOTE: chan size is also variable, would be nice to have a better way of calculating this.
	// Example: len(urls)*2000=max length it could ever be.
	entries := make(chan GoIndexEntry, 1000000)
	var wg sync.WaitGroup
	for i, url := range urls {
		if i < 10 {
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
	}

	wg.Wait()
	close(entries)

	for e := range entries {
		if len(e.Path) <= 1 {
			continue
		}
		_, exists := existingMap[e.Path]
		if !exists {
			existingMap[e.Path] = e
		}
	}
}
