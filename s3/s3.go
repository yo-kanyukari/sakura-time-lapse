package s3

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Getjpgtar 引数で指定したfile pathでS3からタイムラプス用の画像を取得
func Getjpgtar(filePath string, filename string) {
	fmt.Println(filePath)
	//ACCESS_KEYとSECRET_KEYを.envから読む

	creds := credentials.NewStaticCredentials(os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), "")

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("S3_REGION")),
	}))

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to create file %q, %v", filename, err))
		return
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(filePath),
	})
	if err != nil {
		fmt.Println(fmt.Errorf("failed to download file, %v", err))
		return
	}
	fmt.Printf("jpg-file downloaded, %d bytes\n", n)

}

//S3TarPath S3で取得するfileのpathを返す
func S3TarPath(t time.Time) (path string) {
	min := 0
	hour := 0
	if t.Minute() < 10 {
		min = 5
		hour = t.Hour() - 1
	} else {
		min = (t.Minute() / 10) - 1
		hour = t.Hour()
	}
	return fmt.Sprintf("takumi/jpg/%d/%02d%02d%02d/%02d%1d.tar", t.Year(), t.Year()%100, t.Month(), t.Day(), hour, min)

}
