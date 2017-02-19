## git-commit-atom

**Author:** Matthew Moreno <br \>
**Contact:**  [matthew.andres.moreno@gmail.com](mailto:matthew.andres.moreno@gmail.com) <br \>
**Version:** February 18, 2017

------------
## Description

Allows git commit files to be edited in the current Atom pane, avoiding having
to launch another instance of Atom. The approach taken is directly inspired by
AJ Foster's "git-commit-atom.sh", [presented on his personal blog](https://aj-foster.com/2016/git-commit-atom/). However, this project is
implemented in Go in hopes of gaining portability and reliability.

------------
## Implementation

This project has two components: a standalone Go script that acts as the editor
called by Git during the commit process and CoffeeScript to be added to your
Atom `init.coffee` script. When the standalone Go script is activated, it opens
the `COMMIT_EDITMSG` file in the current Atom pane. When that file is closed,
Atom appends a "magic marker" to the end of the `COMMIT_EDITMSG` file. The Go
script, which is listening to the end of the `COMMIT_EDITMSG` file, recognizes
the "magic marker" and terminates, ending the commit edit session. In addition,
the Go script listens for user input at the terminal; the commit session can
also be ended by entering "quit" or "done" at the terminal. (This functionality
allows the standalone script to function in some capacity without the Atom init
script component in place).

------------
## Usage

To begin, you must have Go up and running on your machine, with a valid `GOPATH` and, if you are using macos, you might also need to set your `GOBIN`. This can be accomplished by adding the following lines to your `.bash_profile`.
~~~bash
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOPATH/bin
~~~

Next, add the following to your Atom init script (`init.coffee`). The following is directly quoted from AJ Foster.

~~~CoffeeScript
# This writes a magic token to the end of a commit message. We expect this to
# be run when the commit message editor has been closed.
#
commit_msg_notifier = (path) ->
  process = require("child_process")
  process.exec("echo \"##ATOM EDIT COMPLETE##\" >> " + path.replace /(\s)/g, '\\$1')

# The following looks at all new editors. If the editor is for a COMMIT_EDITMSG
# file, it sets up a callback for a magic token to be written when the editor
# is closed.
#
setup_commit_msg_notifier = (editor) ->
  if editor.buffer?.file?.getBaseName() == "COMMIT_EDITMSG"
    path = editor.buffer.file.getPath()
    editor.onDidDestroy ->
      commit_msg_notifier(path)

  # Return this, else weird things may happen. Anyone understand why?
  true

# The following looks at all new editors. If the editor is for a MERGE_MSG
# file, it sets up a callback for a magic token to be written when the editor
# is closed.
#
setup_merge_msg_notifier = (editor) ->
  if editor.buffer?.file?.getBaseName() == "MERGE_MSG"
    path = editor.buffer.file.getPath()
    editor.onDidDestroy ->
      commit_msg_notifier(path)

  # Return this, else weird things may happen. Anyone understand why?
  true


# Set up for all editors to be screened for commit messages.
atom.workspace.observeTextEditors(setup_commit_msg_notifier)
atom.workspace.observeTextEditors(setup_merge_msg_notifier)
~~~

Install the Go script and then configure Git to use the Go script.
~~~bash
go install github.com/mmore500/git-commit-atom
git config --global core.editor "git-commit-atom"
~~~
If you have preexisting Git repositories, you might have to use
~~~bash
git config core.editor "git-commit-atom"
~~~
on them.
