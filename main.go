package main

import (
	"io/ioutil"
	"os"
	"strings"
	"bufio"
	"fmt"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
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

func main() {
	dat, err := ioutil.ReadFile(getFileName())
	check(err)

	reader := bufio.NewReader(os.Stdin)
	lines := strings.Split(strings.TrimSpace(string(dat)), "\n")
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

	fmt.Printf("Quiz done! Result: %d/%d", goodAnswers, len(quiz))
}