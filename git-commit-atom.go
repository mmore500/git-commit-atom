package main

import (
    "fmt"
    "bufio"
    "os"
    "path"
    "strings"
    "os/exec"
    "github.com/howeyc/fsnotify"
    "github.com/urfave/cli"
)

// little helper function
func check(e error) {
    if e != nil {
        panic(e)
    }
}


func setup_monitor_terminal(done_terminal chan bool){
  go func() {
    reader := bufio.NewReader(os.Stdin)
    for {
      text, _ := reader.ReadString('\n')
      if strings.Contains(text, "quit") || strings.Contains(text, "done") {
        done_terminal <- true
      }
    }
  }()
}

// if we are working with a COMMIT_EDITMSG file,
// Atom should be configured to tag it with the "magic" marker
// ###ATOM EDIT COMPLETE### so we know when the commit is complete
func handle_COMMIT_EDITMSG(filename string){
  done_atom := make(chan bool)
  watcher, file := setup_monitor_file(done_atom, filename)
  defer watcher.Close()
  defer file.Close()

  done_terminal := make(chan bool)
  setup_monitor_terminal(done_terminal)

  exec.Command("atom", filename).Run()
  select {
   case <-done_atom:
     fmt.Println("Commit completed by Atom.")
   case <-done_terminal:
     fmt.Println("Commit completed by terminal.")
  }
}

func setup_monitor_file(done_atom chan bool, filename string) (watcher *fsnotify.Watcher, file *os.File){
  // setup the watcher
  watcher, err := fsnotify.NewWatcher()
  check(err)

  // setup the file scanner
  file, err = os.Open(filename)
  check(err)

  // Process events
  go func() {
      for {
          select {
          case ev := <-watcher.Event:
            if ev.IsModify() {
              scanner := bufio.NewScanner(file)
              // scanner.Scan() advances to the next token returning false if an error was encountered
              var line string
              for scanner.Scan() {
                line = scanner.Text()
              }
              if strings.Contains(line, "##ATOM EDIT COMPLETE##") {
                done_atom <- true
              }
            }
          case err := <-watcher.Error:
            check(err)
          }
      }
  }()

  err = watcher.Watch(filename)
  check(err)

  return
}

func main(){
  app := cli.NewApp()
  app.Name = "git-commit-atom"
  app.Usage =
    `git-commit-atom filename
    signal commit completed by:
      entering "quit" or "done" at the terminal or
      appending ###ATOM EDIT COMPLETE### to the COMMIT_EDITMSG`
  app.Action = func(c *cli.Context) error {
    filename := c.Args().Get(0)


    if path.Base(filename) == "COMMIT_EDITMSG" || path.Base(filename) == "MERGE_MSG" {
      handle_COMMIT_EDITMSG(filename)
    } else {
      exec.Command("atom", "--wait", filename).Run()
    }
    return nil
  }

  app.Run(os.Args)

}
