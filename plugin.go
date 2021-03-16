package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/appleboy/drone-git-push/repo"
	"github.com/davecgh/go-spew/spew"
	"github.com/gookit/color"
)

type (
	// Netrc structure
	Netrc struct {
		Machine  string
		Login    string
		Password string
	}

	// Commit structure
	Commit struct {
		Author Author
	}

	// Author structure
	Author struct {
		Name  string
		Email string
	}

	// Config structure
	Config struct {
		Key           string
		Remote        string
		RemoteName    string
		Branch        string
		LocalBranch   string
		Path          string
		Force         bool
		FollowTags    bool
		SkipVerify    bool
		Commit        bool
		CommitMessage string
		EmptyCommit   bool
		NoVerify      bool
		CopySrcLst    []string
		DefaultSrc    string
	}

	// Plugin Structure
	Plugin struct {
		Netrc  Netrc
		Commit Commit
		Config Config
	}
)

// Exec starts the plugin execution.
func (p Plugin) Exec() error {
	if err := p.HandlePath(); err != nil {
		return err
	}
	color.Yellow.Println("HandlePath")

	if err := p.WriteConfig(); err != nil {
		return err
	}
	color.Yellow.Println("WriteConfig")

	if err := p.WriteKey(); err != nil {
		return err
	}

	color.Yellow.Println("WriteKey")

	if err := p.WriteNetrc(); err != nil {
		return err
	}

	color.Yellow.Println("WriteNetrc")

	if os.Getenv("PLUGIN_SHOW_ENV") != "" {
		cmd := exec.Command("ls", "-a")

		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Execute Shell:%s failed with error:%s", cmd, err.Error())
			return nil
		}
		fmt.Printf("Execute Shell:%s finished with output:\n%s", cmd, string(output))

		cmd2 := exec.Command("env")

		output2, err := cmd2.Output()
		if err != nil {
			fmt.Printf("Execute Shell:%s failed with error:%s", cmd2, err.Error())
			return nil
		}
		fmt.Printf("Execute Shell:%s finished with output:\n%s", cmd2, string(output2))
		fmt.Println(spew.Sdump(p))
	}

	if err := p.WriteToken(); err != nil {
		return err
	}
	color.Yellow.Println("WriteToken")

	if err := p.Clone(); err != nil {
		return err
	}

	color.Yellow.Println("Clone")

	if err := p.CopyFile(); err != nil {
		return err
	}

	if err := p.HandleCommit(); err != nil {
		return err
	}
	color.Yellow.Println("HandleCommit")

	if err := p.HandleRemote(); err != nil {
		return err
	}

	color.Yellow.Println("HandleRemote")

	if err := p.HandlePush(); err != nil {
		return err
	}

	color.Yellow.Println("HandlePush")

	return p.HandleCleanup()
}

// WriteConfig writes all required configurations.
func (p Plugin) WriteConfig() error {
	if err := repo.GlobalName(p.Netrc.Login).Run(); err != nil {
		return err
	}

	if err := repo.GlobalUser(p.Commit.Author.Email).Run(); err != nil {
		return err
	}

	if p.Config.SkipVerify {
		if err := repo.SkipVerify().Run(); err != nil {
			return err
		}
	}

	return nil
}

// WriteKey writes the private SSH key.
func (p Plugin) WriteKey() error {
	return repo.WriteKey(
		p.Config.Key,
	)
}

// WriteNetrc writes the netrc config.
func (p Plugin) WriteNetrc() error {
	return repo.WriteNetrc(
		p.Netrc.Machine,
		p.Netrc.Login,
		p.Netrc.Password,
	)
}

// WriteToken writes token.
func (p Plugin) WriteToken() error {
	var err error

	p.Config.Remote, err = repo.WriteToken(
		p.Config.Remote,
		p.Netrc.Login,
		p.Netrc.Password,
	)

	return err
}

// HandleRemote adds the git remote if required.
func (p Plugin) HandleRemote() error {
	// color.Red.Println(spew.Sdump(p))
	if p.Config.Remote != "" {
		if err := execute(repo.RemoteAdd(p.Config.RemoteName, p.Config.Remote)); err != nil {
			return err
		}
	}

	return nil
}

// HandlePath changes to a different directory if required
func (p Plugin) HandlePath() error {

	if p.Config.Path == "" {
		return nil
	}

	color.Red.Println(os.Getwd())

	splitedRemote := strings.Split(p.Config.Remote, "/")
	dirName := ""
	if len(splitedRemote) > 0 {
		dirName = splitedRemote[len(splitedRemote)-1]
	}
	dirName = strings.ReplaceAll(dirName, ".git", "")
	dirName = filepath.Join(p.Config.Path, dirName)
	color.Yellow.Println(dirName)

	// path is not exist
	if _, err := os.Stat(dirName); err == nil {
		os.RemoveAll(dirName)
	}

	// if _, err := os.Stat(p.Config.Path); os.IsNotExist(err) {
	// 	_ = os.MkdirAll(dirName, os.ModePerm)

	// }
	_ = os.MkdirAll(dirName, os.ModePerm)

	if err := os.Chdir(dirName); err != nil {
		return err
	}

	return nil
}

func (p Plugin) CopyFile() error {

	for _, filename := range p.Config.CopySrcLst {
		filename := strings.TrimSpace(filename)
		if strings.HasPrefix(filename, "/") {
			continue
		}

		fromPath := filepath.Join(p.Config.DefaultSrc, filename)

		// color.Green.Println(fromPath, p.Config.DefaultSrc)
		if err := execute(repo.CopyFile(fromPath)); err != nil {
			continue
		}
	}

	return nil
}

// HandleCommit commits dirty changes if required.
func (p Plugin) HandleCommit() error {
	if p.Config.Commit {
		if err := execute(repo.Add()); err != nil {
			return err
		}

		cmtMsg := fmt.Sprintf("%s-%s", os.Getenv("DRONE_COMMIT_SHA")[:10], p.Config.CommitMessage)

		if err := execute(repo.TestCleanTree()); err != nil {
			// changes to commit
			if err := execute(repo.ForceCommit(cmtMsg, p.Config.NoVerify)); err != nil {
				return err
			}
		} else { // no changes
			if p.Config.EmptyCommit {
				// no changes but commit anyway
				if err := execute(repo.EmptyCommit(cmtMsg, p.Config.NoVerify)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// HandlePush pushs the changes to the remote repo.
func (p Plugin) HandlePush() error {
	var (
		name       = p.Config.RemoteName
		local      = p.Config.LocalBranch
		branch     = p.Config.Branch
		force      = p.Config.Force
		followtags = p.Config.FollowTags
	)

	return execute(repo.RemotePushNamedBranch(name, local, branch, force, followtags))
}

// HandleCleanup does eventually do some cleanup.
func (p Plugin) HandleCleanup() error {
	if p.Config.Remote != "" {
		if err := execute(repo.RemoteRemove(p.Config.RemoteName)); err != nil {
			return err
		}
	}

	return nil
}

func (p Plugin) Clone() error {
	if p.Config.Remote != "" {
		if err := execute(repo.RemoteCloneNamedBranch(p.Config.Remote, p.Config.Branch)); err != nil {
			return err
		}
	}

	return nil
}
