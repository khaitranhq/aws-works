package server_connect

import (
	"fmt"
	"os"
	"strings"

	// "os/exec"

	"github.com/khaitranhq/survey"
	"github.com/atotto/clipboard"
	"github.com/khaitranhq/aws-works/internal/aws/ec2"
	"github.com/khaitranhq/aws-works/internal/common"
	"github.com/khaitranhq/aws-works/internal/util"
)

const EC2_KEYS_DIRECTORY = "/aws-works/keys"

func selectInstance(profile, region string) ec2.Instance {

	instances := ec2.GetRunningInstances(profile, region)

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

func getSSHCommand(
	region, keyPairFolder, user, profile, instanceId string,
	publicIp *string,
) string {
	keyPairDirectory := fmt.Sprintf("%s/%s@%s", keyPairFolder, user, instanceId)

	if publicIp == nil {
		sshCommand := fmt.Sprintf(
			"ssh -i %s -o ProxyCommand='aws ec2-instance-connect open-tunnel --instance-id %s --profile %s --region %s' %s@%s",
			keyPairDirectory,
			instanceId,
			profile,
			region,
			user,
			instanceId,
		)
		return sshCommand
	}

	sshCommand := fmt.Sprintf("ssh -i %s %s@%s", keyPairDirectory, user, *publicIp)
	return sshCommand
}

func getConnectMethod() string {
	prompt := &survey.Select{
		Message: "Choose the connection method",
		Options: []string{"AWS System Manager", "SSH"},
		Default: "AWS System Manager",
	}

	selectedMethod := "AWS System Manager"
	survey.AskOne(prompt, &selectedMethod)
	return selectedMethod
}

func ConnectServerTask() {
	serverLocations := []string{"AWS", "Other"}
	selectServerLocationPrompt := &survey.Select{
		Message: "Select the location of server",
		Options: serverLocations,
	}
	selectedServerLocation := "AWS"
	survey.AskOne(selectServerLocationPrompt, &selectedServerLocation)

	if selectedServerLocation == "AWS" {
		profile := common.SelectAwsProfile()
		region := common.SelectRegion(profile)

		// Check existence of instance key
		homeUserDir, _ := os.UserHomeDir()

		keyPairFolder := fmt.Sprintf("%s/.ssh/%s/%s", homeUserDir, profile, region)

		if _, err := os.Stat(keyPairFolder); err != nil {
			// Create keys folder
			err := os.MkdirAll(keyPairFolder, 0755)
			if err != nil {
				util.ErrorPrint(err.Error())
				os.Exit(1)
			}
		}

		instance := selectInstance(profile, region)
		connectionMethod := getConnectMethod()

		if connectionMethod == "SSH" {
			user := selectUser(keyPairFolder, *instance.InstanceId)
			command := getSSHCommand(
				region,
				keyPairFolder,
				user,
				profile,
				*instance.InstanceId,
				instance.PublicIp,
			)
			clipboard.WriteAll(command)
			fmt.Println("Copyied SSH command to clipboard")
		}

		if connectionMethod == "AWS System Manager" {
			connectCommand := fmt.Sprintf(
				"aws ssm start-session --target %s --profile %s --region %s",
				*instance.InstanceId,
				profile,
				region,
			)
			clipboard.WriteAll(connectCommand)
		}
	} else if selectedServerLocation == "Other" {

	}
}
