# gatoraid
Blog Aggregator in Go
---
## Intalling
* In order to get proper functioning you need to install postgres. Install postgres as follows:
  * macOS
    ```macOS bash
      $brew install postgresql@15
    ```
  * Linux
    ```linux bash
      sudo apt update
      sudo apt install postgresql postgresql-contrib
    ```

* Then the postgress server needs to be instantiated which can be done using:
  * macOS `brew services start postgresql@15`
  * Linux: `sudo service postgresql start`

* In order to compile the program go install `golang`. :P
* Along with `go` you'd need (optionally) `sqlc` for generating the sql go queries, and `goose` for setting up database in `postgres`.

## Setup
* Create a file `~/.gatoraidconfig.json` and add following details. And add following details:
  ```json
  {
    "db_url": <postgres connection link>,
    "current_user_name": _admin
  }
  ```
  * The `_admin` is a placeholder name until a user is registered.

* Setup postgres
  * After installing `goose`, run `goose <dburl> up` in `/sql/schema` directory.
---
## Installing
* You can generate the required binary by running `go build` in root directory, now you can use
  ```bash
    ./gatoraid <command> [arguements]
    ```
  to start using the program.
* In order to run the command from anywhere, you can also install the program by running `go install`.
  ```bash
    gatoraid <command> [arguements]
  ```
---
## Commands
* `register` pass register command with a unique username to generate and login as that user.
* `users` lists all registered users, with `(current)` badge in front of currently logged in user.
* `login` pass login command with unique, already registered user name. Once registered, if someone else is logged in you can log back in as your own user.
* `addfeed` passed with 2 arguement, Feed Name, and valid url for said rss feed. Checks and adds that feed to the system.
* `feeds` lists all the feeds in system, along with the users who first added the said feed.
* `follow` passed with one argument, a valid url which is already available in system--accessible by feeds command--allows currently logged in user to follow other users' feed.
* `following` lists all the feed the currently logged in user is following.
* `unfollow` passed with one argument, a valid url which is already available in system--accessible by feeds command--allows currently logged in user to unfollow said feed.
* `agg` updates the posts of the feeds that currently logged in user is following.
* `browse` provided optionally with a numeric arguement limit, lists the latest *limit* number of posts. Failure to provide *limit* will result in listing all the followed feeds' posts.
* `reset` (Maintainance Artifact) Deletes all the data from the database.

---
