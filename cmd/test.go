package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	multiInput "ragCli/cmd/ui/multiInput"
	textinput "ragCli/cmd/ui/textInput"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

type Options struct {
	Server  string
	Git     string
	AppName *textinput.Output
}

type Project struct {
	projectName string
	gitRepo     string
	AbsolutPath string
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

func (p *Project) Create(wg *sync.WaitGroup, options Options) {
	appDir := fmt.Sprintf("%s/%s", p.AbsolutPath, p.projectName)
	if _, err := os.Stat(p.AbsolutPath); err == nil {
		if err := os.Mkdir(appDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}

	// Determine what repo to pull
	if options.Git == "Yes" {

		switch options.Server {
		case "AWS Lambda":
			fmt.Println("did it hit here?")
			gitClone("https://github.com/Melkeydev/rag-stack-lambda.git", appDir)
		case "AWS Fargate":
			gitClone("https://github.com/Melkeydev/rag-stack-fargate.git", p.AbsolutPath)
		default:
			cobra.CheckErr("No Aws Deploy option selected")
		}
	}

	if err := os.RemoveAll(fmt.Sprintf("%s/.git", appDir)); err != nil {
		cobra.CheckErr(err)
	}

	wg.Done()
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "test",
	Long:  `test long form`,
	Run: func(cmd *cobra.Command, args []string) {

		options := Options{
			AppName: &textinput.Output{},
		}

		steps := initSteps(&options)

		program := tea.NewProgram(textinput.InitialTextInputModel(options.AppName, "What is the name of your application?"))
		if _, err := program.Run(); err != nil {
			cobra.CheckErr(err)
		}

		for _, step := range steps.Steps {
			s := &multiInput.Selection{}
			program = tea.NewProgram(multiInput.InitialModelMulti(step.Options, s, step.Headers))

			if _, err := program.Run(); err != nil {
				cobra.CheckErr(err)
			}

			*step.Field = s.Choice
		}

		project := Project{}
		currentWorkingDir, err := os.Getwd()
		if err != nil {
			cobra.CheckErr(err)
		}

		project.AbsolutPath = currentWorkingDir
		project.projectName = options.AppName.Output

		var initIOWg sync.WaitGroup
		initIOWg.Add(1)
		go project.Create(&initIOWg, options)

		initIOWg.Wait()
		program.ReleaseTerminal()
	},
}
