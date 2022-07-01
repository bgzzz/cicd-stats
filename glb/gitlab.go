package glb

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"

	gg "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type GLB struct {
	token  string
	client *gitlab.Client
}

func NewGLB(token string) (*GLB, error) {
	git, err := gitlab.NewClient(token)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create client")
	}

	return &GLB{client: git, token: token}, nil
}

func (glb *GLB) GetAllRepos(groupID int) ([]*gitlab.Project, error) {
	projects, _, err := glb.client.Groups.ListGroupProjects(groupID, &gitlab.ListGroupProjectsOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to list Projects")
	}

	// for _, project := range projects {
	// 	log.Println(project.PathWithNamespace)
	// }

	subgroups, _, err := glb.client.Groups.ListSubGroups(groupID, &gitlab.ListSubGroupsOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to list subgroups")
	}

	for _, subgroup := range subgroups {
		// log.Printf("subgroup: %v", subgroup.ID)
		// log.Printf("subgroup: %v", subgroup.Name)

		projectInSubgroup, err := glb.GetAllRepos(subgroup.ID)

		if err != nil {
			return nil, errors.Wrapf(err, "unable fetch from subgroup %v", subgroup.ID)
		}
		projects = append(projects, projectInSubgroup...)
	}

	return projects, nil
}

func (glb *GLB) ForceCloneProjects(projects []*gitlab.Project) ([]string, error) {

	if _, err := os.Stat("./repos"); !os.IsNotExist(err) {
		log.Println("removing repos")
		if err := os.RemoveAll("./repos"); err != nil {
			return []string{}, err
		}
	}

	dirs := []string{}
	for _, project := range projects {
		fmt.Println("#", project.HTTPURLToRepo)
		fmt.Printf("\"%d\",\n", project.ID)

		tmp := strings.Split(project.HTTPURLToRepo, "/")
		folderPath := fmt.Sprintf("./repos/%s",
			strings.Replace(tmp[len(tmp)-1], ".git", "", -1))
		_, err := gg.PlainClone(folderPath,
			false, &gg.CloneOptions{
				// The intended use of a GitHub personal access token is in replace of your password
				// because access tokens can easily be revoked.
				// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
				Auth: &http.BasicAuth{
					Username: "abc123", // yes, this can be anything except an empty string
					Password: glb.token,
				},
				URL:      project.HTTPURLToRepo,
				Progress: os.Stdout,
			})

		if err != nil {
			if err.Error() == "remote repository is empty" {
				continue
			}
			return []string{}, err
		}

		dirs = append(dirs, folderPath)
	}

	return dirs, nil
}
