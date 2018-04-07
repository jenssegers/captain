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
	"sync"

	"github.com/sahilm/fuzzy"
	"github.com/urfave/cli"
)

type project struct {
	Path string
	Name string
}

func scan(wg *sync.WaitGroup, folder string, depth int, results chan project) {
	defer wg.Done()

	// Get all files and subdirectories in this directory.
	files, _ := ioutil.ReadDir(folder)
	var directories []string

	for _, file := range files {
		path := folder + "/" + file.Name()

		// Add subdirectories to list of yet to be scanned directories.
		if file.IsDir() {
			directories = append(directories, path)
			continue
		}

		// Search for docker-compose.yml file.
		if !file.IsDir() && file.Name() == "docker-compose.yml" {
			results <- project{
				Path: filepath.Dir(path),
				Name: filepath.Dir(path),
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

func projects() []project {
	usr, _ := user.Current()
	channel := make(chan project)
	var wg sync.WaitGroup

	wg.Add(1)
	go scan(&wg, usr.HomeDir, 5, channel)

	// Turn channel into slice.
	projects := []project{}
	go func() {
		for project := range channel {
			projects = append(projects, project)
		}
	}()

	wg.Wait()

	return projects
}

func match(projects []project, pattern string) (project, error) {
	dict := make(map[string]project)
	for _, project := range projects {
		dict[project.Path] = project
	}

	paths := make([]string, 0, len(projects))
	for _, project := range projects {
		paths = append(paths, project.Path)
	}

	matches := fuzzy.Find(pattern, paths)

	if len(matches) > 0 {
		path := matches[0].Str
		return dict[path], nil
	}

	return projects[0], errors.New("no match found")
}

func up(project project, daemon bool) error {
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

func down(project project) error {
	cmd := exec.Command("docker-compose", "down")
	cmd.Dir = project.Path
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func search(pattern string) (project, error) {
	return match(projects(), pattern)
}

func main() {
	app := cli.NewApp()
	app.Name = "captain"
	app.Usage = "Start and stop docker compose projects"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"up", "sail"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "daemon, d",
					Usage: "Start project in daemon mode",
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
				up(project, c.Bool("daemon"))

				return nil
			},
		},
		{
			Name:    "stop",
			Aliases: []string{"down"},
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
