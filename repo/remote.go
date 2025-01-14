package repo

import (
	"os/exec"
)

// RemoteRemove drops the defined remote from a git repo.
func RemoteRemove(name string) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"remote",
		"rm",
		name)

	return cmd
}

// RemoteAdd adds an additional remote to a git repo.
func RemoteAdd(name, url string) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"remote",
		"add",
		name,
		url)

	return cmd
}

// RemotePush pushs the changes from the local head to a remote branch..
func RemotePush(remote, branch string, force bool, followtags bool) *exec.Cmd {
	return RemotePushNamedBranch(remote, "HEAD", branch, force, followtags)
}

// RemotePushNamedBranch puchs changes from a local to a remote branch.
func RemotePushNamedBranch(remote, localbranch string, branch string, force bool, followtags bool) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"push",
		remote,
		localbranch+":"+branch)

	if force {
		cmd.Args = append(
			cmd.Args,
			"--force")
	}

	if followtags {
		cmd.Args = append(
			cmd.Args,
			"--follow-tags")
	}

	return cmd
}

func RemoteCloneNamedBranch(remote, branch string) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"clone",
		"-b",
		branch,
		// "--depth=1",
		remote,
		".",
	)

	return cmd
}

func GitTag(tagName string) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"tag",
		tagName,
	)
	return cmd
}

func TagPush(remote string) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"push",
		remote,
		"--tags",
	)
	return cmd
}

func ForcePush(remote string) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"push",
		"-f",
	)
	return cmd
}
