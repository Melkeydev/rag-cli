package program

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true).Padding(1)
)

const bye = `⠀⠀⠀⢀⡴⠟⠛⢷⡄⠀⣠⠞⠋⠉⠳⡄⠀⠀⠀⠀
⠀⠀⠀⣸⠁⠀⠀⠈⣧⢰⠇⠀⠀⠀⢠⡇⠀⠀⠀⠀
⠀⠀⠀⠸⣆⠀⠀⠀⠘⣿⠀⠀⠀⠀⡞⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠹⣦⠀⠀⠀⠘⡄⠀⠀⠀⡇⠀⠀⠀⠀⠀
⠀⠀⠀⡴⠚⠙⠳⣀⡴⠂⠁⠒⢄⠀⢿⡀⠀⠀⠀⠀
⠀⠀⢸⡇⠀⢀⠔⠉⠀⠀⠀⡀⠀⠂⠘⣇⠀⠀⠀⠀
⠀⠀⠀⢳⡀⠘⢄⣀⣀⠠⠶⠄⠀⠀⠀⡿⠀⠀⠀⠀
⠀⠀⠀⠀⠻⣕⠂⠁⠀⠀⠀⠀⠀⠀⣰⡇⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠹⣅⠒⠀⠒⠂⠐⠒⢉⡼⢁⣤⠀⠀⠀
⠀⢀⣼⣻⣆⣤⢈⣙⣒⠶⠶⣶⣞⢿⣸⡟⠳⣾⣂⣤
⠸⢿⣽⠏⢠⡿⠋⣿⣭⢁⣈⣿⣽⣆⠙⢷⣄⠙⠋⠁
⠀⠀⠀⠀⠛⠃⠰⠿⣤⠄⠀⠸⠷⠟⠀⠀⠁⠀⠀⠀`

type Project struct {
	ProjectName string
	GitRepo     string
	AbsolutPath string
}

type Program struct {
	Loading bool
	Exit    bool
	Project Project
	OS      string
}

func (p *Program) ExitRag(tprogram *tea.Program) {
	if p.Exit {
		fmt.Printf("%s\n", logoStyle.Render(bye))
		tprogram.ReleaseTerminal()
		os.Exit(1)
	}
}

func executeCmd(name string, args []string, dir string) error {
	command := exec.Command(name, args...)
	command.Dir = dir
	var out bytes.Buffer
	command.Stdout = &out
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

func gitClone(repoUrl string, appDir string) {
	if err := executeCmd("git",
		[]string{"clone", "--depth", "1", "-b", "main", repoUrl, "."},
		appDir); err != nil {
		cobra.CheckErr(err)
	}
}

func checkAWSInstall() error {
	if err := executeCmd("aws",
		[]string{"--version"}, ""); err != nil {
		fmt.Println("this is the error", err)
		cobra.CheckErr(err)
		return err
	}

	return nil
}

func (p *Project) Create(wg *sync.WaitGroup, git, server string) {
	appDir := fmt.Sprintf("%s/%s", p.AbsolutPath, p.ProjectName)
	if _, err := os.Stat(p.AbsolutPath); err == nil {
		if err := os.Mkdir(appDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}

	// Determine what repo to pull
	if git == "Yes" {
		switch server {
		case "AWS Lambda":
			gitClone("https://github.com/Melkeydev/rag-stack-lambda.git", appDir)
		case "AWS Fargate":
			gitClone("https://github.com/Melkeydev/rag-stack-fargate.git", appDir)
		default:
			cobra.CheckErr("No Aws Deploy option selected")
		}
	}

	if err := os.RemoveAll(fmt.Sprintf("%s/.git", appDir)); err != nil {
		cobra.CheckErr(err)
	}

	wg.Done()
}

func (p *Project) InstallAWSCli(wg *sync.WaitGroup, system string) {

	if system == "darwin" {
		// Download the AWS CLI installer
		if err := executeCmd("curl",
			[]string{"https://awscli.amazonaws.com/AWSCLIV2.pkg", "-o", "AWSCLIV2.pkg"},
			""); err != nil {
			cobra.CheckErr(err)
		}

		// Install the AWS CLI
		if err := executeCmd("sudo",
			[]string{"installer", "-pkg", "AWSCLIV2.pkg", "-target", "/"},
			""); err != nil {
			cobra.CheckErr(err)
		}

	} else if system == "linux" {
		if err := executeCmd("curl",
			[]string{"https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip", "-o", "awscliv2.zip"},
			""); err != nil {
			fmt.Printf("%s\n", "Problem curling https://awscli.amazonaws.com")
		}
		if err := executeCmd("unzip",
			[]string{"awscliv2.zip"},
			""); err != nil {
			fmt.Printf("%s\n", "Problem Running unzip command on awscliv2.zip")
		}
		if err := executeCmd("sudo",
			[]string{"./aws/install"},
			""); err != nil {
			fmt.Printf("%s\n", "Found existing version of AWS CLI")
		}

		// Cleanup: delete awscliv2.zip and aws directory
		if err := os.Remove("awscliv2.zip"); err != nil {
			fmt.Println(err)
		}
		if err := os.RemoveAll("aws"); err != nil {
			fmt.Println(err)
		}
	}

	wg.Done()
}
