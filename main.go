package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bgzzz/cicd-stats/glb"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := createCtlApp().Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func createCtlApp() *cli.App {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "hosting-platform",
			Usage:       "hosting-platform is platform for hosting the code in the remote repo",
			Value:       "gitlab",
			DefaultText: "gitlab",
			EnvVars:     []string{"HOSTING_PLATFORM"},
		},
		&cli.StringFlag{
			Name:        "projects-folder",
			Usage:       "projects-folder it the folder that contains all needed git projects to evaluate",
			Value:       "./",
			DefaultText: "./",
			EnvVars:     []string{"PROJECTS_FOLDER"},
		},
		&cli.StringFlag{
			Name:    "gitlab-token",
			Usage:   "gitlab-token is used to access gitlab",
			EnvVars: []string{"GITLAB_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "gitlab-project-id",
			Usage:   "gitlab-project-id determines project id of the gitlab",
			EnvVars: []string{"GITLAB_PROJECT_ID"},
		},
	}

	return &cli.App{
		Name:  "cicd-stats",
		Usage: "cicd-stats go over the projects in the folder and uses git to calculate and further expose metrics",
		Flags: flags,
		Commands: []*cli.Command{
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "update all the projects",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "evaluate",
				Aliases: []string{"e"},
				Usage:   "evaluate git metrics and print it out",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "evaluate-api",
				Aliases: []string{"e-api"},
				Usage:   "evaluate api metrics and print it out",
				Action: func(c *cli.Context) error {
					token := c.String("gitlab-token")
					projectID, err := strconv.ParseInt(c.String("gitlab-project-id"),
						10, 64)
					if err != nil {
						return err
					}

					repos, err := glb.NewGLB(token)
					if err != nil {
						return err
					}

					projects, err := repos.GetAllRepos(int(projectID))
					if err != nil {
						return err
					}

					for _, project := range projects {
						fmt.Println("#", project.HTTPURLToRepo)
						fmt.Printf("\"%d\",\n", project.ID)
					}

					return nil
				},
			},
			{
				Name:    "all",
				Aliases: []string{"e-api"},
				Usage:   "updates, evaluate all metrics and ship it to db",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}
}
