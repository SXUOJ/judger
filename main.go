package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/isther/sandbox/container"
	"github.com/isther/sandbox/pkg/cgroup"
	"github.com/isther/sandbox/pkg/pipe"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `duck is a fake docker`

func main() {
	app := cli.NewApp()
	app.Name = "duck"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(c *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var (
	initCommand = cli.Command{
		Name:  "init",
		Usage: "Init container process run user's process in container",
		Action: func(context *cli.Context) error {
			logrus.Infof("init come on")

			err := container.RunContainerInitProcess()
			return err
		},
	}
	runCommand = cli.Command{
		Name: "run",
		Usage: `Create a container with namespace and cgroup limit
			   duck run -ti [command]`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "ti",
				Usage: "enable tty",
			},
			cli.StringFlag{
				Name:  "m",
				Usage: "memory limit",
			},
			cli.StringFlag{
				Name:  "cpuset",
				Usage: "cpuset limit",
			},
			cli.StringFlag{
				Name:  "cpu",
				Usage: "cpu limit",
			},
		},
		Action: func(context *cli.Context) error {
			if len(context.Args()) < 1 {
				return fmt.Errorf("Missing container command")
			}
			var cmdArray []string
			for _, arg := range context.Args() {
				cmdArray = append(cmdArray, arg)
			}

			tty := context.Bool("ti")
			resConf := &cgroup.ResourceConfig{
				MemoryLimit: context.String("m"),
				CpuSet:      context.String("cpuset"),
				CpuQuota:    context.String("cpu"),
			}

			Run(tty, cmdArray, resConf)
			return nil
		},
	}
)

func Run(tty bool, comArray []string, res *cgroup.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	//create cgroup amd limit source
	builder := cgroup.NewBuilder().AddType("V2")
	cg, err := builder.Build("test")
	defer cg.Destroy()
	if err != nil {
		logrus.Errorf("Build cgroup error")
		return
	}

	cpuq, _ := strconv.ParseUint(res.CpuQuota, 10, 64)
	cg.SetCPUQuota(cpuq)
	cg.SetCPUSet(res.CpuSet)
	mem, _ := strconv.ParseUint(res.MemoryLimit, 10, 64)
	cg.SetMemoryLimit(mem)
	cg.AddProc(uint64(parent.Process.Pid))

	pipe.SendInitCommand(comArray, writePipe)
	parent.Wait()
	os.Exit(0)
}
