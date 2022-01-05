package main
import (
    "fmt"
    // "os"

    // "github.com/aws/aws-sdk-go/aws"
    // // "github.com/aws/aws-sdk-go/aws/awserr"
    // "github.com/aws/aws-sdk-go/aws/session"
    // "github.com/aws/aws-sdk-go/service/iam"
    "github.com/Team-Homo-Novus/IAM-lister/aws/iam"
)


func main() {
    result = iam.ListUsers()
    fmt.Println(result)
}