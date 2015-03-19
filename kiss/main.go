package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/daviddengcn/go-colortext"
	"github.com/gophergala/go_ne/core"
	//"github.com/mgutz/ansi"
)

var (
	taskName   = flag.String("task", "", "which task to run")
	configFile = flag.String("config", ".kiss.yml", "path to config file")
	group      = flag.String("group", "", "defines the server group for which the task should run")
)

func printOutput(str string, a ...interface{}) {
	ct.ChangeColor(ct.Green, false, ct.None, false)
	fmt.Println(fmt.Sprintf(str, a))
	ct.ResetColor()
}

func printError(str string, a ...interface{}) {
	ct.ChangeColor(ct.Red, false, ct.None, false)
	fmt.Println(fmt.Sprintf(str, a))
	ct.ResetColor()
}

func main() {
	flag.Parse()

	config, err := core.NewConfig()
	if err != nil {
		fail(err)
	}

	err = config.Load(*configFile)
	if err != nil {
		fail(err)
	}

	checkFlags()

	hosts, err := config.GetServerGroup(*group)
	if err != nil {
		fail(err)
	}

	//fmt.Println(ansi.Color(fmt.Sprintf("Running tasks for group `%v`", *group), "green"))
	printOutput("Running tasks for group `%v`", *group)

	for _, host := range hosts {
		//fmt.Println(ansi.Color(fmt.Sprintf("Selecting host `%v`", host.Host), "green"))
		printOutput("Selecting host `%v`", host.Host)

		var runner core.Runner
		if host.RunLocally {
			runner, err = core.NewLocalRunner()
			if err != nil {
				fail(err)
			}
			go core.LogOutput(runner.(*core.Local))

		} else {
			runner, err = core.NewRemoteRunner(host)
			if err != nil {
				fail(err)
			}
		}

		//fmt.Println(ansi.Color(fmt.Sprintf("Executing `%v`", *taskName), "green"))
		printOutput("Executing `%v`", *taskName)

		err = core.RunTask(runner, config, *taskName)
		if err != nil {
			fail(err)
		} else {
			//fmt.Println(ansi.Color("Tasks completed successfully", "green"))
			printOutput("Tasks completed successfully")
		}
	}
}

func checkFlags() {
	if len(*group) == 0 {
		fail(errors.New("Please select the target server group by passing the `-group=name` flag"))
	}

	if len(*taskName) == 0 {
		fail(errors.New("Please select the task to execute by passing the `-task=name` flag"))
	}
}

func fail(err error) {
	//fmt.Println(ansi.Color(fmt.Sprintf("Task failed: %v", err), "red"))
	printError("Task failed: %v", err)

	os.Exit(1)
}
