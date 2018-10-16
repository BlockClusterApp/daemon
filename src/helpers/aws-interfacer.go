package helpers

import (
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func getAWSSession(awsCreds *dtos.AWSCreds) *session.Session {
	creds := credentials.NewStaticCredentials(awsCreds.AccessKeyID, awsCreds.SecretAccessKey, "")

	var awsConfig = &aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: creds,
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		GetLogger().Printf("Error creating aws session %s", err.Error())
		return nil
	}
	return session.Must(sess, err)
}

func GetAuthorizationToken() string{

	bc := GetBlockclusterInstance()
	bcAwsCreds := bc.GetAWSCreds()

	var awsCreds = &dtos.AWSCreds{
		AccessKeyID: bcAwsCreds.AccessKeys.AccessKeyId,
		SecretAccessKey: bcAwsCreds.AccessKeys.SecretAccessKey,
	}

	client := ecr.New(getAWSSession(awsCreds))

	var params = &ecr.GetAuthorizationTokenInput{}
	var clientId string
	params.SetRegistryIds([]*string{aws.String(fmt.Sprintf("402432300121.dkr.ecr.us-west-2.amazonaws.com/%s-webapp", clientId))})

	output, err := client.GetAuthorizationToken(params)

	if err != nil {
		GetLogger().Printf("Error getting authentication token from aws %s", err.Error())
		return ""
	}

	return *output.AuthorizationData[0].AuthorizationToken
}
