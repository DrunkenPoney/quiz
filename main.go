package main

import (
	"io/ioutil"
	"os"
	"strings"
	"bufio"
	"fmt"
	"strconv"
	"time"
)

type Arguments struct {
	limit int
	csv string
	help bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func Args() Arguments {
	args := Arguments{
		limit: 30,
		csv: "problems.csv",
		help: false,
	}

	var prev string
	for _, arg := range os.Args[1:] {
		switch strings.TrimSpace(prev) {
		case "-f":
		case "--file":
		case "--csv":
			args.csv = arg
			break;
		case "-l":
		case "--limit":
			i, err := strconv.ParseUint(arg, 10, 64)
			check(err)
			args.limit = int(i)
			break;
		case "--help":
		case "-h":
			fmt.Println("Usage of quiz:")
			fmt.Println("  -f --file --csv [string]")
			fmt.Println("   		a csv file in the format of 'question,answer' (default \"problems.csv\")")
			fmt.Println("  -l --limit [int]")
			fmt.Println("  		the time limit for the quiz in seconds (default : 30)")
			fmt.Println("  -h --help")
			fmt.Println("  		shows this message")
			os.Exit(0)
			break;
		}
		prev = arg
	}

	return args
}

func getFileName() string {
	isNext := false
	for _, arg := range os.Args[1:] {
		if isNext {
			return arg;
		}
		isNext = strings.Index(strings.TrimSpace(arg), "-f") == 0
	}
	return "problems.csv"
}

func doQuiz(lines []string, c chan <- int) {
	reader := bufio.NewReader(os.Stdin)
	quiz := make([][]string, len(lines)+1)
	goodAnswers := 0
	for i, line := range lines {
		quiz[i] = make([]string, 3)
		split := strings.Split(line, ",")
		quiz[i][0], quiz[i][1] = split[0], split[1]
		fmt.Print(quiz[i][0] + "=")
		quiz[i][2], _ = reader.ReadString('\n')
		quiz[i][2] = strings.TrimSpace(quiz[i][2])
		if quiz[i][1] == quiz[i][2] {
			goodAnswers++
		}
		fmt.Println(quiz[i])
	}
	c <- goodAnswers
}

func main() {
	args := Args()
	dat, err := ioutil.ReadFile(args.csv)
	check(err)
	c := make(chan int)
	lines := strings.Split(strings.TrimSpace(string(dat)), "\n")
	go doQuiz(lines, c)
	timeout := time.After(time.Duration(args.limit) * time.Second)
	select {
		case goodAnswers := <- c:
			fmt.Printf("Quiz done! Result: %d/%d", goodAnswers, len(lines))
		case <- timeout:
			fmt.Println("\nTimed out!")
		return
	}
}
