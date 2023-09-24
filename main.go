package main

import (
	"github.com/AlecAivazis/survey/v2"
	instance_connect "github.com/khaitranhq/aws-works/internal/tasks/instance-connect"
)

type Task struct {
	Description string
	Command     func()
}

func selectTask(tasks []Task) Task {
	tasksDescription := []string{}
	for _, task := range tasks {
		tasksDescription = append(tasksDescription, task.Description)
	}

	selectedTaskDecription := ""
	prompt := &survey.Select{
		Message: "Choose a task",
		Options: tasksDescription,
	}
	survey.AskOne(prompt, &selectedTaskDecription)

	var selectedTask Task
	for _, task := range tasks {
		if task.Description == selectedTaskDecription {
			selectedTask = task
		}
	}
	return selectedTask
}

func main() {
	tasks := []Task{{
		Description: "1. Connect to instance",
		Command:     instance_connect.ConnectInstance,
	}}

	selectedTask := selectTask(tasks)
	selectedTask.Command()
}
