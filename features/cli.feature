@phase-1
Feature: The command line
  The rules every command of a project made from this template follows, held by the
  built binary. The harness stamps the binary it builds with a known version and
  commit, as a release's build is stamped with the release's, so a scenario can tell
  a stamp from the dev version of a build without one. A usage failure is exit 2 and
  names what was wrong in words a person reads, never kong's own such as "EOL".

  @ID-CLI-01 @count-cli
  Scenario: --version prints the version stamped at build
    When go-template-itos runs with "--version"
    Then it exits with code 0
    And the first line of its standard output is "go-template-itos" and the stamped version

  @ID-CLI-02 @count-cli
  Scenario: version prints what --version prints
    When go-template-itos runs with "version"
    Then it exits with code 0
    And the first line of its standard output is "go-template-itos" and the stamped version
    And the second line of its standard output is "commit" and the stamped commit

  # The commit is the one a release's build stamps, else the revision the go tool
  # records for a build in a checkout, else there is none and no second line.
  @ID-CLI-03 @count-cli
  Scenario: --version prints the build's commit on its second line
    When go-template-itos runs with "--version"
    Then it exits with code 0
    And the second line of its standard output is "commit" and the stamped commit

  @ID-CLI-04 @count-cli
  Scenario: an unknown command exits 2, naming it and the command meant
    When go-template-itos runs with "coun"
    Then it exits with code 2
    And its error output says "coun"
    And its error output says "count"

  @ID-CLI-05 @count-cli
  Scenario: an unknown flag exits 2, naming it
    Given the file "notes.txt" contains "one two"
    When go-template-itos runs with "count notes.txt --bogus"
    Then it exits with code 2
    And its error output says "--bogus"

  # kong takes a switch given a value; the mappers and the hook on kong's parse
  # refuse it in our words instead.
  @ID-CLI-06 @count-cli
  Scenario: a switch given a value exits 2, naming it
    Given the file "notes.txt" contains "one two"
    When go-template-itos runs with "count notes.txt --json=yes"
    Then it exits with code 2
    And its error output says "--json"
    And its error output does not say "EOL"

  # Decision 1: a flag that is not cumulative is given once. A switch and its
  # --no- pair count as one flag, so giving both is giving it twice; the
  # environment variable is not the flag, and still takes a value beside it.
  @ID-CLI-08 @T-14
  Scenario: a switch given twice exits 2, naming it
    Given the file "notes.txt" contains "one two"
    When go-template-itos runs with "count notes.txt --json --json"
    Then it exits with code 2
    And its error output says "--json"
    And its standard output does not say "words"

  @ID-CLI-09 @T-14
  Scenario: a switch with its --no- pair exits 2, naming it
    Given the file "notes.txt" contains "one two"
    When go-template-itos runs with "count notes.txt --no-json --json"
    Then it exits with code 2
    And its error output says "--json"
    And its standard output does not say "words"

  @ID-CLI-07 @count-cli
  Scenario: version takes no argument and refuses one
    When go-template-itos runs with "version extra"
    Then it exits with code 2
    And its error output says "extra"
