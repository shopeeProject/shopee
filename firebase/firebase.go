package firebaseOp

import (
	"context"
	"io/ioutil"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func FirebaseOp1(c *gin.Context) {
	configStorage := &firebase.Config{
		StorageBucket: "xxxxxxxxxxxxxxxxx.appspot.com",
	}
	opt := option.WithCredentialsFile("C:\\Users\\KARTHIK SURINENI\\Desktop\\shopee\\firebase.json")

	app, err := firebase.NewApp(context.Background(), configStorage, opt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	docName := "abc.txt"

	wc := bucket.Object(docName).NewWriter(context.Background())
	// if _, err = io.Copy(wc, uploadedFile); err != nil {
	// 	log.Error(err)
	// 	return
	// }
	err = wc.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// ******* DOWNLOAD ********
// downloads a file from cloud
func downloadFromCloudStorage(fileURI string) ([]byte, error) {

	configStorage := &firebase.Config{
		StorageBucket: "fxxxxxxxxxxx.appspot.com",
	}
	opt := option.WithCredentialsFile("./config/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), configStorage, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		return nil, err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return nil, err
	}

	rc, err := bucket.Object(fileURI).NewReader(context.Background())
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil

}
