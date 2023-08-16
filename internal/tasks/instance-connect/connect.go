package instance_connect

import (
	"fmt"
	"os"
	"strings"

	// "os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/khaitranhq/aws-works/internal/aws/ec2"
	"github.com/khaitranhq/aws-works/internal/common"
	"github.com/khaitranhq/aws-works/internal/util"
)

const EC2_KEYS_DIRECTORY = "/aws-works/keys"

func selectInstance(profile, region string) ec2.Instance {

	instances := ec2.GetInstances(profile, region)

	// Select instance
	selectInstanceOptions := []string{}
	for _, instance := range instances {
		selectInstanceOptions = append(
			selectInstanceOptions,
			*instance.InstanceId+" "+*instance.InstanceName,
		)
	}

	prompt := &survey.Select{
		Message: "Choose an instance",
		Options: selectInstanceOptions,
	}

	selectedInstanceOption := ""
	survey.AskOne(prompt, &selectedInstanceOption)

	selectedInstance := ec2.Instance{}
	for _, instance := range instances {
		if selectedInstanceOption == *instance.InstanceId+" "+*instance.InstanceName {
			selectedInstance = instance
		}
	}
	return selectedInstance
}

func savePrivateKey(fileDirectory string) {
	prompt := &survey.Multiline{
		Message: "Enter private key of user:",
	}

	privateKey := ""
	survey.AskOne(prompt, &privateKey)

	// Create file
	file, err := os.Create(fileDirectory)
	if err != nil {
		util.ErrorPrint(err.Error())
		os.Exit(1)
	}

	err = os.Chmod(fileDirectory, 0600)
	if err != nil {
		util.ErrorPrint(err.Error())
	}

	// Write file
	_, err = file.WriteString(privateKey)
	if err != nil {
		util.ErrorPrint(err.Error())
		os.Exit(1)
	}

	fmt.Println("Saved SSH key")
}

func selectUser(keyFolderDirectory string, instanceId string) string {
	files, err := os.ReadDir(keyFolderDirectory)

	if err != nil {
		util.ErrorPrint(err.Error())
		os.Exit(1)
	}

	usersOfInstance := []string{}
	for _, file := range files {
		if strings.Contains(file.Name(), instanceId) {
			splittedFileName := strings.Split(file.Name(), "@")
			usersOfInstance = append(usersOfInstance, splittedFileName[0])
		}
	}
	usersOfInstance = append(usersOfInstance, "New user")

	selectUserPrompt := &survey.Select{
		Message: "Select a user to connect",
		Options: usersOfInstance,
	}

	selectedUser := ""
	survey.AskOne(selectUserPrompt, &selectedUser)

	if selectedUser == "New user" {
		newUser := ""
		newUserPrompt := &survey.Input{
			Message: "Enter new user",
		}
		survey.AskOne(newUserPrompt, &newUser)

		savePrivateKey(keyFolderDirectory + "/" + newUser + "@" + instanceId)
		return newUser
	}

	return selectedUser
}

func ConnectInstance() {
	profile := common.SelectAwsProfile()
	region := common.SelectRegion()

	// Check existence of instance key
	configDir, configDirErr := os.UserConfigDir()
	if configDirErr != nil {
		util.ErrorPrint(configDirErr.Error())
		os.Exit(1)
	}

	keyPairFolder := configDir + EC2_KEYS_DIRECTORY + "/" + profile + "/" + region

	_, err := os.Stat(
		keyPairFolder,
	)

	if err != nil {
		// Create keys folder
		err := os.MkdirAll(keyPairFolder, 0755)
		if err != nil {
			util.ErrorPrint(err.Error())
			os.Exit(1)
		}
	}

	instance := selectInstance(profile, region)
	user := selectUser(keyPairFolder, *instance.InstanceId)

	fmt.Println("Default connection type is SSH")

	keyPairDirectory := keyPairFolder + "/" + user + "@" + *instance.InstanceId

	os.Setenv("AWS_WORKS_SSH_DIRECTORY", keyPairDirectory)
	sshCommand := "ssh -i " + keyPairDirectory + " " + user + "@" + *instance.PublicIp + "\n"
	fmt.Printf("SSH Command: %s", sshCommand)
}
