package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"runtime"
	"sync"
	"flag"
)

var (
	prefixFlag = flag.String("p", "", "prefix for search")
	bucketFlag = flag.String("b", "", "bucket for search")
	regionFlag = flag.String("r", "", "region for search")
	endpointFlag = flag.String("e", "", "endpoint for bucket")
	accessFlag = flag.String("a", "", "access key")
	secretFlag = flag.String("s", "", "secret key")
	cpu      = flag.Int("cpus", runtime.NumCPU(), "Number of CPUs to use. Defaults to number of processors.")
	sequentialFlag = flag.Bool("seq", false, "Run queries sequentially when true")
)

func listS3(wg *sync.WaitGroup, index int, results chan<- []string) {

	var creds *credentials.Credentials
	if *accessFlag != "" && *secretFlag != "" {
		creds = credentials.NewStaticCredentials(*accessFlag, *secretFlag, "")
	} else {
		creds = credentials.AnonymousCredentials
	}
	sess := session.New(aws.NewConfig().WithCredentials(creds).WithRegion(*regionFlag).WithEndpoint(*endpointFlag).WithS3ForcePathStyle(true))

	prefix := fmt.Sprintf("%s%x", *prefixFlag, index)

	svc := s3.New(sess)
	inputparams := &s3.ListObjectsInput{
		Bucket: aws.String(*bucketFlag),
		Prefix: aws.String(prefix),
	}

	result := make([]string, 0, 1000)

	svc.ListObjectsPages(inputparams, func(page *s3.ListObjectsOutput, lastPage bool) bool {

		for _, value := range page.Contents {
			result = append(result, *value.Key)
		}

		if lastPage {
			results <- result
			wg.Done()
			return false
		} else {
			return true
		}
	})
}

func listPrefixes() (map[string]bool, error) {

	var wg sync.WaitGroup
	var results = make(chan []string)

	for i := 0x0; i <= 0xf; i++ {
		wg.Add(1)

		go func(index int) {
			listS3(&wg, index, results)
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	prefixHash := make(map[string]bool)
	for result := range results {
		for _, r := range result {
			prefixHash[r] = true
		}
	}

	return prefixHash, nil
}

func listPrefixesSequentially() (map[string]bool, error) {

	var wg sync.WaitGroup
	var results = make(chan []string)
	prefixHash := make(map[string]bool)

	for i := 0x0; i <= 0xf; i++ {
		wg.Add(1)

		go func(index int) {
			listS3(&wg, index, results)
		}(i)

		result := <-results
		for _, r := range result {
			prefixHash[r] = true
		}
	}

	return prefixHash, nil
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*cpu)
	if *prefixFlag == "" || *bucketFlag == "" || *regionFlag == "" || *endpointFlag == "" {
		fmt.Println("Bad arguments")
		return
	}

	var list map[string]bool
	if *sequentialFlag {
		list, _ = listPrefixesSequentially()
	} else {
		list, _ = listPrefixes()
	}

	fmt.Println("Number of replies:", len(list))
}
