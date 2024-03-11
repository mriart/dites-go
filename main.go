// Dites
// Starts a web server that returns catalan dites
// Marc Riart Solans, 202403

package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	gen    = []string{}
	feb    = []string{}
	mar    = []string{}
	apr    = []string{}
	may    = []string{}
	jun    = []string{}
	jul    = []string{}
	aug    = []string{}
	sep    = []string{}
	oct    = []string{}
	nov    = []string{}
	dec    = []string{}
	none   = []string{}
	todays = []string{}
	all    = []string{}

	noMonth  bool
	noSeason bool
)

// Init initializes all variables, previosly defined as global.
// The function reads line by line the file dites.txt and fills the following slices:
//
//	gen - dec, contains sentences of the month
//	none, the rest, no month associated
//	todays, contains those that are of a particular day (Sant Jordi, Nadal...)
//	all, contains all, for quick searches
func init() {
	// Open the file
	file, err := os.Open("dites.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a regular expression for the particular dites of the day like Sant Jordi, Nadal...
	pattern := `(sant jordi|sant joan|tots sants|nadal|\([0-9])`
	re := regexp.MustCompile(pattern)

	// Create a new scanner
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		noMonth = false
		noSeason = false
		line := scanner.Text()
		lineLower := strings.ToLower(line)

		switch {
		case strings.Contains(lineLower, "gener"):
			gen = append(gen, line)
		case strings.Contains(lineLower, "febrer"):
			feb = append(feb, line)
		case strings.Contains(lineLower, "març"):
			mar = append(mar, line)
		case strings.Contains(lineLower, "abril"):
			apr = append(apr, line)
		case strings.Contains(lineLower, "maig"):
			may = append(may, line)
		case strings.Contains(lineLower, "juny"):
			jun = append(jun, line)
		case strings.Contains(lineLower, "juliol"):
			jul = append(jul, line)
		case strings.Contains(lineLower, "agost"):
			aug = append(aug, line)
		case strings.Contains(lineLower, "setembre"):
			sep = append(sep, line)
		case strings.Contains(lineLower, "octubre"):
			oct = append(oct, line)
		case strings.Contains(lineLower, "novembre"):
			nov = append(nov, line)
		case strings.Contains(lineLower, "desembre"):
			dec = append(dec, line)
		default:
			noMonth = true
		}

		switch {
		case strings.Contains(lineLower, "hivern"):
			dec = append(dec, line)
			gen = append(gen, line)
			feb = append(feb, line)
		case strings.Contains(lineLower, "primavera"):
			mar = append(mar, line)
			apr = append(apr, line)
			may = append(may, line)
		case strings.Contains(lineLower, "estiu"):
			jun = append(jun, line)
			jul = append(jul, line)
			aug = append(aug, line)
		case strings.Contains(lineLower, "tardor"):
			sep = append(sep, line)
			oct = append(oct, line)
			nov = append(nov, line)
		default:
			noSeason = true
		}

		if noMonth && noSeason {
			none = append(none, line)
		}

		if re.MatchString(lineLower) {
			todays = append(todays, line)
		}

		all = append(all, line)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir("./res"))))
	http.HandleFunc("/today", handlerToday)
	http.HandleFunc("/month", handlerMonth)
	http.HandleFunc("/misc", handlerMisc)
	http.HandleFunc("/search", handlerSearch)
	http.HandleFunc("/searchall", handlerSearchAll)

	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handlerToday(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, ditaOfToday())
}

func handlerMonth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, randomDitaMonth())
}

func handlerMisc(w http.ResponseWriter, r *http.Request) {
	idx := rand.IntN(len(none))
	resp := none[idx]

	fmt.Fprint(w, resp)
}

func handlerSearch(w http.ResponseWriter, r *http.Request) {
	// Parse the URL query string
	params := r.URL.Query()

	// Extract the "pattern" parameter. Returns a []string (every parameter can have several comma-separated values)
	pattern, ok := params["pattern"]
	if !ok {
		fmt.Fprintf(w, "Missing 'pattern' parameter in query string")
		return
	}
	fmt.Println(pattern[0])

	// Breaks pattern[0] into all words, and place in a slice
	words := strings.Fields(pattern[0])

	// Loops over slice all to find the words. Fills a that contains the indices of the hits (hitIndices)
	hitIndices := []int{}
	for i, line := range all {
		if containsAllSubstrs(line, words) {
			hitIndices = append(hitIndices, i)
		}
	}

	// Search for a random index, return the line value of []all of that index
	resp := ""
	if len(hitIndices) > 0 {
		idx := rand.IntN(len(hitIndices))
		resp += all[hitIndices[idx]] + "<br>" + "..." + "<br>"
	}

	fmt.Fprint(w, resp)
}

func handlerSearchAll(w http.ResponseWriter, r *http.Request) {
	resp := ""

	// Parse the URL query string
	params := r.URL.Query()

	// Extract the "pattern" parameter. Returns a []string (every parameter can have several comma-separated values)
	pattern, ok := params["pattern"]
	if !ok {
		fmt.Fprintf(w, "Missing 'pattern' parameter in query string")
		return
	}
	fmt.Println(pattern[0])

	// Breaks pattern[0] into all words, and place in a slice
	words := strings.Fields(pattern[0])

	// Looks in slice of strings all for all words
	for _, line := range all {
		if containsAllSubstrs(line, words) {
			resp += line + "<br>"
		}
	}

	fmt.Fprint(w, resp)
}

// Return dita of today
func ditaOfToday() string {
	resp := ""
	day := time.Now().Day()
	month := int(time.Now().Month())
	todayIs := ""

	// Analyze if today is a special day like Nadal
	switch {
	case day == 23 && month == 4:
		todayIs = "sant jordi"
	case day == 24 && month == 6:
		todayIs = "sant joan"
	case day == 1 && month == 11:
		todayIs = "tots sants"
	case day == 25 && month == 12:
		todayIs = "nadal"
	default:
		todayIs = "(" + fmt.Sprintf("%d", day) + " " + monthCat(month)
	}

	// Look if todayIs it is in todays. If so, build all today dites and return
	for _, v := range todays {
		if strings.Contains(strings.ToLower(v), todayIs) {
			resp += v + "<br>"
		}
	}

	if resp != "" {
		return resp
	}

	// If it didn't return, it is because today is a normal day. Return a random dita of the month
	return randomDitaMonth()
}

// Return random dita for the current month
func randomDitaMonth() string {
	idx := 0
	resp := ""
	month := int(time.Now().Month())

	// If it didn't return, it is because today is a normal day.Return a random dita of the month
	switch month {
	case 1:
		idx = rand.IntN(len(gen))
		resp = gen[idx]
	case 2:
		idx = rand.IntN(len(feb))
		resp = feb[idx]
	case 3:
		idx = rand.IntN(len(mar))
		resp = mar[idx]
	case 4:
		idx = rand.IntN(len(apr))
		resp = apr[idx]
	case 5:
		idx = rand.IntN(len(may))
		resp = may[idx]
	case 6:
		idx = rand.IntN(len(jun))
		resp = jun[idx]
	case 7:
		idx = rand.IntN(len(jul))
		resp = jul[idx]
	case 8:
		idx = rand.IntN(len(aug))
		resp = aug[idx]
	case 9:
		idx = rand.IntN(len(sep))
		resp = sep[idx]
	case 10:
		idx = rand.IntN(len(oct))
		resp = oct[idx]
	case 11:
		idx = rand.IntN(len(nov))
		resp = nov[idx]
	case 12:
		idx = rand.IntN(len(dec))
		resp = dec[idx]
	default:
		resp = ""
	}
	return resp
}

// Translates months (1, 2...) into catalan (gener, febrer...)
func monthCat(month int) string {
	switch month {
	case 1:
		return "gener"
	case 2:
		return "febrer"
	case 3:
		return "març"
	case 4:
		return "abril"
	case 5:
		return "maig"
	case 6:
		return "juny"
	case 7:
		return "juliol"
	case 8:
		return "agost"
	case 9:
		return "setembre"
	case 10:
		return "octubre"
	case 11:
		return "novembre"
	case 12:
		return "desembre"
	default:
		return "not a month"
	}
}

// Analyzes if a string contains (case-insensitive) all substrings. Returns a boolean
func containsAllSubstrs(str string, substrings []string) bool {
	strLower := strings.ToLower(str)
	allPresent := true

	for _, substring := range substrings {
		if !strings.Contains(strLower, strings.ToLower(substring)) {
			allPresent = false
			break
		}
	}
	return allPresent
}
