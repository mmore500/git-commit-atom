## git-commit-atom
------------
Want to use Atom's handy `COMMIT_EDITMSG` syntax highlighting?
🙋

Tired of waiting on Atom to open a new window with the `--wait` option?
:hourglass: :sleeping:

Together with [sister Atom package `git-edit-atom`](https://atom.io/packages/git-edit-atom), this Go script allows Git commit files to be conveniently edited in the current editor pane... avoiding the launch of another instance of Atom!
:star2: :smirk:

![A screenshot of git-edit-atom and git-commit-atom in action together](https://thumbs.gfycat.com/BaggyFreshBoaconstrictor-size_restricted.gif)

## Prerequisites
 * A working installation of Go tools (need it? see [here](https://golang.org/doc/install)).
 * be sure your `GOBIN` is included in your `PATH`! If your `PATH` isn't already properly configured, try adding the following lines to your `.bash_profile`.
 ~~~bash
 export GOPATH=$HOME/go
 export GOBIN=$GOPATH/bin
 export PATH=$PATH:$GOPATH/bin
 ~~~
 * [Sister Atom package `git-edit-atom`](https://atom.io/packages/git-edit-atom), highly recommended.

## Installation
Install the Go script and then configure Git to use the Go script.
At the terminal, you'll need to run the following commands.
~~~bash
go install github.com/mmore500/git-commit-atom
git config --global core.editor "git-commit-atom"
~~~
If you have preexisting Git repositories, you might have to use
~~~bash
git config core.editor "git-commit-atom"
~~~
on them.

## Usage
Once `git-commit-atom` is configured as Git's editor, Git `COMMIT_EDITMSG`, `TAG_EDITMSG`, `MERGE_MSG`, `git-rebase-todo`, and `.diff` files will open in the current pane of Atom.
To complete the message editing process simply close the tab (`cmd-w` is convenient) if the Atom package `git-edit-atom` is installed or, if not, enter `quit` or `done` at the terminal.

## Implementation
This project has two components: a standalone Go script that acts as the editor called by Git during the commit process and the Atom package `git-edit-atom`.
When the standalone Go script is activated, it opens the `COMMIT_EDITMSG` file in the current Atom pane.
When that file is closed, Atom appends a "magic marker" (`## ATOM EDIT COMPLETE##`) to the end of the `COMMIT_EDITMSG` file.
The Go script, which is listening to the end of the `COMMIT_EDITMSG` file, recognizes the "magic marker" and terminates, ending the commit edit session.
In addition, the Go script listens for user input at the terminal.
The commit session can also be ended by entering `quit` or `done`.
(This functionality allows the standalone script to function in some capacity without the Atom package in place).
This approach is directly inspired by AJ Foster's "git-commit-atom.sh", [presented on his personal blog](https://aj-foster.com/2016/git-commit-atom/).
However, this project is implemented in Go and as an Atom package in hopes of gaining portability and reliability.
