<p align="center">
  <img alt="Gopher" src="logo.svg" height="250" />
  <h3 align="center">Updater</h3>
  <p align="center">Semantic updater and migrator for GoLang executables.</p>
  <p align="center">
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
    <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-77%25-brightgreen.svg?longCache=true&style=flat)</a>
    <a href="https://goreportcard.com/report/github.com/ainsleyclark/updater"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/ainsleyclark/updater"></a>
    <a href="https://pkg.go.dev/github.com/ainsleyclark/updater"><img src="https://godoc.org/github.com/ainsleyclark/updater?status.svg" alt="GoDoc"></a>
  </p>
</p>

## Why? [![start with why](https://img.shields.io/badge/start%20with-why%3F-brightgreen.svg?style=flat)](http://www.ted.com/talks/simon_sinek_how_great_leaders_inspire_action)

Updater aims to unify semantic migrations and executable updates in GoLang. You can seamlessly update executables and 
run any SQL database migrations that depend on the specific version of the application. Callbacks can be passed
to each migration allowing you to edit environment variables, or a directory structure.

## Installation

```bash
go get -u github.com/ainsleyclark/updater
```

## Example

### Creating the updater
To create an updater, simply call `updater.New()` with options `Updater.Options{}` 

```go
u, err := updater.New(updater.Options{
    GithubURL: "https://github.com/ainsleyclark/verbis", // The URL of the Git Repos
    Version:       "v0.0.1", // The currently running version
    Verify:        false, // Updates will be verified by checking the new exec with -version
    DB:            nil, // Pass in an sql.DB for a migration
})

if err != nil {
    log.Fatal(err)
}

status, err := u.Update(fmt.Sprintf("verbis_v0.0.2_%s_%s.zip", runtime.GOOS, runtime.GOARCH))
if err != nil {
    return
}

fmt.Println(status)
```

## Adding a migration
Migrations are stored in memory, so you can call `AddMigration` from anywhere with as version number, SQL statement
(optional) and CallBack functions (optional). 

```go
func init() {
	err := updater.AddMigration(&updater.SQL{
		Version:      "v0.0.2", // The version of the migration
		SQL:    strings.NewReader("UPDATE my_table SET 'title' WHERE id = 1"),
		CallBackUp:   func() error { return nil }, // Runs on up of migration.
		CallBackDown: func() error { return nil }, // Runs on error of migration.
		Stage:        updater.Patch, // Can be Patch, Major or Minor.
	})

	if err != nil {
		log.Fatal(err)
	}
}
```

## Credits

Shout out to [go-rocket-update](https://github.com/mouuff/go-rocket-update) for providing an excellent API for self updating executables.


