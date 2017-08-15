package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"gopkg.in/asaskevich/govalidator.v6"
)

func main() {
	grafanaURL := flag.String("grafanaURL", "", "The URL of the grafana server to be backed up")
	s3Bucket := flag.String("s3Bucket", "", "The name of the S3 bucket where the backup should be stored")
	backupPath := flag.String("backupPath", "", "The path of the backup from which to restore (within the S3 bucket)")
	useSharedConfig := flag.Bool("useSharedConfig", false, "Controls whether to use the ~/.aws shared config")

	flag.Parse()

	if *grafanaURL == "" || *s3Bucket == "" || *backupPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if !govalidator.IsURL(*grafanaURL) {
		log.Fatalf("Invalid grafanaURL: %v", *grafanaURL)
	}

	// S3 client
	var sess *session.Session
	if *useSharedConfig {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		sess = session.Must(session.NewSession())
	}
	client := s3.New(sess)

	// Get the list of dashboards from specified bucket/path combination
	resp, err := client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(*s3Bucket),
		Prefix: aws.String(*backupPath),
	})
	if err != nil {
		log.Fatalf("Error retrieving backup dashboard list: %v", err)
	}

	// Restore requested dashboards (all dashboards by default)
	createDashboardURL := *grafanaURL + "/api/dashboards/db"
	for _, dashboard := range resp.Contents {

		// get the dashboard from S3
		dashboardPath := *dashboard.Key
		output, err := client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(*s3Bucket),
			Key:    aws.String(dashboardPath),
		})
		if err != nil {
			log.Fatalf("Error retrieving dashboard with path %v: %v", dashboardPath, err)
		}

		// send it to grafana
		resp, err := http.Post(createDashboardURL, "application/json", output.Body)
		if err != nil {
			log.Fatalf("Error retrieving dashboard with path %v: %v", dashboardPath, err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Received non-OK response while retrieving POSTing with path %v: %v", dashboardPath, resp)
		}
	}
}
