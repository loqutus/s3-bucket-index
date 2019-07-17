package main

import (
	"bytes"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dustin/go-humanize"
	"html/template"
	"os"
)

type File struct {
	Name string
	Size string
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
	var tpl bytes.Buffer
	d := Files{domainName, []File{}}
	tmpl := template.New("index")
	tmpl, err = template.ParseFiles("./index.html")
	if err != nil{
		return "template.ParseFiles error", err
	}
	for _, item := range resp.Contents {
		var f File
		f.Name = string(*item.Key)
		f.Size = humanize.Bytes(uint64(*item.Size))
		d.Files = append(d.Files, f)
	}
	err = tmpl.Execute(&tpl, d)
	if err != nil{
		return "template rendering error", err
	}
	return tpl.String(), nil
}

func main() {
	lambda.Start(handler)
}
