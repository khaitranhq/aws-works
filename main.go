package main

import (
	server_connect "github.com/khaitranhq/aws-works/internal/tasks/server-connect"
	"github.com/khaitranhq/survey"
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
		Description: "1. Connect to a server",
		Command:     server_connect.ConnectServerTask,
	}}

	selectedTask := selectTask(tasks)
	selectedTask.Command()
}
