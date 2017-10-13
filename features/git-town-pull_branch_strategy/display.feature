Feature: passing an invalid option to the pull strategy configuration

  As a user or tool configuring Git Town's pull branch strategy
  I want to know what the existing value for the pull-strategy is
  So I can decide whether to I want to adjust it.


  Scenario: default setting
    When I run `git-town pull-branch-strategy`
    Then Git Town prints
      """
      rebase
      """


  Scenario: explicit rebase
    Given the "pull-branch-strategy" configuration is set to "rebase"
    When I run `git-town pull-branch-strategy`
    Then Git Town prints
      """
      rebase
      """


  Scenario: explicit merge
    Given the "pull-branch-strategy" configuration is set to "merge"
    When I run `git-town pull-branch-strategy`
    Then Git Town prints
      """
      merge
      """
