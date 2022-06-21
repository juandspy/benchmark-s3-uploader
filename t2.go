package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	//"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func writeColumnNames(writer *csv.Writer, colNames []string) error {
	err := writer.Write(colNames)
	if err != nil {
		fmt.Println(err, "Write column names to CSV")
		return err
	}
	return nil
}

const MAX = 10000000

func writeTableContent(writer *csv.Writer, tableName string, colNames []string) error {
	// now we know column types, time to perform export
	columns := []string{"val1", "val2", "val3"}

	for i := 0; i < MAX; i++ {
		err := writer.Write(columns)
		if err != nil {
			fmt.Println(err, "write one row")
			return err
		}
	}
	return nil
}

func StoreTable(ctx context.Context,
	minioClient *minio.Client, bucketName string, tableName string) error {
	colNames := []string{"foo", "bar", "baz"}

	buffer := new(bytes.Buffer)

	// initialize CSV writer
	writer := csv.NewWriter(buffer)

	err := writeColumnNames(writer, colNames)
	if err != nil {
		return err
	}

	err = writeTableContent(writer, tableName, colNames)
	if err != nil {
		return err
	}

	writer.Flush()

	reader := io.Reader(buffer)

	size := buffer.Len()

	options := minio.PutObjectOptions{ContentType: "text/csv"}
	objectName := string(tableName) + ".csv"
	_, err = minioClient.PutObject(ctx, bucketName, objectName, reader, int64(size), options)
	if err != nil {
		return err
	}

	//buffer.Reset()
	return nil
}

func NewS3Connection() (*minio.Client, context.Context, error) {
	var endpoint string = "127.0.0.1:9000"
	ctx := context.Background()

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			"minio",
			"minio123", ""),
		Secure: false,
	})

	// check if client has been constructed properly
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return minioClient, ctx, nil
}

func main() {
	minioClient, context, err := NewS3Connection()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	total_csvs, err := strconv.ParseInt(
		os.Getenv("TOTAL_CSVS"), 10, 0,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var i int64
	for i = 0; i < total_csvs; i++ {
		err = StoreTable(context, minioClient, "test", "testX")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
