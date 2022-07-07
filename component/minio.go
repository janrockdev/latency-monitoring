package component

import (
	"context"
	"flag"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io/ioutil"
	"lmon/common"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func getFilesfromFolder(folder string) []string {
	var fileList []string
	items, _ := ioutil.ReadDir(folder)
	for _, item := range items {
		if item.IsDir() {
			subitems, _ := ioutil.ReadDir(item.Name())
			for _, subitem := range subitems {
				if !subitem.IsDir() {
					fmt.Println(item.Name() + "/" + subitem.Name())
				}
			}
		} else {
			// handle file there
			if string(item.Name()[0]) != "." {
				fileList = append(fileList, folder+"/"+item.Name())
			}
		}
	}
	return fileList
}

func startMinio() {
	ctx := context.Background()
	endpoint := "127.0.0.1:9000"
	accessKeyID := "user"
	secretAccessKey := "xixook2c"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false, //useSSL
	})
	if err != nil {
		common.Logr.Fatalln(err)
	}

	bucketName := "test1"
	location := "eu-west-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			common.Logr.Warnf("Bucket %v already exists!", bucketName)
		} else {
			common.Logr.Fatalln(err)
		}
	} else {
		common.Logr.Infof("Successfully created bucket: %s", bucketName)
	}

	// Upload the zip file
	objectName := "doc1.pdf"
	contentType := "application/pdf"

	// Upload the zip file with FPutObject
	startTotal := time.Now().UnixMilli()
	for _, filePath := range getFilesfromFolder("/Users/janrock/GolandProjects/lmon/source") {
		start := time.Now().UnixMilli()
		info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
		stop := time.Now().UnixMilli()
		if err != nil {
			common.Logr.Fatalln(err)
		}
		common.Logr.Infof("Successfully uploaded %s of size %d. Result time: %v ms.", objectName, info.Size, stop-start)
	}
	stopTotal := time.Now().UnixMilli()
	common.Logr.Infof("Successfully uploaded <s> of size <d>. Result time: %v ms.", stopTotal-startTotal)
}

func run(
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
) error {
	startMinio()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
	}()

	return nil
}

func RunMinio() error {
	common.Logr.Infof("Starting Min.IO...")
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	if err := run(cancel, &wg); err != nil {
		return fmt.Errorf("Error when starting task: %v", err)
	}

	waitSig(ctx, cancel)
	wg.Wait()

	return nil
}

func waitSig(ctx context.Context, cancel func()) {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	select {
	case sig := <-gracefulStop:
		common.Logr.Debugf("Caught signal name=%v", sig)
		common.Logr.Debugf("Closing client connections")
		cancel()
	case <-ctx.Done():
		return
	}
}
