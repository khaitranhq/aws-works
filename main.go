package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dlclark/regexp2"
)

const AWS_PROFILE_CONFIG_DIRECTORY = "/.aws/config"

func findAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func getAwsProfiles() []string {
	homeDir, _ := os.UserHomeDir()
	awsProfileConfigData, err := os.ReadFile(homeDir + AWS_PROFILE_CONFIG_DIRECTORY)

	if err != nil {
		fmt.Printf("\033[1;31m%s\033[0m", "AWS config file not existed!")
		os.Exit(1)
	}

	profiles := []string{}

	// Check if there is default profile
	re := regexp2.MustCompile("\\[default\\]", 0)
	if isMatch, _ := re.MatchString(string(awsProfileConfigData)); isMatch {
		profiles = append(profiles, "default")
	}

	// Get names of other profiles
	re = regexp2.MustCompile("(?<=profile\\s)[a-zA-Z]*", 0)
	matches := findAllString(re, string(awsProfileConfigData))
	profiles = append(profiles, matches...)

	return profiles
}

func selectAwsProfile() string {
	profiles := getAwsProfiles()
	prompt := &survey.Select{
		Message: "Choose an AWS profile",
		Options: profiles,
	}

	selectedProfile := ""
	survey.AskOne(prompt, &selectedProfile)
	return selectedProfile
}

func selectRegion() string {
	regions := []string{
		"af-south-1",
		"ap-east-1",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-south-1",
		"ap-south-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-southeast-3",
		"ap-southeast-4",
		"ca-central-1",
		"cn-north-1",
		"cn-northwest-1",
		"eu-central-1",
		"eu-central-2",
		"eu-north-1",
		"eu-south-1",
		"eu-south-2",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"il-central-1",
		"me-central-1",
		"me-south-1",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-gov-east-1",
		"us-gov-west-1",
		"us-west-1",
		"us-west-2",
	}

	selectedRegion := ""
	prompt := &survey.Select{
		Message: "Choose a region",
		Options: regions,
	}

	survey.AskOne(prompt, &selectedRegion)
	return selectedRegion
}

type Task struct {
	Describe string
	Command  string
}

func selectTask(profile string, region string) string {
	tasks := []Task{{
		Describe: "1. List instance ids, public IPs, key pairs of EC2 instances",
		Command:  "aws ec2 describe-instances --filters Name=instance-state-name,Values=running --query 'Reservations[*].Instances[*].[Tags[?Key==`Name`].Value,PublicIpAddress,InstanceId]' --output table --profile " + profile,
	}}

	tasksDescribe := []string{}
	for _, task := range tasks {
		tasksDescribe = append(tasksDescribe, task.Describe)
	}

	selectedTask := ""
	prompt := &survey.Select{
		Message: "Choose a task",
		Options: tasksDescribe,
	}
	survey.AskOne(prompt, &selectedTask)
	return selectedTask
}

func main() {

}
