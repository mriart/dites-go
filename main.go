// Dites
// Starts a web server that returns catalan dites
// Marc Riart Solans, 202402

package main

import (
	"bufio"
	"fmt"
	"html/template"
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

	noMonth  bool
	noSeason bool
)

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
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/today", handlerToday)
	http.HandleFunc("/misc", handlerMisc)
	http.HandleFunc("/search", handlerSearch)
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))

	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp := ditaOfToday()

	// Define data to pass to the template
	data := map[string]string{"Dita": resp}

	// Parse the template file and execute it
	tmpl, _ := template.ParseFiles("index.html")
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}

func handlerToday(w http.ResponseWriter, r *http.Request) {
	resp := randomDitaMonth()
	fmt.Fprintf(w, "%s", resp)
}

func handlerMisc(w http.ResponseWriter, r *http.Request) {
	idx := rand.IntN(len(none))
	resp := none[idx]

	fmt.Fprintf(w, "%s", resp)
}

func handlerSearch(w http.ResponseWriter, r *http.Request) {
	resp := ""
	// Parse the URL query string
	params := r.URL.Query()

	// Extract the "pattern" parameter. Returns a []string (every parameter can have several comma-separated values)
	pattern, ok := params["pattern"]
	if !ok {
		fmt.Fprintf(w, "Missing 'pattern' parameter in query string")
		return
	}
	fmt.Println(pattern)

	// Revisit dites.txt line by line. If found the pattern, append to the response
	file, err := os.Open("dites.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineLower := strings.ToLower(line)
		patternLower := strings.ToLower(pattern[0])

		if strings.Contains(lineLower, patternLower) {
			resp += line + "<br>"
		}
	}
	fmt.Fprintf(w, "%s", resp)
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
