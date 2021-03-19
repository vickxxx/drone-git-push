package repo

import (
	"os/exec"
)

func Rmfile(src string) *exec.Cmd {
	cmd := exec.Command(
		"rm",
		"-rf",
		src,
	)

	return cmd
}

func CopyFile(src string) *exec.Cmd {
	cmd := exec.Command(
		"cp",
		"-r",
		src,
		".",
	)

	return cmd
}

func ClearFile(filename string) *exec.Cmd {
	cmd := exec.Command(
		"clear.sh",
		filename,
	// "git",
	// "filter-branch",
	// "--force",
	// "--index-filter",
	// "\"",
	// "git",
	// "rm",
	// "-rf",
	// "--cached",
	// "--ignore-unmatch",
	// "app",
	// "\"",
	// "--prune-empty",
	// "--tag-name-filter",
	// "cat",
	// "--",
	// "--all",
	)

	return cmd
}
