![](https://raw.githubusercontent.com/b4b4r07/screenshots/master/gch/logo.png)

[![Build Status](https://img.shields.io/travis/b4b4r07/gch.svg?style=flat-square)][travis]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GitHub release](http://img.shields.io/github/release/b4b4r07/gch.svg?style=flat-square)][release]

[travis]: https://travis-ci.org/b4b4r07/gch
[license]: https://raw.githubusercontent.com/b4b4r07/dotfiles/master/doc/LICENSE-MIT.txt
[release]: https://github.com/b4b4r07/gch/releases

List the changes that are applied to the repository in the `$GOPATH/src` directory.

## Description

`ghq` provides a way to organize remote repository clones, like `go get` does. However, when we are operating multiple projects, there is a disadvantage that the progress of the project in `$GOPATH` is hard to watch or check.

`gch` recursively list up a project directory in the `$GOPATH`, and get that `git status`.

*==> [`ghq`](https://github.com/motemen/ghq) + [`peco`](https://github.com/peco/peco) + [`gch`](https://github.com/b4b4r07/gch) = BEST Development Environment !*

***DEMO:***

![](https://raw.githubusercontent.com/b4b4r07/screenshots/master/gch/demo.png)

## Requirements

- Go 1.3+
- git
- (ghq)
- (peco)

## Usage

List modified project in all $GOPATHs.

```console
$ gch
$GOPATH/src/github.com/b4b4r07/cdinterface
 M cdinterface.sh
?? .keep
$GOPATH/src/github.com/b4b4r07/gomi
 M README.md
$GOPATH/src/github.com/b4b4r07/xtime
 M README.md
```

View paths only.

```console
$ gch -l
/home/b4b4r07/src/github.com/b4b4r07/cdinterface
/home/b4b4r07/src/github.com/b4b4r07/gomi
/home/b4b4r07/.go/src/github.com/b4b4r07/xtime
```

For more information, see `gch --help`.

| Simbol | Meaning |
|---|---|
| A | addition of a file |
| C | copy of a file into a new one |
| D | deletion of a file |
| M | modification of the contents or mode of a file |
| R | renaming of a file |
| T | change in the type of the file |
| U | file is unmerged (you must complete the merge before it can be committed) |
| X | "unknown" change type (most probably a bug, please report it) |
| ?? | untracked |
| ... | ... |

For more information, see `man git-diff-files`.

## Installation

	$ go get github.com/b4b4r07/gch

or

	$ ghq get b4b4r07/gch
	$ go install $GOPATH/src/github.com/b4b4r07/gch

## License

[MIT](https://raw.githubusercontent.com/b4b4r07/dotfiles/master/doc/LICENSE-MIT.txt) © BABAROT (a.k.a. b4b4r07)
