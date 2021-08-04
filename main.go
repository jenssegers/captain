package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/sahilm/fuzzy"
	"github.com/urfave/cli"
)

type Project struct {
	Path string
	Name string
}

type Config struct {
	Root      string
	Blacklist []string
	Depth     int
}

var config Config

func init() {
	usr, _ := user.Current()

	config = Config{
		Blacklist: []string{usr.HomeDir + "/Library", usr.HomeDir + "/Applications"},
		Root:      usr.HomeDir,
		Depth:     5,
	}

	depth, ok := os.LookupEnv("CAPTAIN_DEPTH")
	if ok {
		depth, err := strconv.Atoi(depth)
		if err == nil {
			config.Depth = depth
		}
	}

	root, ok := os.LookupEnv("CAPTAIN_ROOT")
	if ok {
		config.Root = root
	}
}

func scan(wg *sync.WaitGroup, folder string, depth int, results chan Project) {
	defer wg.Done()

	// Get all files and subdirectories in this directory.
	files, _ := ioutil.ReadDir(folder)
	var directories []string

	for _, file := range files {
		path := folder + "/" + file.Name()

		// Add subdirectories to list of yet to be scanned directories.
		if file.IsDir() {
			// Check if folder is in blacklist.
			for _, blacklist := range config.Blacklist {
				if blacklist == path {
					continue
				}
			}

			directories = append(directories, path)
			continue
		}

		// Search for docker-compose.yml file.
		if !file.IsDir() && (file.Name() == "docker-compose.yml" || file.Name() == "docker-compose.yaml") {
			results <- Project{
				Path: filepath.Dir(path),
				Name: strings.Trim(strings.Replace(filepath.Dir(path), config.Root, "", 1), "/"),
			}

			// No need to continue scan other subdirectories
			return
		}
	}

	// If no docker-compose.yml/docker-compose.yaml file was found, scan all subdirectories that we found.
	if depth > 1 {
		for _, folder := range directories {
			wg.Add(1)
			go scan(wg, folder, depth-1, results)
		}
	}

	return
}

func projects() []Project {
	channel := make(chan Project)
	var wg sync.WaitGroup

	wg.Add(1)
	go scan(&wg, config.Root, config.Depth, channel)

	// Turn channel into slice.
	projects := []Project{}
	go func() {
		for project := range channel {
			projects = append(projects, project)
		}
	}()

	wg.Wait()

	return projects
}

func match(projects []Project, pattern string) (Project, error) {
	dict := make(map[string]Project)
	for _, project := range projects {
		dict[project.Name] = project
	}

	list := make([]string, 0, len(projects))
	for _, project := range projects {
		list = append(list, project.Name)
	}

	matches := fuzzy.Find(pattern, list)

	if len(matches) > 0 {
		name := matches[0].Str
		return dict[name], nil
	}

	return projects[0], errors.New("no match found")
}

func dc(project Project, arg ...string) error {
	cmd := exec.Command("docker-compose", arg...)
	cmd.Dir = project.Path
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func search(pattern string) (Project, error) {
	return match(projects(), pattern)
}

func main() {
	app := cli.NewApp()
	app.Name = "captain"
	app.Usage = "Start and stop docker compose projects"
	app.Version = "0.3.2"

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"up", "sail"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "detach, d",
					Usage: "Start project in detached mode",
				},
			},
			Usage: "Start a docker compose project",
			Action: func(c *cli.Context) error {
				project, err := search(c.Args().Get(0))

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				fmt.Println("Starting " + project.Name + "\n")
				if c.Bool("detach") {
					dc(project, "up", "-d")
				} else {
					dc(project, "up")
				}

				return nil
			},
		},
		{
			Name:    "stop",
			Aliases: []string{"down", "dock"},
			Usage:   "Stop a docker compose project",
			Action: func(c *cli.Context) error {
				project, err := search(c.Args().Get(0))

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				fmt.Println("Stopping " + project.Name + "\n")
				dc(project, "stop")

				return nil
			},
		},
		{
			Name:  "restart",
			Usage: "Restart a docker compose project",
			Action: func(c *cli.Context) error {
				project, err := search(c.Args().Get(0))

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				fmt.Println("Restarting " + project.Name + "\n")
				dc(project, "restart")

				return nil
			},
		},
		{
			Name:  "build",
			Usage: "Build a docker compose project",
			Action: func(c *cli.Context) error {
				project, err := search(c.Args().Get(0))

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				fmt.Println("Building " + project.Name + "\n")
				dc(project, "build")

				return nil
			},
		},
		{
			Name:  "logs",
			Usage: "View container output from a docker compose project",
			Action: func(c *cli.Context) error {
				project, err := search(c.Args().Get(0))

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				dc(project, "logs")

				return nil
			},
		},
		{
			Name:  "abandon",
			Usage: "Stop all running docker containers",
			Action: func(c *cli.Context) error {
				fmt.Println("Stopping all containers")
				cmd := exec.Command("sh", "-c", "docker ps -q | xargs -n 1 -P 8 -I {} docker stop {}")
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				return cmd.Run()
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List available docker compose projects",
			Action: func(c *cli.Context) error {
				for _, project := range projects() {
					fmt.Println(project.Name)
				}
				return nil
			},
		},
		{
			Name:  "exec",
			Usage: "Executing command in a running service container",
			Action: func(c *cli.Context) error {
				project, err := search(c.Args()[0])

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				service := c.Args()[1]

				if len(service) <= 0 {
					fmt.Println("Missing service name")
					return nil
				}

				args := c.Args().Tail()[1:(len(c.Args()) - 1)]

				fmt.Println("Executing command in " + project.Name + " " + service + "\n")

				args = append([]string{"exec", service}, args...)
				dc(project, args...)
				return nil
			},
		},
		{
			Name:  "run",
			Usage: "Executing command as a new service container",
			Action: func(c *cli.Context) error {
				project, err := search(c.Args()[0])

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				service := c.Args()[1]

				if len(service) <= 0 {
					fmt.Println("Missing service name")
					return nil
				}

				args := c.Args().Tail()[1:(len(c.Args()) - 1)]

				fmt.Println("Running command in " + project.Name + " " + service + "\n")

				args = append([]string{"run", "--rm", service}, args...)
				dc(project, args...)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
