package common

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dlclark/regexp2"
	"github.com/khaitranhq/aws-works/internal/config"
	"github.com/khaitranhq/aws-works/internal/util"
)

var regions = []string{
	"ap-southeast-1",
	"af-south-1",
	"ap-east-1",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-northeast-3",
	"ap-south-1",
	"ap-south-2",
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

func GetDefaultRegion(profile string) string {
	homeDir, _ := os.UserHomeDir()
	awsProfileConfigData, err := os.ReadFile(homeDir + config.AWS_PROFILE_CONFIG_DIRECTORY)

	if err != nil {
		util.ErrorPrint(err.Error())
		os.Exit(1)
	}

	re := regexp2.MustCompile(
		fmt.Sprintf("(?<=\\[profile %s\\]\\nregion = )[a-zA-Z0-9-]*", profile),
		0,
	)
	if match, _ := re.FindStringMatch(string(awsProfileConfigData)); match != nil {
		return match.String()
	}
	return regions[0]
}

func sortRegion(defaultRegion string) []string {
	defaultRegionIndex := 0
	for i, region := range regions {
		if region == defaultRegion {
			defaultRegionIndex = i
			break
		}
	}

	return append(
		[]string{defaultRegion},
		append(regions[:defaultRegionIndex], regions[defaultRegionIndex+1:]...)...,
	)
}

func SelectRegion(profile string) string {
	defaultRegion := GetDefaultRegion(profile)
	fmt.Println(defaultRegion)
	selectedRegion := ""
	sortedRegion := sortRegion(defaultRegion)
	prompt := &survey.Select{
		Message: "Choose a region",
		Options: sortedRegion,
	}

	survey.AskOne(prompt, &selectedRegion)
	return selectedRegion

}
