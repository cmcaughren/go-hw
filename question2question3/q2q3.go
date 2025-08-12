/*
Write a function in Go, Java, or C# that returns the number of occurrences of a specified digit in a series
between 2 numbers that have been supplied, and whose behaviour is controlled by a supplied
parameter. Ideally the function signature should have the following parameters:
• SeriesStart – the number that begins the series, which can be a number between -(253-1) and
(253-1)
• SeriesEnd – the number that ends the series, which can be a number between -(253-1) and (253-
1)
• SeriesIncrement – the increment to be used to determine the individual elements in the series
• SpecifiedDigit – the digit that you want the number of occurrences of in the series
• SeriesType – an identifier that affects which items in the series to analyze, it can have 1 of the
following values:
o 1 – analyze all elements in the series
o 2 – analyze only even numbered elements in the series
o 3 – analyze only odd numbered elements in the series
As an example, if I called this function with the following parameters:
• SeriesStart = 1
• SeriesEnd = 11
• SeriesIncrement = 1
• SpecifiedDigit = 1
• SeriesType = 1
I would expect that the function would return 4 for the number of occurrences of the digit “1”, as it
would be present in the series elements 1, 10, and 11 and it occurs 4 times, once in each of “1” and
“10”, and twice in “11”.

*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SeriesInfo struct {
	SeriesStart     int64
	SeriesEnd       int64
	SeriesIncrement int64
	SpecifiedDigit  int
	SeriesType      int
}

func main() {

	// Set up HTTP routes
	r := mux.NewRouter()
	r.HandleFunc("/{SeriesStart}/{SeriesEnd}/{SeriesIncrement}/{SpecifiedDigit}/{SeriesType}", Series).Methods("GET")

	//Start Server
	port := ":8080"
	fmt.Printf("Starting server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))

}

func Series(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, err := parseSeriesInfo(vars)
	if err != nil {
		fmt.Printf("%v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sum := 0
	string_sd := rune(s.SpecifiedDigit + '0')

	for i := s.SeriesStart; i <= s.SeriesEnd; i += s.SeriesIncrement {
		if (s.SeriesType == 1) ||
			((s.SeriesType == 2) && (i%2 == 0)) ||
			((s.SeriesType == 3) && (i%2 != 0)) {
			string_i := strconv.FormatInt(i, 10)
			for _, r := range string_i {
				if r == string_sd {
					sum++
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sum)
}

func parseSeriesInfo(vars map[string]string) (SeriesInfo, error) {
	var s SeriesInfo
	var err error

	s.SeriesStart, err = strconv.ParseInt(vars["SeriesStart"], 10, 64)
	if err != nil {
		return s, fmt.Errorf("SeriesStart requires a number between -(2^53 - 1) and (2^53 - 1)")
	}

	s.SeriesEnd, err = strconv.ParseInt(vars["SeriesEnd"], 10, 64)
	if err != nil {
		return s, fmt.Errorf("SeriesEnd requires a number between -(2^53 - 1) and (2^53 - 1)")
	}

	s.SeriesIncrement, err = strconv.ParseInt(vars["SeriesIncrement"], 10, 64)
	if err != nil {
		return s, fmt.Errorf("SeriesIncrement requires a number between -(2^53 - 1) and (2^53 - 1)")
	}

	s.SpecifiedDigit, err = strconv.Atoi(vars["SpecifiedDigit"])
	if err != nil || s.SpecifiedDigit > 9 || s.SpecifiedDigit < 0 {
		return s, fmt.Errorf("SeriesDigit must be an integer 0 - 9")
	}

	s.SeriesType, err = strconv.Atoi(vars["SeriesType"])
	if err != nil || s.SeriesType > 3 || s.SeriesType < 1 {
		errMsg := `SeriesType must be one of: 
			1 – analyze all elements in the series
			2 – analyze only even numbered elements in the series
			3 – analyze only odd numbered elements in the series
		`
		return s, fmt.Errorf("%s", errMsg)
	}
	return s, nil
}
