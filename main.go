package main 

import(
	"github.com/docopt/docopt-go"
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"strconv"
	"math"
	"time"
)

type Result struct {
	Slope int
	Length int
	Visited bool
}

type SkiMap struct {
	Data []int
	Width int
	Height int
}

func main() {
	usage := `Ski Map - Longest and Steepest Run.
Usage:
  skimap [--input=<file>]
  skimap -h | --help
  skimap --version

Options:
  -h --help       Show this screen.
  --version       Show version.
  --input=<file>  File to Parse [default: map.txt]
`

	//Parse the program usage and retrieve arguments map
	arguments, _ := docopt.Parse(usage, nil, true, "Ski Map 1.0", false)

	skiMap := readSkiMap(arguments["--input"].(string))

	start := time.Now()
	skiPath := findLongestAndSteepestPath(skiMap)
	elapsed := time.Since(start)

	fmt.Printf("Operation took: %s\n", elapsed)
	fmt.Printf("Email: %d%d@redmart.com\n", skiPath.Length, skiPath.Slope)
}

func findLongestAndSteepestPath(skiMap *SkiMap) Result {
	result := make([]Result, skiMap.Width*skiMap.Height)

	for k, _ := range skiMap.Data {
		//Start sloping from each starting point if there is no result for it
		if result[k].Length == 0 {
			result[k] = slopeFrom(skiMap, k, &result)
			//Count the first element (origin)
			result[k].Length += 1
		}
	}

	return maxResult(result)
}

func maxResult(arr []Result) Result {
	max := Result{0,0, false}

	for i := 0; i < len(arr); i++ {
		//Replace max if greater
		if arr[i].Length > max.Length {
			max = arr[i]
		//Replace max if equal but of greater slope
		} else if arr[i].Length == max.Length {
			if arr[i].Slope > max.Slope {
				max = arr[i]
			}
		}
	}

	return max
}


func calcResult(skiMap *SkiMap, result *Result, results *[]Result, start, finish int) {
	//Calculate only if path is possible (n,s,e,w) and if slope of finish is lower
	if finish >= 0 && skiMap.Data[finish] < skiMap.Data[start] {
		//If finish is already visited, use the calculated values
		if (*results)[finish].Visited == true {
			result.Length += (*results)[finish].Length + 1
			result.Slope += skiMap.Data[start] - skiMap.Data[finish] + (*results)[finish].Slope
		//If not visited, calculate the next spot's values
		} else {
			res := slopeFrom(skiMap, finish, results)
			result.Length += res.Length + 1
			result.Slope += skiMap.Data[start] - skiMap.Data[finish] + res.Slope
			//We have visited this spot
			res.Visited = true
			//Save result so we don't have to calculate this again
			(*results)[finish] = res
		}
	} 
}

func slopeFrom(skiMap *SkiMap, start int, results *[]Result) Result {

	//Calculate who is the next value on the graph
	n := north(start, skiMap.Width, skiMap.Height)
	s := south(start, skiMap.Width, skiMap.Height)
	e := east(start, skiMap.Width, skiMap.Height)
	w := west(start, skiMap.Width, skiMap.Height)

	nResult, sResult, wResult, eResult := Result{}, Result{}, Result{}, Result{}

	//Calculate neighboring values
	calcResult(skiMap, &nResult, results, start, n)
	calcResult(skiMap, &sResult, results, start, s)
	calcResult(skiMap, &eResult, results, start, e)
	calcResult(skiMap, &wResult, results, start, w)

	arr := []Result{nResult, sResult, wResult, eResult}
	//Result will be the maximum value with maximum slope
	res := maxResult(arr)

	return res
} 

func north(position, width, height int) int {
	if position < width {
		return -1
	} else {
		return int(position - width)
	}
}

func south(position, width, height int) int {
	if position + width >= width * height {
		return -1
	} else {
		return int(position + width)
	}
}

func east(position, width, height int) int {
	if math.Mod(float64(position+1), float64(width)) == 0 {
		return -1
	} else {
		return int(position+1)
	}
}

func west(position, width, height int) int {
	if math.Mod(float64(position), float64(width)) == 0 {
		return -1
	} else {	
		return int(position-1)
	}
}

func readSkiMap(file string) *SkiMap {

	skiMap := &SkiMap{}

	//Open specified file
	f, err := os.Open(file)
	
	if err != nil {
		log.Fatal("Could not open file:", file, err)
	}

	//Read the header info to know our matrix size
	scanner := bufio.NewScanner(f)
	scanner.Scan()

	//Parse the header information
	header := strings.Split(scanner.Text(), " ")

	if len(header) != 2 {
		log.Fatal("Invalid header. Header elements must be 2. Found ", len(header) , " elements.")
	}	

	if height, err := strconv.ParseInt(header[0], 10, 64); err != nil {
		log.Fatal("Header:", err)
	} else {
		skiMap.Height = int(height)
	}

	if width, err := strconv.ParseInt(header[1], 10, 64); err != nil {
		log.Fatal("Header:", err)
	} else {
		skiMap.Width = int(width)
	}

	//Parse the Matrix	
	skiMap.Data = []int{}

	row := 0
	//Go through each line
	for scanner.Scan() {
		line := scanner.Text()
		row++
		wordScanner := bufio.NewScanner(strings.NewReader(line))
		wordScanner.Split(bufio.ScanWords)
		col := 0
		//Go through each number in line
		for wordScanner.Scan() {
			col++
			num, err := strconv.ParseInt(wordScanner.Text(), 10, 64)

			if err != nil {
				log.Fatal("Could not read number at:", row, col)
			}

			skiMap.Data = append(skiMap.Data, int(num))
		}
		//Line does not comply matrix header number of elements
		if col != int(skiMap.Width) {
			log.Fatal("Number of elements does not match header at line ", row)
		}
	}

	//Matrix does not comply with number of rows
	if row != int(skiMap.Height) {
		log.Fatal("Number of lines does not match header, ", row)
	}

	//There was an error reading the file
	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading file:",file, err)
	}

	return skiMap
}