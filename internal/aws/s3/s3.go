package s3

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/khaitranhq/aws-works/internal/util"
)

type Bucket string

func GetBuckets(profile string) []Bucket {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))

	if err != nil {
		util.ErrorPrint(err.Error())
		os.Exit(1)
	}

	client := s3.NewFromConfig(cfg)

	result, err := client
}
