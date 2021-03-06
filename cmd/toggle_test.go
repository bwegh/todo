package cmd

import (
	"fmt"
	. "github.com/svenfuchs/todo/io"
	. "github.com/svenfuchs/todo/test"
	"testing"
	"time"
)

func TestCmdToggleByIdFound(t *testing.T) {
	now := time.Now().Format("2006-01-02")
	src := NewMemoryIo("# Comment\n- foo [1]\nx bar [2]")
	out := NewMemoryIo("")
	args := &Args{Ids: []int{1}}

	ToggleCmd{Cmd{args, src, out}}.Run()
	actual := src.ReadLines()
	expected := []string{"# Comment", fmt.Sprintf("x foo done:%s [1]", now), "x bar [2]"}

	AssertEqual(t, actual, expected)
}

func TestCmdToggleByTextFound(t *testing.T) {
	src := NewMemoryIo("# Comment\n- foo [1]\nx bar done:2015-12-13 [2]")
	out := NewMemoryIo("")
	args := &Args{Text: "bar"}

	ToggleCmd{Cmd{args, src, out}}.Run()
	actual := src.ReadLines()
	expected := []string{"# Comment", "- foo [1]", "- bar [2]"}

	AssertEqual(t, actual, expected)
}

func TestCmdToggleByProjectsFound(t *testing.T) {
	now := time.Now().Format("2006-01-02")
	src := NewMemoryIo("- foo +baz [1]\nx bar +baz [2]")
	out := NewMemoryIo("")
	args := &Args{Projects: []string{"baz"}}

	ToggleCmd{Cmd{args, src, out}}.Run()
	actual := src.ReadLines()
	expected := []string{fmt.Sprintf("x foo +baz done:%s [1]", now), "- bar +baz [2]"}

	AssertEqual(t, actual, expected)
}

func TestCmdToggleReportsToggled(t *testing.T) {
	src := NewMemoryIo("# Comment\n- foo [1]\nx bar done:2015-12-13 [2]")
	out := NewMemoryIo("")
	args := &Args{Ids: []int{2}, Report: true}

	ToggleCmd{Cmd{args, src, out}}.Run()
	actual := out.ReadLines()
	expected := []string{"Toggled 1 item:", "", "- bar [2]"}

	AssertEqual(t, actual, expected)
}
