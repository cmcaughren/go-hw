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

	// Set up HTTP route
	r := mux.NewRouter()
	r.HandleFunc("/{SeriesStart}/{SeriesEnd}/{SeriesIncrement}/{SpecifiedDigit}/{SeriesType}", Series).Methods("GET")

	//Start Server
	port := ":8080"
	fmt.Printf("Starting server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))

}

func Series(w http.ResponseWriter, r *http.Request) {
	//variables passed in from http endpoint
	vars := mux.Vars(r)

	//validate the user provided info
	s, err := parseSeriesInfo(vars)
	if err != nil {
		fmt.Printf("%v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := countDigitInSeries(s.SeriesStart, s.SeriesEnd, s.SeriesIncrement, s.SpecifiedDigit, s.SeriesType)

	//format json http response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// convert http string -> int or int64, validate
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

func countDigitInSeries(seriesStart, seriesEnd, seriesIncrement int64, specifiedDigit, seriesType int) int {
	sum := 0

	//convert the int specified digit to a "rune",
	//add '0' to get the ascii number 1, not the control char 1
	string_sd := rune(specifiedDigit + '0')

	for i := seriesStart; i <= seriesEnd; i += seriesIncrement {
		if (seriesType == 1) || //analyze all the numbers
			((seriesType == 2) && (i%2 == 0)) || //analyze only the even elements
			((seriesType == 3) && (i%2 != 0)) { //analyze only the odd elements
			string_i := strconv.FormatInt(i, 10) //convert int64 -> string
			for _, r := range string_i {         //search each "rune" (char) in the string
				if r == string_sd {
					sum++ //increment counter as needed
				}
			}
		}
	}
	return sum
}
