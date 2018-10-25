package helpers

import (
	"encoding/base64"
	"fmt"
	"github.com/BlockClusterApp/daemon/src/dtos"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/getsentry/raven-go"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func getAWSSession(awsCreds *dtos.AWSCreds) *session.Session {
	creds := credentials.NewStaticCredentials(awsCreds.AccessKeyID, awsCreds.SecretAccessKey, "")

	var awsConfig = &aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: creds,
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		raven.CaptureError(err, map[string]string{
			"AccessKeyID": awsCreds.AccessKeyID,
		})
		GetLogger().Printf("Error creating aws session %s", err.Error())
		return nil
	}
	return session.Must(sess, err)
}

func GetAuthorizationToken() string {

	bc := GetBlockclusterInstance()
	bcAwsCreds := bc.GetAWSCreds()

	var awsCreds = &dtos.AWSCreds{
		AccessKeyID:     bcAwsCreds.AccessKeys.AccessKeyId,
		SecretAccessKey: bcAwsCreds.AccessKeys.SecretAccessKey,
	}

	client := ecr.New(getAWSSession(awsCreds))

	var registryIds []*string

	for _, i := range bcAwsCreds.RegistryIds {
		registryIds = append(registryIds, aws.String(i))
	}

	var params = &ecr.GetAuthorizationTokenInput{}
	params.SetRegistryIds(registryIds)

	output, err := client.GetAuthorizationToken(params)

	password := *output.AuthorizationData[0].AuthorizationToken
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "AWS", password)))
	repositories := "402432300121.dkr.ecr.us-west-2.amazonaws.com"
	email := fmt.Sprintf("%s@enterprise.blockcluster.io", bcAwsCreds.ClientID)
	dockerConfig := fmt.Sprintf("{\"auths\":{\"%s\": {\"username\": \"AWS\", \"password\": \"%s\", \"email\": \"%s\", \"auth\":\"%s\"}}}", repositories, password, email, auth)

	if err != nil {
		raven.CaptureError(err, map[string]string{
			"AccessKeyID": awsCreds.AccessKeyID,
		})
		GetLogger().Printf("Error getting authentication token from aws %s", err.Error())
		return ""
	}

	return base64.StdEncoding.EncodeToString([]byte(dockerConfig))
}
