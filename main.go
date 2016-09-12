package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

var url = flag.String("url", "", "Jira search URL")
var ticketRegexp = regexp.MustCompile(`^\s*([a-zA-Z]+-\d+):?.*`)
var errInvalidLine = errors.New("invalid line")

func main() {
	flag.Parse()

	if len(*url) == 0 {
		fmt.Fprintln(os.Stderr, "Search URL must be provided")
		flag.Usage()
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)

	count := 0
	for scanner.Scan() {
		ticket, err := parseTicket(scanner.Text())
		if err == nil {
			fmt.Printf("Opening ticket %s...\n", ticket)
			if err = openTicket(*url, ticket); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to open ticket:", err)
			}
			count += 1
		} else {
			fmt.Fprintln(os.Stderr, "Invalid line:", scanner.Text())
		}
	}

	fmt.Println("Ticket count:", count)

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func parseTicket(line string) (string, error) {
	matches := ticketRegexp.FindStringSubmatch(line)
	if matches == nil {
		return "", errInvalidLine
	}
	return matches[1], nil
}

func openTicket(url, ticket string) error {
	cmd := exec.Command("open", url+ticket)
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}
