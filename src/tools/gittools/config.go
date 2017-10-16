/*
This file contains functionality around storing configuration settings
inside Git's metadata storage for the repository.
*/

package gittools

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Originate/git-town/src/exit"
	"github.com/Originate/git-town/src/tools/command"
)

// AddToPerennialBranches adds the given branch as a perennial branch
func AddToPerennialBranches(branchName string) {
	SetPerennialBranches(append(GetPerennialBranches(), branchName))
}

// DeleteAllAncestorBranches removes all Git Town ancestor entries
// for all branches from the configuration.
func DeleteAllAncestorBranches() {
	for _, key := range getConfigurationKeysMatching("^git-town-branch\\..*\\.ancestors$") {
		removeConfigurationValue(key)
	}
}

// DeleteParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func DeleteParentBranch(branchName string) {
	removeConfigurationValue("git-town-branch." + branchName + ".parent")
}

// GetAncestorBranches returns the names of all parent branches for the given branch,
// beginning but not including the parennial branch from which this hierarchy was cut.
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func GetAncestorBranches(branchName string) []string {
	value := getLocalConfigurationValue("git-town-branch." + branchName + ".ancestors")
	if value == "" {
		return []string{}
	}
	return strings.Split(value, " ")
}

// GetChildBranches returns the names of all branches for which the given branch
// is a parent.
func GetChildBranches(branchName string) (result []string) {
	for _, key := range getConfigurationKeysMatching("^git-town-branch\\..*\\.parent$") {
		parent := getLocalConfigurationValue(key)
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	return
}

// GetConfigurationValue returns the given configuration value,
// from either global or local Git configuration
func GetConfigurationValue(key string) (result string) {
	return command.New("git", "config", key).Output()
}

// GetGlobalConfigurationValue returns the global git configuration value for the given key
func GetGlobalConfigurationValue(key string) (result string) {
	if hasConfigurationValue("global", key) {
		result = command.New("git", "config", "--global", key).Output()
	}
	return
}

// GetMainBranch returns the name of the main branch.
func GetMainBranch() string {
	return getLocalConfigurationValue("git-town.main-branch-name")
}

// GetParentBranch returns the name of the parent branch of the given branch.
func GetParentBranch(branchName string) string {
	return getLocalConfigurationValue("git-town-branch." + branchName + ".parent")
}

// GetPerennialBranches returns all branches that are marked as perennial.
func GetPerennialBranches() []string {
	result := getLocalConfigurationValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// GetPullBranchStrategy returns the currently configured pull branch strategy.
func GetPullBranchStrategy() string {
	return GetLocalConfigurationValueWithDefault("git-town.pull-branch-strategy", "rebase")
}

// GetRemoteOriginURL returns the URL for the "origin" remote.
// In tests this value can be stubbed.
func GetRemoteOriginURL() string {
	if os.Getenv("GIT_TOWN_ENV") == "test" {
		mockRemoteURL := getLocalConfigurationValue("git-town.testing.remote-url")
		if mockRemoteURL != "" {
			return mockRemoteURL
		}
	}
	return command.New("git", "remote", "get-url", "origin").Output()
}

// GetRemoteUpstreamURL returns the URL of the "upstream" remote.
func GetRemoteUpstreamURL() string {
	return command.New("git", "remote", "get-url", "upstream").Output()
}

// GetURLHostname returns the hostname contained within the given Git URL.
func GetURLHostname(url string) string {
	hostnameRegex, err := regexp.Compile("(^[^:]*://([^@]*@)?|git@)([^/:]+).*")
	exit.OnWrap(err, "Error compiling hostname regular expression")
	matches := hostnameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return matches[3]
}

// GetURLRepositoryName returns the repository name contains within the given Git URL.
func GetURLRepositoryName(url string) string {
	hostname := GetURLHostname(url)
	repositoryNameRegex, err := regexp.Compile(".*" + hostname + "[/:](.+)")
	exit.OnWrap(err, "Error compiling repository name regular expression")
	matches := repositoryNameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return strings.TrimSuffix(matches[1], ".git")
}

// HasGlobalConfigurationValue returns whether there is a global git configuration for the given key
func HasGlobalConfigurationValue(key string) bool {
	return command.New("git", "config", "-l", "--global", "--name").OutputContainsLine(key)
}

// HasCompiledAncestorBranches returns whether the Git Town configuration
// contains a cached ancestor list for the branch with the given name.
func HasCompiledAncestorBranches(branchName string) bool {
	return len(GetAncestorBranches(branchName)) > 0
}

// HasRemote returns whether the current repository contains a Git remote
// with the given name.
func HasRemote(name string) bool {
	return command.New("git", "remote").OutputContainsLine(name)
}

// IsMainBranch returns whether the branch with the given name
// is the main branch of the repository.
func IsMainBranch(branchName string) bool {
	return branchName == GetMainBranch()
}

// RemoveAllConfiguration removes all Git Town configuration
func RemoveAllConfiguration() {
	command.New("git", "config", "--remove-section", "git-town").Output()
}

// SetAncestorBranches stores the given list of branches as ancestors
// for the given branch in the Git Town configuration.
func SetAncestorBranches(branchName string, ancestorBranches []string) {
	setConfigurationValue("git-town-branch."+branchName+".ancestors", strings.Join(ancestorBranches, " "))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func SetMainBranch(branchName string) {
	setConfigurationValue("git-town.main-branch-name", branchName)
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func SetParentBranch(branchName, parentBranchName string) {
	setConfigurationValue("git-town-branch."+branchName+".parent", parentBranchName)
}

// SetPerennialBranches marks the given branches as perennial branches
func SetPerennialBranches(branchNames []string) {
	setConfigurationValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func SetPullBranchStrategy(strategy string) {
	setConfigurationValue("git-town.pull-branch-strategy", strategy)
}

// UpdateOffline updates whether Git Town is in offline mode
func UpdateOffline(value bool) {
	setGlobalConfigurationValue("git-town.offline", strconv.FormatBool(value))
}

// UpdateShouldHackPush updates whether the current repository is configured to push
// freshly created branches up to the origin remote.
func UpdateShouldHackPush(value bool) {
	setConfigurationValue("git-town.hack-push-flag", strconv.FormatBool(value))
}

// Helpers

// GetConfigurationValueWithDefault returns the given configuration value,
// or the given default value if none is found.
func GetConfigurationValueWithDefault(key, defaultValue string) string {
	value := GetConfigurationValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getLocalConfigurationValue returns the given configuration value
// only from the local Git configuration
func getLocalConfigurationValue(key string) (result string) {
	if hasConfigurationValue("local", key) {
		result = command.New("git", "config", "--local", key).Output()
	}
	return
}

// GetLocalConfigurationValueWithDefault returns the given configuration value,
// and the given default value if it isnt given.
func GetLocalConfigurationValueWithDefault(key, defaultValue string) string {
	value := getLocalConfigurationValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getConfigurationKeysMatching(toMatch string) (result []string) {
	configRegexp, err := regexp.Compile(toMatch)
	exit.OnWrapf(err, "Error compiling configuration regular expression (%s): %v", toMatch, err)
	lines := command.New("git", "config", "-l", "--local", "--name").Output()
	for _, line := range strings.Split(lines, "\n") {
		if configRegexp.MatchString(line) {
			result = append(result, line)
		}
	}
	return
}

func hasConfigurationValue(location, key string) bool {
	return command.New("git", "config", "-l", "--"+location, "--name").OutputContainsLine(key)
}

func setConfigurationValue(key, value string) {
	command.New("git", "config", key, value).Run()
}

func setGlobalConfigurationValue(key, value string) {
	command.New("git", "config", "--global", key, value).Run()
}

func removeConfigurationValue(key string) {
	command.New("git", "config", "--unset", key).Run()
}