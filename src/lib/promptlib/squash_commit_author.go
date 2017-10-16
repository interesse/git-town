package promptlib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Originate/git-town/src/exit"
	"github.com/Originate/git-town/src/tools/cfmt"
	"github.com/Originate/git-town/src/tools/command"
	"github.com/Originate/git-town/src/tools/gittools"
	"github.com/Originate/git-town/src/tools/prompttools"
	"github.com/Originate/git-town/src/tools/stringtools"
	"github.com/fatih/color"
)

// GetSquashCommitAuthor gets the author of the supplied branch.
// If the branch has more than one author, the author is queried from the user.
func GetSquashCommitAuthor(branchName string) string {
	authors := getBranchAuthors(branchName)
	if len(authors) == 1 {
		return authors[0].NameAndEmail
	}
	cfmt.Printf(squashCommitAuthorHeaderTemplate, branchName)
	printNumberedAuthors(authors)
	fmt.Println()
	return askForAuthor(authors)
}

// Helpers

type branchAuthor struct {
	NameAndEmail    string
	NumberOfCommits string
}

var squashCommitAuthorHeaderTemplate = `
Multiple people authored the '%s' branch.
Please choose an author for the squash commit.

`

func askForAuthor(authors []branchAuthor) string {
	for {
		fmt.Print("Enter user's number or a custom author (default: 1): ")
		author, err := parseAuthor(GetUserInput(), authors)
		if err == nil {
			return author
		}
		prompttools.PrintError(err.Error())
	}
}

func getBranchAuthors(branchName string) (result []branchAuthor) {
	output := command.New("git", "shortlog", "-s", "-n", "-e", gittools.GetMainBranch()+".."+branchName).Output()
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result = append(result, branchAuthor{NameAndEmail: parts[1], NumberOfCommits: parts[0]})
	}
	return
}

func parseAuthor(userInput string, authors []branchAuthor) (string, error) {
	numericRegex, err := regexp.Compile("^[0-9]+$")
	exit.OnWrap(err, "Error compiling numeric regular expression")

	if numericRegex.MatchString(userInput) {
		return parseAuthorNumber(userInput, authors)
	}
	if userInput == "" {
		return authors[0].NameAndEmail, nil
	}
	return userInput, nil
}

func parseAuthorNumber(userInput string, authors []branchAuthor) (string, error) {
	index, err := strconv.Atoi(userInput)
	exit.OnWrap(err, "Error parsing string to integer")
	if index >= 1 && index <= len(authors) {
		return authors[index-1].NameAndEmail, nil
	}
	return "", errors.New("Invalid author number")
}

func printNumberedAuthors(authors []branchAuthor) {
	boldFmt := color.New(color.Bold)
	for index, author := range authors {
		stat := stringtools.Pluralize(author.NumberOfCommits, "commit")
		cfmt.Printf("  %s: %s (%s)\n", boldFmt.Sprintf("%d", index+1), author.NameAndEmail, stat)
	}
}