package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const ComfigYaml string = "config.yaml_"

type Config struct {
	Storage struct {
		S3 struct {
			BucketName string `yaml:"bucket" envconfig:"STORAGE_S3_BUCKET"`
			RegionName string `yaml:"region" envconfig:"STORAGE_S3_REGION"`
			KeyID      string `yaml:"key_id" envconfig:"STORAGE_S3_KEY"`
			AccessKey  string `yaml:"access_key" envconfig:"STORAGE_S3_SECRET"`
			ACL        string `yaml:"acl" envconfig:"STORAGE_S3_ACL"`
		}
		Files string `yaml:"files" envconfig:"STORAGE_FILES"`
	}
}

type Reader interface {
	Read() (*Config, error)
}

func main() {
	var configReader Config

	cfgData, err := ioutil.ReadFile(ComfigYaml)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(cfgData, &configReader)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Bucket: ", configReader.Storage.S3.BucketName)
	log.Printf("Bucket: %v", configReader.Storage.S3.BucketName)
	// Create a single AWS session
	sessions, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(configReader.Storage.S3.RegionName),
			Credentials: credentials.NewStaticCredentials(configReader.Storage.S3.KeyID, configReader.Storage.S3.AccessKey, ""),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	// Open file
	testfileupload, err := ioutil.ReadFile(configReader.Storage.Files)
	if err != nil {
		log.Fatal(err)
	}
	// Uploads
	log.Print("Upload file: ")
	_, err = s3.New(sessions).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(configReader.Storage.S3.BucketName),
		Key:                  aws.String(configReader.Storage.Files),
		ACL:                  aws.String(configReader.Storage.S3.ACL),
		Body:                 bytes.NewReader(testfileupload),
		ContentType:          aws.String(http.DetectContentType(testfileupload)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Print("OK")
	log.Print("Check file in bucket: ")
	checktestfile := &s3.GetObjectInput{
		Bucket: aws.String(configReader.Storage.S3.BucketName),
		Key:    aws.String(configReader.Storage.Files),
	}
	_, err = s3.New(sessions).GetObject(checktestfile)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("OK")
	log.Print("Remove file: ")
	deletetestfile := &s3.DeleteObjectInput{
		Bucket: aws.String(configReader.Storage.S3.BucketName),
		Key:    aws.String(configReader.Storage.Files),
	}
	_, err = s3.New(sessions).DeleteObject(deletetestfile)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("OK")
}
