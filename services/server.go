package services

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type counter struct {
	count     uint64
	timeStamp time.Time
}

type responseJson struct {
	Count string `json:"count"`
}

func NewRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getSliceOfValuesFromFile(filename string) ([]int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	var value string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value = fmt.Sprintf("%s", scanner.Text())
	}
	splitValues := strings.Split(value, ",")
	if len(splitValues) != 2 {
		return nil, errors.New("invalid file format")
	}

	// convert the first value to int
	counterValue, err := strconv.Atoi(splitValues[0])
	if err != nil {
		return nil, err
	}

	// check the unix format of the timestamp
	timestampIntValue, err := strconv.ParseInt(splitValues[1], 10, 64)
	if err != nil {
		return nil, errors.New("invalid file format")
	}
	return []int64{int64(counterValue), timestampIntValue}, nil
}

// create a function that stores a value into a file
func storeValue(value string) error {
	// open the file
	file, err := os.OpenFile(storeFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	err = file.Truncate(0)
	_, err = file.Seek(0, 0)

	// current time stamp
	now := time.Now().Unix()
	nowString := fmt.Sprintf("%d", now)
	// write some text line-by-line to file
	_, err = file.WriteString(value + "," + nowString + "\n")
	if err != nil {
		return err
	}

	// save changes
	err = file.Sync()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func assertCounterPath(path string) error {
	if path == CounterPath {
		return nil
	} else {
		return errors.New("invalid path")
	}
}

// define a function to add the trailing slash
func addTrailingSlash(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	}
	return path + "/"
}

var counterObj counter

func validateTimestampWithWindow(window float64, currentTime, timestamp int64) bool {
	timeStampInTimeVersion := time.Unix(timestamp, 0)
	currentTimeTimeVersion := time.Unix(currentTime, 0)

	diff := currentTimeTimeVersion.Sub(timeStampInTimeVersion)

	if diff.Seconds() > 0 && diff.Seconds() < window {
		return true
	}
	return false
}

func counterObjectValidator(c *counter, currentTime time.Time) {
	if validateTimestampWithWindow(window, currentTime.Unix(), c.timeStamp.Unix()) {
		counterObj.count = c.count
		counterObj.timeStamp = time.Unix(c.timeStamp.Unix(), 0)
	} else {
		counterObj.count = uint64(0)
		counterObj.timeStamp = time.Now()
	}
}

func ApiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()

		addTrailingSlash(r.URL.Path)

		if err := assertCounterPath(r.URL.Path); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// if it is the first time or the server was shutdown
		if counterObj.count == 0 {
			// check if the value exists in the file
			valuesFromStoredFile, err := getSliceOfValuesFromFile(storeFile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			counterObjectValidator(&counter{
				count:     uint64(valuesFromStoredFile[0]),
				timeStamp: time.Unix(valuesFromStoredFile[1], 0),
			}, currentTime)

		} else {
			counterObjectValidator(&counterObj, currentTime)
		}

		// add 1 to the counter
		atomic.AddUint64(&counterObj.count, 1)
		// load the counter value
		loadedValue := atomic.LoadUint64(&counterObj.count)

		// store the counter value to the file
		go func() {
			err := storeValue(fmt.Sprintf("%d", loadedValue))
			if err != nil {
				fmt.Println(err)
			}
		}()

		// write the response headers
		w.Header().Set(ContentType, ContentTypeValue)
		w.Header().Set(XContentTypeOptions, NoSniff)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseJson{Count: fmt.Sprintf("%d", loadedValue)})

	})
}

func Logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestNumber).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Tracing is a middleware that traces the request
func Tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(XRequestId)
			if requestID == Empty {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestNumber, requestID)
			w.Header().Set(XRequestId, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
