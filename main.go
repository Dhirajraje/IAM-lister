package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

//	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Response struct {
	UserId     string `json:"user_id"`
	User       string `json:"user_name"`
	LastUsedOn string `json:"last_used_on"`
	CreatedOn  string `json:"created_on"`
	Arn        string `json:"user_arn"`
}

func IamLister() string {
	sess, err := session.NewSession(&aws.Config{})
	days, _ := strconv.Atoi(os.Getenv("DAYS"))
	// Handle client error
	if err == nil {
		// Create a IAM service client.
		svciam := iam.New(sess)
		// List users
		result, err := svciam.ListUsers(&iam.ListUsersInput{MaxItems: aws.Int64(1)})
		if err != nil {
			return "Error" + err.Error()
		} else {
			users := result.Users
			isTruncated := *result.IsTruncated
			maker := result.Marker
			for isTruncated {
				result, err := svciam.ListUsers(&iam.ListUsersInput{Marker: maker})
				if err == nil {
					isTruncated = *result.IsTruncated
					maker = result.Marker
					users = append(users, result.Users...)
				} else {
					fmt.Println(err.Error())
					break
				}
			}
			oldUsers := []Response{}
			for _, user := range users {
				timeAgo := time.Now().UTC().Add(time.Duration(-days * 24 * int(time.Hour)))
				if user.PasswordLastUsed != nil {
					if timeAgo.After(*user.PasswordLastUsed) {
						oldUsers = append(oldUsers, Response{
							User:       *user.UserName,
							LastUsedOn: user.PasswordLastUsed.String(),
							CreatedOn:  user.CreateDate.String(),
							Arn:        *user.Arn,
							UserId:     *user.UserId,
						})
					}
				} else {
					if timeAgo.After(*user.CreateDate) {
						oldUsers = append(oldUsers, Response{
							User:       *user.UserName,
							LastUsedOn: "Never used",
							CreatedOn:  user.CreateDate.String(),
							Arn:        *user.Arn,
							UserId:     *user.UserId,
						})
					}
				}
			}
			res, _ := json.MarshalIndent(oldUsers, "", "  ")
			return string(res)
		}
	} else {
		return "Error while initilizing session" + err.Error()
	}
}

func main() {
//	lambda.Start(HandleLambdaEvent)
	fmt.Println(IamLister())
}
