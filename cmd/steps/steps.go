package steps

import textinput "github.com/melkeydev/rag-cli/cmd/ui/textInput"

type StepSchema struct {
	StepName string
	Options  []string
	Headers  string
	Field    *string
}

type Steps struct {
	Steps []StepSchema
}

type Options struct {
	Server    string
	AppName   *textinput.Output
	AWSPrompt string
	AWS       string
}

func InitSteps(options *Options) *Steps {
	steps := &Steps{
		Steps: []StepSchema{
			{
				StepName: "AWS Prompt",
				Options:  []string{"Yes", "No"},
				Headers:  "RAG requires AWS. Do you have AWS CLI installed?",
				Field:    &options.AWSPrompt,
			},
			{
				StepName: "AWS",
				Options:  []string{"Yes", "No thank you"},
				Headers:  "Do you want RAG to install the AWS CLI on your machine?",
				Field:    &options.AWS,
			},
			{
				StepName: "Server",
				Options:  []string{"AWS Lambda", "AWS Fargate"},
				Headers:  "How do you want to deploy your app?",
				Field:    &options.Server,
			},
		},
	}

	return steps
}
