package gittools

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Originate/git-town/src/exit"
	"github.com/Originate/git-town/src/tools/command"
)

// DoesBranchHaveUnmergedCommits returns whether the branch with the given name
// contains commits that are not merged into the main branch
func DoesBranchHaveUnmergedCommits(branchName string) bool {
	return command.New("git", "log", GetMainBranch()+".."+branchName).Output() != ""
}

// GetLocalBranches returns the names of all branches in the local repository,
// ordered alphabetically
func GetLocalBranches() (result []string) {
	for _, line := range strings.Split(command.New("git", "branch").Output(), "\n") {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, line)
	}
	return
}

// GetLocalBranchesWithDeletedTrackingBranches returns the names of all branches
// whose remote tracking branches have been deleted
func GetLocalBranchesWithDeletedTrackingBranches() (result []string) {
	for _, line := range strings.Split(command.New("git", "branch", "-vv").Output(), "\n") {
		line = strings.Trim(line, "* ")
		parts := strings.SplitN(line, " ", 2)
		branchName := parts[0]
		deleteTrackingBranchStatus := fmt.Sprintf("[%s: gone]", GetTrackingBranchName(branchName))
		if strings.Contains(parts[1], deleteTrackingBranchStatus) {
			result = append(result, branchName)
		}
	}
	return
}

// GetLocalBranchesWithMainBranchFirst returns the names of all branches
// that exist in the local repository,
// ordered to have the name of the main branch first,
// then the names of the branches, ordered alphabetically
func GetLocalBranchesWithMainBranchFirst() (result []string) {
	mainBranch := GetMainBranch()
	result = append(result, mainBranch)
	for _, branch := range GetLocalBranches() {
		if branch != mainBranch {
			result = append(result, branch)
		}
	}
	return
}

// GetPreviouslyCheckedOutBranch returns the name of the previously checked out branch
func GetPreviouslyCheckedOutBranch() string {
	command := command.New("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if command.Err() != nil {
		return ""
	}
	return command.Output()
}

// GetTrackingBranchName returns the name of the remote branch
// that corresponds to the local branch with the given name
func GetTrackingBranchName(branchName string) string {
	return "origin/" + branchName
}

// HasBranch returns whether the repository contains a branch with the given name.
// The branch does not have to be present on the local repository.
func HasBranch(branchName string) bool {
	for _, line := range strings.Split(command.New("git", "branch", "-a").Output(), "\n") {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "remotes/origin/", "", 1)
		if line == branchName {
			return true
		}
	}
	return false
}

// HasTrackingBranch returns whether the local branch with the given name
// has a tracking branch.
func HasTrackingBranch(branchName string) bool {
	trackingBranchName := GetTrackingBranchName(branchName)
	for _, line := range strings.Split(command.New("git", "branch", "-r").Output(), "\n") {
		if strings.TrimSpace(line) == trackingBranchName {
			return true
		}
	}
	return false
}

// IsBranchInSync returns whether the branch with the given name is in sync with its tracking branch
func IsBranchInSync(branchName string) bool {
	if HasTrackingBranch(branchName) {
		localSha := GetBranchSha(branchName)
		remoteSha := GetBranchSha(GetTrackingBranchName(branchName))
		return localSha == remoteSha
	}
	return true
}

// ShouldBranchBePushed returns whether the local branch with the given name
// contains commits that have not been pushed to the remote.
func ShouldBranchBePushed(branchName string) bool {
	trackingBranchName := GetTrackingBranchName(branchName)
	command := command.New("git", "rev-list", "--left-right", branchName+"..."+trackingBranchName)
	return command.Output() != ""
}

// Helpers

// GetCurrentBranchNameDuringRebase returns the name of the current branch during a rebase
func GetCurrentBranchNameDuringRebase() string {
	filename := fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory())
	rawContent, err := ioutil.ReadFile(filename)
	exit.On(err)
	content := strings.TrimSpace(string(rawContent))
	return strings.Replace(content, "refs/heads/", "", -1)
}