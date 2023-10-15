package ec2

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/khaitranhq/aws-works/internal/util"
)

type Instance struct {
	InstanceName *string
	InstanceId   *string
	PublicIp     *string
	State        *string
}

func GetInstances(profile, region string) []Instance {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithDefaultRegion(region),
	)

	if err != nil {
		util.ErrorPrint(err.Error())
		os.Exit(1)
	}

	client := ec2.NewFromConfig(cfg)

	result, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		util.ErrorPrint(err.Error())
		os.Exit(1)
	}

	instances := []Instance{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			minialInstanceData := Instance{
				InstanceName: new(string),
				InstanceId:   instance.InstanceId,
				PublicIp:     instance.PublicIpAddress,
				State:        (*string)(&instance.State.Name),
			}
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					minialInstanceData.InstanceName = tag.Value
				}
			}

			instances = append(instances, minialInstanceData)
		}
	}

	return instances
}

func GetRunningInstances(profile, region string) []Instance {
	instances := GetInstances(profile, region)

	runningInstances := []Instance{}
	for _, instance := range instances {
		if *instance.State == "running" {
			runningInstances = append(runningInstances, instance)
		}
	}
	return runningInstances
}
