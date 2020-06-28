package main

import (
	"fmt"
	"strings"
)

func main() {
	if len(jobs) < 1 {
		fmt.Println("=> no registered job yet!")
		return
	}

	fmt.Println("=> started sql2slack successfully ...")
	jobsNames := []string{}
	for k := range jobs {
		jobsNames = append(jobsNames, k)
	}

	fmt.Printf("=> available jobs:[%s]\n", strings.Join(jobsNames, ", "))

	cronhub.Run()
}
