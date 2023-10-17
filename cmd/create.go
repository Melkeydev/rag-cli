package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/melkeydev/rag-cli/cmd/program"
	"github.com/melkeydev/rag-cli/cmd/steps"
	"github.com/melkeydev/rag-cli/cmd/ui/loading"
	multiInput "github.com/melkeydev/rag-cli/cmd/ui/multiInput"
	textinput "github.com/melkeydev/rag-cli/cmd/ui/textInput"

	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var (
	logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true).Padding(1)
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Spin up a full stack project in seconds",
	Long:  "RAG is a stack designed to get you developing full stack applications using React, AWS and Go",
	Run: func(cmd *cobra.Command, args []string) {

		logo := `
		██████╗  █████╗  ██████╗ 
		██╔══██╗██╔══██╗██╔════╝ 
		██████╔╝███████║██║  ███╗
		██╔══██╗██╔══██║██║   ██║
		██║  ██║██║  ██║╚██████╔╝
		╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ 
				`
		fmt.Printf("%s\n", logoStyle.Render(logo))

		options := steps.Options{
			AppName: &textinput.Output{},
		}

		program := &program.Program{
			OS: runtime.GOOS,
		}

		steps := steps.InitSteps(&options)

		tprogram := tea.NewProgram(textinput.InitialTextInputModel(options.AppName, "What is the name of your application?", program))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}
		program.ExitRag(tprogram)

		for _, step := range steps.Steps {
			s := &multiInput.Selection{}

			if step.StepName == "AWS" && options.AWSPrompt == "Yes" {
				continue
			}
			tprogram = tea.NewProgram(multiInput.InitialModelMulti(step.Options, s, step.Headers, program))

			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(err)
			}
			program.ExitRag(tprogram)

			*step.Field = s.Choice
		}

		project := program.Project
		currentWorkingDir, err := os.Getwd()
		if err != nil {
			cobra.CheckErr(err)
		}

		project.AbsolutePath = currentWorkingDir
		project.ProjectName = options.AppName.Output

		var initIOWg sync.WaitGroup

		if options.AWS == "Yes" {
			initIOWg.Add(2)

			go func() {
				defer initIOWg.Done()
				project.InstallAWSCli(&initIOWg, program.OS)
			}()

			go func() {
				defer initIOWg.Done()
				project.Create(&initIOWg, options.Server)
			}()

			tprogram = tea.NewProgram(loading.InitialAnimatedLoading())
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(err)
			}
			program.ExitRag(tprogram)
			initIOWg.Wait()
			tprogram.ReleaseTerminal()

		} else {
			initIOWg.Add(1)
			go project.Create(&initIOWg, options.Server)

			tprogram = tea.NewProgram(loading.InitialAnimatedLoading())
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(err)
			}
			program.ExitRag(tprogram)

			initIOWg.Wait()
			tprogram.ReleaseTerminal()

		}

	},
}
