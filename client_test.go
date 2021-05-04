package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func createProjectsFile(t *testing.T, data string) *Config {
	t.Helper()

	testDir := t.TempDir()
	f, err := os.Create(filepath.Join(testDir, "projects.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	f.WriteString(data)

	return &Config{
		DataDir: testDir,
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		projectsFile string
		want         Projects
	}{
		{
			projectsFile: `{"currentProjectName":"awesome","projects":{"awesome":{"name":"awesome"}}}`,
			want: Projects{
				{Name: "awesome"},
			},
		},
		{
			projectsFile: `{}`,
			want:         Projects{},
		},
		{
			projectsFile: ``,
			want:         Projects{},
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			projectsFilePath := createProjectsFile(t, tt.projectsFile)

			c, err := NewClient(projectsFilePath)
			if err != nil {
				t.Fatal(err)
			}

			got, err := c.List()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestGet(t *testing.T) {
	projectsFile := `{"currentProjectName":"awesome","projects":{"awesome":{"name":"awesome"}}}`

	cases := []struct {
		input string
		want  *Project
	}{
		{
			input: "awesome",
			want:  &Project{Name: "awesome"},
		},
		{
			input: "nonexistent",
			want:  &Project{},
		},
	}

	projectsFilePath := createProjectsFile(t, projectsFile)

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			c, err := NewClient(projectsFilePath)
			if err != nil {
				t.Fatal(err)
			}

			got, err := c.Get(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	projectsFile := `{"currentProjectName":"awesome","projects":{"awesome":{"name":"awesome","workspaces":["awesome-workspace"]}}}`

	cases := []struct {
		input Project
		want  *Project
	}{
		{
			input: Project{
				Name:       "new-project",
				Workspaces: []string{"/workspace-0", "/workspace-1"},
			},
			want: &Project{
				Name:       "new-project",
				Workspaces: []string{"/workspace-0", "/workspace-1"},
			},
		},
		{
			input: Project{
				Name:       "new-project",
				Workspaces: []string{},
			},
			want: &Project{
				Name:       "new-project",
				Workspaces: []string{},
			},
		},
		{
			input: Project{
				Name:       "new-project",
				Workspaces: nil,
			},
			want: &Project{
				Name:       "new-project",
				Workspaces: []string{},
			},
		},
		{
			input: Project{
				Name:       "awesome",
				Workspaces: []string{"awesome-workspace"},
			},
			want: &Project{
				Name:       "awesome",
				Workspaces: []string{"awesome-workspace"},
			},
		},
		{
			input: Project{
				Name:       "awesome",
				Workspaces: []string{"awesome-workspace", "new-awesome-workspace"},
			},
			want: &Project{
				Name:       "awesome",
				Workspaces: []string{"awesome-workspace"},
			},
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			projectsFilePath := createProjectsFile(t, projectsFile)

			c, err := NewClient(projectsFilePath)
			if err != nil {
				t.Fatal(err)
			}

			got, err := c.Create(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestGetCurrentProject(t *testing.T) {
	projectsFile := `{"currentProjectName":"awesome","projects":{"awesome":{"name":"awesome","workspaces":["awesome-workspace"]}}}`

	projectsFilePath := createProjectsFile(t, projectsFile)

	want := &Project{
		Name:       "awesome",
		Workspaces: []string{"awesome-workspace"},
	}

	c, err := NewClient(projectsFilePath)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.GetCurrentProject()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestChangeCurrentProject(t *testing.T) {
	projectsFile := `{"currentProjectName":"awesome-1","projects":{"awesome-1":{"name":"awesome-1","workspaces":["awesome-1-workspace"]},"awesome-2":{"name":"awesome-2","workspaces":["awesome-2-workspace"]}}}`

	cases := []string{
		"awesome-2",
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			projectsFilePath := createProjectsFile(t, projectsFile)

			c, err := NewClient(projectsFilePath)
			if err != nil {
				t.Fatal(err)
			}

			err = c.ChangeCurrentProject(tt)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestChangeCurrentProject_UnknownName(t *testing.T) {
	projectsFile := `{"currentProjectName":"awesome-1","projects":{"awesome-1":{"name":"awesome-1","workspaces":["awesome-1-workspace"]},"awesome-2":{"name":"awesome-2","workspaces":["awesome-2-workspace"]}}}`

	cases := []string{
		"nonexistent",
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			projectsFilePath := createProjectsFile(t, projectsFile)

			c, err := NewClient(projectsFilePath)
			if err != nil {
				t.Fatal(err)
			}

			err = c.ChangeCurrentProject(tt)
			if err == nil {
				t.Fatal("should be error but not")
			}
		})
	}
}
