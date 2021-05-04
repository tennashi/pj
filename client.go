package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	projectsFilePath string
}

type Config struct {
	DataDir string `json:"dataDir"`
}

func (cfg *Config) ensureDataDir() error {
	return os.MkdirAll(cfg.DataDir, 0755)
}

func NewClient(cfg *Config) (*Client, error) {
	err := cfg.ensureDataDir()
	if err != nil {
		return nil, err
	}

	projectsFilePath := filepath.Join(cfg.DataDir, "projects.json")

	f, err := os.OpenFile(projectsFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if fInfo.Size() < int64(len("{}")) {
		f.WriteString("{}")
	}

	return &Client{
		projectsFilePath: projectsFilePath,
	}, nil
}

type ProjectsFile struct {
	CurrentProjectName string             `json:"currentProjectName"`
	Projects           map[string]Project `json:"projects"`
}

type Project struct {
	Name             string   `json:"name"`
	Workspaces       []string `json:"workspaces"`
	CurrentWorkspace string   `json:"currentWorkspace"`
}

func (p *Project) Summary() []string {
	workspaceDirNames := make([]string, 0, len(p.Workspaces))
	for _, workspace := range p.Workspaces {
		workspaceDirNames = append(workspaceDirNames, filepath.Base(workspace))
	}

	return []string{
		p.Name,
		p.CurrentWorkspace,
		strings.Join(workspaceDirNames, ","),
	}
}

type Projects []*Project

func (p Projects) Header() []string {
	return []string{"NAME", "CURRENT_WORKSPACE", "WORKSPACES"}
}

func (p Projects) Summaries() []Summarable {
	if p == nil {
		return nil
	}

	ret := make([]Summarable, 0, len(p))
	for _, proj := range p {
		ret = append(ret, proj)
	}
	return ret
}

func (c *Client) List() (Projects, error) {
	f, err := os.Open(c.projectsFilePath)
	if err != nil {
		return nil, err
	}

	projectsFile := ProjectsFile{}
	d := json.NewDecoder(f)
	err = d.Decode(&projectsFile)
	if err != nil {
		return nil, err
	}

	projects := projectsFile.Projects

	res := make([]*Project, 0, len(projects))
	for _, project := range projects {
		project := project
		res = append(res, &project)
	}

	return res, nil
}

func (c *Client) GetCurrentProject() (*Project, error) {
	f, err := os.Open(c.projectsFilePath)
	if err != nil {
		return nil, err
	}

	projectsFile := ProjectsFile{}
	d := json.NewDecoder(f)
	err = d.Decode(&projectsFile)
	if err != nil {
		return nil, err
	}

	currentProject := projectsFile.Projects[projectsFile.CurrentProjectName]

	return &currentProject, nil
}

func (c *Client) Get(projectName string) (*Project, error) {
	f, err := os.Open(c.projectsFilePath)
	if err != nil {
		return nil, err
	}

	projectsFile := ProjectsFile{}
	d := json.NewDecoder(f)
	err = d.Decode(&projectsFile)
	if err != nil {
		return nil, err
	}

	projects := projectsFile.Projects

	project := projects[projectName]
	return &project, nil
}

func (c *Client) Create(project Project) (*Project, error) {
	f, err := os.Open(c.projectsFilePath)
	if err != nil {
		return nil, err
	}

	projectsFile := ProjectsFile{}
	d := json.NewDecoder(f)
	err = d.Decode(&projectsFile)
	if err != nil {
		return nil, err
	}

	projects := projectsFile.Projects
	if projects == nil {
		projects = map[string]Project{}
	}

	if project, ok := projects[project.Name]; ok {
		return &project, nil
	}

	workspacePathes := make([]string, 0, len(project.Workspaces))
	for _, workspace := range project.Workspaces {
		workspacePath, err := filepath.Abs(filepath.Clean(workspace))
		if err != nil {
			return nil, err
		}

		workspacePathes = append(workspacePathes, workspacePath)
	}

	project.Workspaces = workspacePathes

	projects[project.Name] = project

	projectsFile.Projects = projects
	projectsFile.CurrentProjectName = project.Name

	f, err = os.OpenFile(c.projectsFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	e := json.NewEncoder(f)
	err = e.Encode(projectsFile)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) ChangeCurrentProject(projectName string) error {
	f, err := os.Open(c.projectsFilePath)
	if err != nil {
		return err
	}

	projectsFile := ProjectsFile{}
	d := json.NewDecoder(f)
	err = d.Decode(&projectsFile)
	if err != nil {
		return err
	}

	if _, ok := projectsFile.Projects[projectName]; !ok {
		return errors.New("not found")
	}

	projectsFile.CurrentProjectName = projectName

	f, err = os.OpenFile(c.projectsFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	e := json.NewEncoder(f)
	err = e.Encode(projectsFile)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Update(project *Project) error {
	if project == nil {
		return errors.New("project should not be nil")
	}

	f, err := os.Open(c.projectsFilePath)
	if err != nil {
		return err
	}

	projectsFile := ProjectsFile{}
	d := json.NewDecoder(f)
	err = d.Decode(&projectsFile)
	if err != nil {
		return err
	}

	if _, ok := projectsFile.Projects[project.Name]; !ok {
		return errors.New("not found")
	}

	projectsFile.Projects[project.Name] = *project

	f, err = os.OpenFile(c.projectsFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	e := json.NewEncoder(f)
	err = e.Encode(projectsFile)
	if err != nil {
		return err
	}

	return nil
}
