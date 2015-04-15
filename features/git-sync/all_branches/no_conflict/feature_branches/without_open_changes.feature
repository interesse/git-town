Feature: git sync --all: syncs all feature branches without open changes

  Background:
    Given I have feature branches named "feature1" and "feature2"
    And the following commits exist in my repository
      | BRANCH   | LOCATION         | MESSAGE         | FILE NAME     |
      | main     | remote           | main commit     | main_file     |
      | feature1 | local and remote | feature1 commit | feature1_file |
      | feature2 | local and remote | feature2 commit | feature2_file |
    And I have the following tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | remote-tag | remote   |
    And I am on the "feature1" branch
    When I run `git sync --all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                             |
      | feature1 | git fetch --prune                   |
      |          | git checkout main                   |
      | main     | git rebase origin/main              |
      |          | git checkout feature1               |
      | feature1 | git merge --no-edit origin/feature1 |
      |          | git merge --no-edit main            |
      |          | git push                            |
      |          | git checkout feature2               |
      | feature2 | git merge --no-edit origin/feature2 |
      |          | git merge --no-edit main            |
      |          | git push                            |
      |          | git checkout feature1               |
    And I am still on the "feature1" branch
    And all branches are now synchronized
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE                           | FILE NAME     |
      | main     | local and remote | main commit                       | main_file     |
      | feature1 | local and remote | feature1 commit                   | feature1_file |
      |          |                  | main commit                       | main_file     |
      |          |                  | Merge branch 'main' into feature1 |               |
      | feature2 | local and remote | feature2 commit                   | feature2_file |
      |          |                  | main commit                       | main_file     |
      |          |                  | Merge branch 'main' into feature2 |               |
    And I now have the following tags
      | NAME       | LOCATION         |
      | local-tag  | local            |
      | remote-tag | local and remote |