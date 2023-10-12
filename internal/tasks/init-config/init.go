package init_config

import (
	"fmt"

	"github.com/khaitranhq/aws-works/internal/common"
)

func InitConfig() {
	fmt.Println("Select a profile and an S3 bucket")
	profile := common.SelectAwsProfile()
}
