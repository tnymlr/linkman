[![Build Status](https://travis-ci.org/dikeert/linkman.svg?branch=master)](https://travis-ci.org/dikeert/linkman)

# Linkman

Linkman is a small command line tool to manage bookmarks,
written in Go.

**Features**:

 - Adding bookmarks.
 - Listing bookmarks.
 - Archiving bookmarks.

**Linkman supports**:

 - Automatic title fetching for supplied URLs.
 - Adding multiple bookmarks in one go.
 - Archiving multiple bookmarks in one go.
 - Maintaining multiple list of bookmarks.

## Installation

Right now the only installation option is `go get`:

```
$ go get github.com/dikeert/linkman
```

should download, build and install `linkman` into `bin` directory in `$GOPATH`.

I'll provide other installation methods as soon as I have time to do so,
contributions are welcome.


## Adding bookmarks

To add a bookmark into `default` list use `add` command and supply
one or more URLs:

```
$ linkman add http://url1.example.com http://url2.example.com
```


By default, linkman will go and fetch the webpages for supplied URLs
and store their titles alongside the URLs and does not allow for duplicates.

Additionally, linkman calculates and stores "source" of the supplied URLs,
that is second or third (depending on TLD) level domain name:

| url               | source        |
| ----              | ----          |
| youtube.com       | youtube       |
| stackoverflow.com | stackoverflow |
| domain.co.uk      | domain        |


`add` supports multiple options, that allow to:

 - Skip fetching the title: `--skip-title-fetch`
 - Provide custom title: `-t`, `--title`
 - Allow duplicates: `-f`, `--force`
 - Provide custom list name: `-l`, `--list`

**Example**

With [newsboat](https://newsboat.org/) one can use this command to store
RSS articles into "reading" list:

```
$ linkman add -l "reading" --skip-title-fetch -t "$title" "$url"
```

## Listing bookmarks

To print non-archived bookmarks from `default` list, use `list` command:

```
$ linkman list
```

Output:

```
ID:     2
Source: github
Title:  GitHub - asdine/storm: Simple and powerful toolkit for BoltDB
URL:    https://github.com/asdine/storm#update-an-object
```

`list` command supports multiple options that allow to:

 - include archived bookmarks: `-a`, `--archived`
 - specify output format: `-f`, `--format`
 - show bookmarks from the specified list: `-l`, `--list`
 - show *only* archvied bookmarks: `-A`, `--only-archived`
 - filter out bookmarks that have no title: `-T`, `--require-title`
 - show *only* bookmarks from the specified source: `-s`, `--source`
 - show *only* bookmarks which title contains specified string: 
   `-t`, `--title`

`list` allows for multiple filtering options. Whenever multiple filtering
options are provided, they are combined using `and` operation:

`$ linkman list -l reading -t github` will only show bookmarks from 
"reading" list that have "github" in the title.

### Output format

`list` command allows to specify output format using `-f` and `--format`
options.

Output format accepts templates in
[Go templates](https://golang.org/pkg/text/template/) form.

Each bookmark has several fields that can be specified in output format:

 - `ID`, bookmark identificator, assigned automatically on `add`
 - `URL`, the URL itself
 - `Source`, calculated source string (see `add` command for details)
 - `Title`, the title of the webpage behind the URLs
 - `List`, the list that bookmark belongs to

Output format supports special chars from C, such as `\n`, `\t` and so on.

**Example**

One can can show list of bookmarks using
[dmenu](https://tools.suckless.org/dmenu/) using command like this one:

```
$ linkman list -l reading -f '{{.ID}} {{.Title}}\n' | \
  dmenu -p "Reading list:" -l 30 -i
```

## Archiving bookmarks

to archive bookmarks use `archive` command and provide one or more IDs.

```
$ linkman archive $ID
```

`$ID` here is ID value from `list` output.


## Real life usage example

I use [newsboat](https://newsboat.org/) as my RSS reader. One of the features
of newsboat is support for bookmarks. With a simple keypress one can send
article information to a script. You can read more
[here](https://newsboat.org/releases/2.15/docs/newsboat.html#_bookmarking).

The script that I have is this one:

```bash

url="$1"
title="$2"

linkman add -l "reading" --skip-title-fetch -t "$title" "$url"
notify-send "Reading list:" "Added \"$title\""
```

It creates a new bookmark in `linkman` in the list "reading". It skips fetching
the title because the title is supplied by newsboat so the bookmarking process is
_fast_.

I use this script to create bookmarks for articles that I want to read later.


Now, whenever I feel like reading something, I use script like that:

```bash

local choice=$(linkman list -l reading -f '{{.ID}} {{.Title}}\n' |
    dmenu -p "Reading list:" -l 30 -i -fn "$FONT")

if ! [ "$choice" = "" ]; then
  local id=$(echo "$choice" | awk '{print $1}')
  local title="$(echo "$choice" | cut -d ' ' -f 2-)"

  linkman list -l reading -f '{{.ID}}\t{{.URL}}\n' |
    grep -e "^$id" |
    awk '{print $2}' |
    xargs -I '{}' firefox -P ReadingMode --name FirefoxReadingMode '{}' \
      && linkman archive "$id" \
      && notify-send "Reading list:" "Archived \"$title\""
fi
```

Which uses `linkman` to retrieve list of bookmarks in format:
`{{.ID}} {{.Titlte}}\n` and then passes it to
[dmenu](https://tools.suckless.org/dmenu/) so I could choose one of them.

It then opens the URL from chosen bookmark in special profile in Firefox which
has chromeless window: 

```
#TabsToolbar {
  visibility: collapse;
}

#titlebar {
  visibility: collapse;
}

#sidebar-header {
  visibility: collapse !important;
}

#nav-bar {
  visibility: collapse !important;
}
```

In order to retrieve the URL of the chosen bookmark it
uses `linkman` again, now requesting list of bookmarks in format:
`{{ID.}}\t{{.URL}}\n` and simply greping through it. It also archives the
article by ID whenever I'm done reading it and close Firefox.

This is my setup for distraction-free reading.

