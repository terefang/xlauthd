package main

import (
    "os"

    "github.com/terefang/gommons/pkg/subcmd"
)

const Banner = `------
         _____         __  __    ____ 
   _  __/ /   | __  __/ /_/ /_  / __ \
  | |/_/ / /| |/ / / / __/ __ \/ / / /
 _>  </ / ___ / /_/ / /_/ / / / /_/ / 
/_/|_/_/_/  |_\__,_/\__/_/ /_/_____/
`

func main() {
    exitcode := subcmd.Execute()
    os.Exit(exitcode)
}
