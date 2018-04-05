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
	path string
	name string
}

func scan(wg *sync.WaitGroup, folder string, depth int, results *[]project) {
	defer wg.Done()

	files, _ := ioutil.ReadDir(folder)
	var directories []string

	for _, file := range files {
		path := folder + "/" + file.Name()

		if file.IsDir() {
			directories = append(directories, path)
			continue
		}

		if !file.IsDir() && file.Name() == "docker-compose.yml" {
			*results = append(*results, project{
				path: filepath.Dir(path),
				name: filepath.Dir(path),
			})

			// Stop early, we found a docker compose project.
			return
		}
	}

	// If this was not a docker compose project, we should search the sub directories.
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
	projects := []project{}
	var wg sync.WaitGroup

	wg.Add(1)
	go scan(&wg, usr.HomeDir, 5, &projects)
	wg.Wait()

	return projects
}

func match(projects []project, pattern string) (project, error) {
	dict := make(map[string]project)
	for _, project := range projects {
		dict[project.path] = project
	}

	paths := make([]string, 0, len(projects))
	for _, project := range projects {
		paths = append(paths, project.path)
	}

	matches := fuzzy.Find(pattern, paths)

	if len(matches) > 0 {
		path := matches[0].Str
		return dict[path], nil
	}

	return projects[0], errors.New("no match found")
}

func up(project project) error {
	cmd := exec.Command("docker-compose", "up")
	cmd.Dir = project.path
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
			Usage:   "Start a docker compose project",
			Action: func(c *cli.Context) error {
				project, err := search(c.Args().Get(0))

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				fmt.Println("Starting " + project.name)
				up(project)

				return nil
			},
		},
		{
			Name:  "list",
			Usage: "List available docker compose projects",
			Action: func(c *cli.Context) error {
				for _, project := range projects() {
					fmt.Println(project.name)
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
