package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	preHTML = `	
				<!DOCTYPE html>
				<html>
				<head>
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<link rel="icon" type="image/png" href="/res/favicon.png?"/>
				<link rel="apple-touch-icon" sizes="128x128" href="/res/and-dites.jpg">
				<style>
				p {
						font-family: sans-serif;
						font-size: 24px;
				}
				button {
					background-color: #4CAF50; /* Green */
					border: none;
					color: white;
					padding: 15px 32px;
					text-align: center;
					text-decoration: none;
					display: inline-block;
					font-size: 16px;
					margin: 4px 2px;
					cursor: pointer;
					border-radius: 4px;
					width: 250px;
				}
				input {
					font-size: 16px;
					width: 250px;
				}
				</style>
				<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
				<script>
				$(document).ready(function(){
					$("#btn1").click(function(){
					$("#p1").load("today");
					});
					$("#btn2").click(function(){
					$("#p1").load("misc");       
					});
					$("#btn3").click(function(){
					//pattern_value = $("input:text").val();
					pattern_value = $("#in1").val();
					url = "search?pattern=" + escape(pattern_value);
					$("#p1").load(url);
					});
				});
				</script>
				</head>
				<body>
				<p id="p1" align="center">
				`
	postHTML = `
				</p><br>
				<div align="center">
					<button id="btn1">Una altra dita del dia</button><br><br>
					<button id="btn2">Altres dites</button><br><br>
					<input id="in1" type="text" name="pattern" onkeydown='if (event.keyCode == 13) document.getElementById("btn3").click();'/><br>
						<button id="btn3">Cerca</button>
				</div>
				</body>
				</html>
				`
)

var (
	gen  = []string{}
	feb  = []string{}
	mar  = []string{}
	apr  = []string{}
	may  = []string{}
	jun  = []string{}
	jul  = []string{}
	aug  = []string{}
	sep  = []string{}
	oct  = []string{}
	nov  = []string{}
	dec  = []string{}
	none = []string{}

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

	// Create a new scanner
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		noMonth = false
		noSeason = false
		line := scanner.Text()

		switch {
		case strings.Contains(strings.ToLower(line), "gener"):
			gen = append(gen, line)
		case strings.Contains(strings.ToLower(line), "febrer"):
			feb = append(feb, line)
		case strings.Contains(strings.ToLower(line), "març"):
			mar = append(mar, line)
		case strings.Contains(strings.ToLower(line), "abril"):
			apr = append(apr, line)
		case strings.Contains(strings.ToLower(line), "maig"):
			may = append(may, line)
		case strings.Contains(strings.ToLower(line), "juny"):
			jun = append(jun, line)
		case strings.Contains(strings.ToLower(line), "juliol"):
			jul = append(jul, line)
		case strings.Contains(strings.ToLower(line), "agost"):
			aug = append(aug, line)
		case strings.Contains(strings.ToLower(line), "setembre"):
			sep = append(sep, line)
		case strings.Contains(strings.ToLower(line), "octubre"):
			oct = append(oct, line)
		case strings.Contains(strings.ToLower(line), "novembre"):
			nov = append(nov, line)
		case strings.Contains(strings.ToLower(line), "desembre"):
			dec = append(dec, line)
		default:
			noMonth = true
		}

		switch {
		case strings.Contains(strings.ToLower(line), "hivern"):
			dec = append(dec, line)
			gen = append(gen, line)
			feb = append(feb, line)
		case strings.Contains(strings.ToLower(line), "primavera"):
			mar = append(mar, line)
			apr = append(apr, line)
			may = append(may, line)
		case strings.Contains(strings.ToLower(line), "estiu"):
			jun = append(jun, line)
			jul = append(jul, line)
			aug = append(aug, line)
		case strings.Contains(strings.ToLower(line), "tardor"):
			sep = append(sep, line)
			oct = append(oct, line)
			nov = append(nov, line)
		default:
			noSeason = true
		}

		if noMonth && noSeason {
			none = append(none, line)
		}
	}

	// Check for any errors during the scan
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func main() {
	http.HandleFunc("/", showToday)
	http.HandleFunc("/misc", showMisc)

	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func showToday(w http.ResponseWriter, r *http.Request) {
	var idx int
	var resp string

	switch int(time.Now().Month()) {
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

	fmt.Fprintf(w, "%s", preHTML+resp+postHTML)
}

func showMisc(w http.ResponseWriter, r *http.Request) {
	idx := rand.IntN(len(none))
	resp := none[idx]

	fmt.Fprintf(w, "%s", preHTML+resp+postHTML)
}
