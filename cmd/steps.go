package cmd

type StepSchema struct {
	StepName string
	Options  []string
	Headers  string
	Field    *string
}

type Steps struct {
	Steps []StepSchema
}

func initSteps(options *Options) *Steps {
	steps := &Steps{
		Steps: []StepSchema{
			{
				StepName: "Git",
				Options:  []string{"Yes", "No thank you"},
				Headers:  "Do you want to init a git project?",
				Field:    &options.Git,
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
