package repo

import (
	"os/exec"
)

func CopyFile(src string) *exec.Cmd {
	cmd := exec.Command(
		"cp",
		"-r",
		src,
		".",
	)

	return cmd
}
