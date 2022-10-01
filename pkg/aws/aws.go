package aws

import (
	"context"
	"log"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/config"
	"github.com/sirupsen/logrus"
)

func resetDriftctlAWSVars() {
	os.Unsetenv("DCTL_S3_ACCESS_KEY_ID")
	os.Unsetenv("DCTL_S3_SECRET_ACCESS_KEY")
	os.Unsetenv("DCTL_S3_SESSION_TOKEN")
	//os.Unsetenv("DCTL_S3_PROFILE")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SESSION_TOKEN")
	//os.Unsetenv("AWS_PROFILE")

}
func SetupAWS(config config.DritfctlRun) {

	resetDriftctlAWSVars()
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		logrus.Panic(err)
	}
	switch config.AWSConfig.TFStateMethod {
	case "assumeRole":
		stsSvcState := sts.NewFromConfig(cfg)
		credsState := stscreds.NewAssumeRoleProvider(stsSvcState, config.AWSConfig.TFStateRole)
		dtclCreds, err := credsState.Retrieve(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		os.Setenv("DCTL_S3_ACCESS_KEY_ID", dtclCreds.AccessKeyID)
		os.Setenv("DCTL_S3_SECRET_ACCESS_KEY", dtclCreds.SecretAccessKey)
		os.Setenv("DCTL_S3_SESSION_TOKEN", dtclCreds.SessionToken)
		logrus.Info("Using assume role to check tf+s3 state")
	case "direct":
		logrus.Info("Using default AWS configuration to check tf+s3 state")
	case "awsprofile":
		os.Setenv("DCTL_S3_ACCESS_KEY_ID", config.AWSConfig.TFStateProfile)
	default:
		logrus.Warnf("tfStateMethod %s is not correct. Using default AWS configuration to check tf+s3 state", config.AWSConfig.TFStateMethod)
	}
	switch config.AWSConfig.InfraScanMethod {
	case "assumeRole":
		stsSvcScan := sts.NewFromConfig(cfg)
		credsScan := stscreds.NewAssumeRoleProvider(stsSvcScan, config.AWSConfig.InfraScanRole)
		awsCreds, err := credsScan.Retrieve(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		os.Setenv("AWS_ACCESS_KEY_ID", awsCreds.AccessKeyID)
		os.Setenv("AWS_SECRET_ACCESS_KEY", awsCreds.SecretAccessKey)
		os.Setenv("AWS_SESSION_TOKEN", awsCreds.SessionToken)
	case "direct":
		logrus.Info("Using default AWS configuration to scan AWS resources")
	default:
		logrus.Warnf("InfraScanMethod %s is not correct. Using default AWS configuration to scan AWS resources", config.AWSConfig.TFStateMethod)
	}
}
