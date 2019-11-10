package awsapp

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/antik9/aws-web-app/internal/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var ImageTemplate string = `
<!DOCTYPE html>
<html lang="en"><head></head>
    <body><img src="data:image/jpg;base64,{{.Image}}"></body>
</html>
`

func GetPicture(w http.ResponseWriter, r *http.Request) {
	defer db.SaveIpToDatabase(r.Header.Get("X-FORWARDED-FOR"))
	sess := session.Must(
		session.NewSession(
			&aws.Config{
				Region: aws.String(endpoints.EuNorth1RegionID),
			},
		),
	)
	svc := s3.New(sess)

	result, err := svc.GetObjectWithContext(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("otus-test"),
		Key:    aws.String("NGC6543.jpg"),
	})
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	defer result.Body.Close()

	fullLength := *result.ContentLength
	var buckets int64
	buffer := make([]byte, fullLength)

	for buckets < fullLength {
		n, err := result.Body.Read(buffer[buckets:])
		if err != nil {
			break
		}
		buckets += int64(n)
	}

	original_image, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	writeImageWithTemplate(w, &original_image)
}

func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}