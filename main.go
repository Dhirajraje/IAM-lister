package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type response struct {
	users    []*iam.User
	messgage string
}

func HandleLambdaEvent() {
	sess, err := session.NewSession(&aws.Config{})
	days, err := strconv.Atoi(os.Getenv("DAYS"))
	if err != nil {
		fmt.Println("DAYS not configured correctly, moving forward with 30 days")
		days = 30
	}
	// Handle client error
	if err == nil {
		// Create a IAM service client.
		svciam := iam.New(sess)
		// List users
		result, err := svciam.ListUsers(&iam.ListUsersInput{MaxItems: aws.Int64(1)})

		if err != nil {
			fmt.Println("Error", err)
			return
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
					fmt.Print(err)
					break
				}
			}
			oldUsers := []*iam.User{}
			for _, user := range users {
				timeAgo := time.Now().UTC().Add(time.Duration(days * 24 * int(time.Hour)))
				if user.PasswordLastUsed != nil {
					if timeAgo.Before(*user.PasswordLastUsed) {
						oldUsers = append(oldUsers, user)
					}
				} else {
					if timeAgo.Before(*user.CreateDate) {
						oldUsers = append(oldUsers, user)
					}
				}

			}

			fmt.Print(oldUsers)
			// for i, user := range result.Users {
			// 	if user == nil {
			// 		continue
			// 	}
			// 	fmt.Print(user)
			// 	fmt.Printf("%d user %s created %v\n", i, *user.UserName, user.CreateDate)
			// }
		}

	} else {
		fmt.Print("Error while initilizing session", err)
	}

}

func main() {
	lambda.Start(HandleLambdaEvent)
}
