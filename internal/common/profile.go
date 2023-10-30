package common

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dlclark/regexp2"
	"github.com/khaitranhq/aws-works/internal/util"
)

const AWS_PROFILE_CONFIG_DIRECTORY = "/.aws/config"

func GetAwsProfiles() []string {
	homeDir, _ := os.UserHomeDir()
	awsProfileConfigData, err := os.ReadFile(homeDir + AWS_PROFILE_CONFIG_DIRECTORY)

	if err != nil {
		util.ErrorPrint(err.Error())
		os.Exit(1)
	}

	profiles := []string{}

	// Check if there is default profile
	re := regexp2.MustCompile("\\[default\\]", 0)
	if isMatch, _ := re.MatchString(string(awsProfileConfigData)); isMatch {
		profiles = append(profiles, "default")
	}

	// Get names of other profiles
	re = regexp2.MustCompile("(?<=profile\\s)[a-zA-Z-]*", 0)
	matches := util.FindAllMatchedSubString(re, string(awsProfileConfigData))
	profiles = append(profiles, matches...)

	return profiles
}

func SelectAwsProfile() string {
	profiles := GetAwsProfiles()
	prompt := &survey.Select{
		Message: "Choose an AWS profile",
		Options: profiles,
	}

	selectedProfile := ""
	survey.AskOne(prompt, &selectedProfile)
	return selectedProfile
}
