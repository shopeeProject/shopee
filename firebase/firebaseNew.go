package firebaseOp

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/url"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/shopeeProject/shopee/util"
	"google.golang.org/api/option"
)

func returnFailResponse(err string) util.DataResponse {
	return util.DataResponse{
		Success: false,
		Message: "Error while performing Upload " + err,
	}
}

func UploadImageAndGetUrl(file *multipart.FileHeader, fileName string, userEmail string) util.DataResponse {
	// Set up Firebase Admin SDK
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("./../firebase.json"))
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
		return returnFailResponse(err.Error())
	}

	// Create a new storage client
	client, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
		return returnFailResponse(err.Error())
	}

	// Specify the bucket name
	bucketName := "brave-theater-255512.appspot.com"
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Fatalf("Failed to get bucket: %v", err)
		return returnFailResponse(err.Error())
	}

	// Read the image file
	// filePath := fileName
	// file, err := fi
	src, err := file.Open()
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		return returnFailResponse(err.Error())
	}
	defer src.Close()

	// Read file data
	data, err := ioutil.ReadAll(src)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return returnFailResponse(err.Error())
	}

	// Upload the file to Firebase Storage
	objectName := "images/" + fileName + userEmail
	obj := bucket.Object(objectName)
	writer := obj.NewWriter(ctx)
	if _, err := writer.Write(data); err != nil {
		log.Fatalf("Failed to write to bucket: %v", err)
		return returnFailResponse(err.Error())
	}
	if err := writer.Close(); err != nil {
		log.Fatalf("Failed to close writer: %v", err)
		return returnFailResponse(err.Error())
	}

	// Construct the download URL
	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", bucketName, url.QueryEscape(objectName))
	fmt.Printf("File uploaded successfully. Download URL: %s\n", downloadURL)
	return util.DataResponse{
		Success: true,
		Message: "File Uploaded successfully",
		Data:    map[string]string{"downloadURL": downloadURL},
	}
}

// Function to upload file to Firebase Storage
func UploadFile(bucket *storage.BucketHandle, file *multipart.FileHeader, userEmail string) util.DataResponse {
	var ctx = context.Background()
	f, err := file.Open()
	if err != nil {
		return returnFailResponse(err.Error())
	}
	defer f.Close()
	file.Filename = userEmail + file.Filename
	objectName := "shopee/images/" + file.Filename

	// Create a writer to the bucket with the file name
	wc := bucket.Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return returnFailResponse(err.Error())
	}
	if err := wc.Close(); err != nil {
		return returnFailResponse(err.Error())
	}
	bucketName := "brave-theater-255512.appspot.com"
	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", bucketName, url.QueryEscape(objectName))
	fmt.Printf("File uploaded successfully. Download URL: %s\n", downloadURL)
	return util.DataResponse{
		Success: true,
		Message: "File Uploaded successfully",
		Data:    map[string]string{"downloadURL": downloadURL},
	}

}
