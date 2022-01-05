package iam

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/iam"
)

func ListUsers() {
    sess, _ := session.NewSession(
        &aws.Config{
        // Region: aws.String("us-west-2")
    },
    )
    svc := iam.New(sess)
    result, _ := svc.ListUsers(&iam.ListUsersInput{
        MaxItems: aws.Int64(1),
    })
    fmt.Println(result)
	return result
}
