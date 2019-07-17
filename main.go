package main

import (
	"errors"
	"fmt"
	"io"
	//"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"html/template"
	"os"
)

type File struct {
	Name string
	Size uint64
}

type Files struct {
	Domain string
	Files  []File
}

func handler() (string, error) {
	bucketName, ok := os.LookupEnv("AWS_BUCKET")
	if ok == false {
		return "AWS_BUCKET is not set", errors.New("AWS_BUCKET is not set")
	}
	domainName, ok := os.LookupEnv("DOMAIN_NAME")
	if ok == false {
		return "DOMAIN_NAME is not set", errors.New("DOMAIN_NAME is not set")
	}
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		return "Unable to list items in bucket", err
	}
	var w io.Writer
	d := Files{domainName, []File{}}
	tmpl := template.New("index")
	tmpl, err = template.ParseFiles("./index.html")
	if err != nil{
		return "template.ParseFiles error", err
	}
	for _, item := range resp.Contents {
		f := File {string(*item.Key), uint64(*item.Size)}
		f.Name = string(*item.Key)
		f.Size = uint64(*item.Size)
		d.Files = append(d.Files, f)
	}
	fmt.Println(d)
	err = tmpl.Execute(w, d)
	if err != nil{
		return "template rendering error", err
	}
	fmt.Println(w)
	return "Hello ƛ!", nil
}

func main() {
	//lambda.Start(handler)
	out, err := handler()
	if err != nil {
		fmt.Println(out)
		fmt.Println(err)
		os.Exit(1)
	}
}
