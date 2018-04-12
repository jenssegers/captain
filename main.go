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
		if !file.IsDir() && file.Name() == "docker-compose.yml" {
			results <- Project{
				Path: filepath.Dir(path),
				Name: strings.Replace(filepath.Dir(path), config.Root, "", 1),
			}

			// No need to continue scan other subdirectories
			return
		}
	}

	// If no docker-compose.yml file was found, scan all subdirectories that we found.
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

func up(project Project, daemon bool) error {
	var cmd *exec.Cmd

	if daemon {
		cmd = exec.Command("docker-compose", "up", "-d")
	} else {
		cmd = exec.Command("docker-compose", "up")
	}

	cmd.Dir = project.Path
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func down(project Project) error {
	cmd := exec.Command("docker-compose", "stop")
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
	app.Version = "0.2.1"

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
				up(project, c.Bool("detach"))

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
				down(project)

				return nil
			},
		},
		{
			Name:  "abandon",
			Usage: "Stop all running docker containers",
			Action: func(c *cli.Context) error {
				fmt.Println("Stopping all containers\n")
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
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
