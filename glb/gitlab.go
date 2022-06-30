package glb

import (
	"log"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type GLB struct {
	client *gitlab.Client
}

func NewGLB(token string) (*GLB, error) {
	git, err := gitlab.NewClient(token)
	log.Println(token)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create client")
	}

	return &GLB{client: git}, nil
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
