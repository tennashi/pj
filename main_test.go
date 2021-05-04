package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebdah/goldie/v2"
)

type testCommand struct {
	args         []string
	wantExitCode int
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}

	return nil
}

func testRun(t *testing.T, cwd, dataDir string, commands []testCommand) {
	t.Helper()

	curDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.Chdir(curDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = os.Chdir(cwd)
	if err != nil {
		t.Fatal(err)
	}

	for _, command := range commands {
		os.Args = append([]string{command.args[0], "--data-dir", dataDir}, command.args[1:]...)

		got := run()

		if diff := cmp.Diff(command.wantExitCode, got); diff != "" {
			t.Fatalf("exit code mismatch (-want +got):\n%s", diff)
		}
	}
}

func testCommands(t *testing.T, commands []testCommand) {
	t.Helper()

	tmpDir := t.TempDir()

	tmpStdout, err := os.CreateTemp(tmpDir, "stdout.*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tmpStderr, err := os.CreateTemp(tmpDir, "stderr.*.txt")
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = tmpStdout
	os.Stderr = tmpStderr

	err = copyFile("testdata/input/projects.json", filepath.Join(tmpDir, "projects.json"))
	if err != nil {
		t.Fatal(err)
	}

	testRun(t, "/tmp", tmpDir, commands)

	caseName := strings.ToLower(strings.ReplaceAll(t.Name(), "/", "_"))

	g := goldie.New(t)

	g.WithFixtureDir("testdata/stdout")
	gotStdout, err := os.ReadFile(os.Stdout.Name())
	if err != nil {
		t.Fatal(err)
	}
	g.Assert(t, caseName, gotStdout)

	g.WithFixtureDir("testdata/stderr")
	gotStderr, err := os.ReadFile(os.Stderr.Name())
	if err != nil {
		t.Fatal(err)
	}
	g.Assert(t, caseName, gotStderr)
}

func TestRun(t *testing.T) {
	cases := []testCommand{
		{args: []string{"pj"}, wantExitCode: 0},
		{args: []string{"pj", "-h"}, wantExitCode: 0},
		{args: []string{"pj", "init"}, wantExitCode: 0},
		{args: []string{"pj", "init", "awesome-project"}, wantExitCode: 0},
		{args: []string{"pj", "list"}, wantExitCode: 0},
		{args: []string{"pj", "get"}, wantExitCode: 0},
		{args: []string{"pj", "get", "awesome-project"}, wantExitCode: 0},
		{args: []string{"pj", "get", "awesome-project"}, wantExitCode: 0},
		{args: []string{"pj", "current"}, wantExitCode: 0},
	}

	for _, tt := range cases {
		t.Run(strings.Join(tt.args, "_"), func(t *testing.T) {
			testCommands(t, []testCommand{tt})
		})
	}
}

func TestRun_InitializeMultipleProject(t *testing.T) {
	testCmds := []testCommand{
		{args: []string{"pj", "init", "awesome-project"}, wantExitCode: 0},
		{args: []string{"pj", "init", "new-awesome-project"}, wantExitCode: 0},
	}

	testCommands(t, testCmds)
}

func TestRun_ChangeCurrentProject(t *testing.T) {
	testCmds := []testCommand{
		{args: []string{"pj", "init", "awesome-project"}, wantExitCode: 0},
		{args: []string{"pj", "change", "awesome-project"}, wantExitCode: 0},
	}

	testCommands(t, testCmds)
}

func TestRun_AddWorkspace(t *testing.T) {
	testCmds := []testCommand{
		{args: []string{"pj", "init", "awesome-project"}, wantExitCode: 0},
		{args: []string{"pj", "workspace", "add"}, wantExitCode: 0},
	}

	testCommands(t, testCmds)
}
