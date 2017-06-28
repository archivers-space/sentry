package archive

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/multiformats/go-multihash"
	"io"
	"net/http"
)

// File is a buffered byte slice often made from a GET response body.
// It provides easy hash-calculation & storage to S3
// TODO - depricate, use s3-datastore, or, uh... the distributed web
type File struct {
	Url  string
	Data []byte
	Hash string
}

func FileUrl(url *Url) string {
	// f := &File{Hash: hash}
	// filename, err := f.Filename()
	// if err != nil {
	// 	return ""
	// }
	// return fmt.Sprintf("https://s3.amazonaws.com/%s/%s/%s", AwsS3BucketName, AwsS3BucketPath, filename)
	return fmt.Sprintf("https://content.archivers.space/urls/%s", url.Id)
}

// NewFileFromRes generates a new file by consuming & closing a given response body
func NewFileFromRes(url string, res *http.Response) (*File, error) {
	f := &File{Url: url}

	buf := &bytes.Buffer{}

	if _, err := io.Copy(buf, res.Body); err != nil {
		fmt.Println("copy error: %s")
		return nil, err
	}
	defer res.Body.Close()

	f.Data = buf.Bytes()

	if err := f.calcHash(); err != nil {
		return nil, fmt.Errorf("err calculating hash for url: %s error: %s", f.Url, err.Error())
	}
	return f, nil
}

// Filename returns the name of the file, which is it's sha2-256 hash
func (f *File) Filename() (string, error) {
	if f.Data == nil && f.Hash == "" {
		return "", fmt.Errorf("no data or hash for filename")
	}

	if f.Hash == "" {
		if err := f.calcHash(); err != nil {
			return "", err
		}
	}

	// lop the multihash prefix off for storage purposes so the files don't all
	// start with "1120", this has the effect of returning valid sha256 hashes
	// as filenames
	return f.Hash[4:], nil
}

// GetS3 reads a given file's hash from S3
func (f *File) GetS3() error {
	filename, err := f.Filename()
	if err != nil {
		return err
	}

	svc := s3.New(session.New(&aws.Config{
		Region:      aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(AwsAccessKeyId, AwsSecretAccessKey, ""),
	}))

	res, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(AwsS3BucketName),
		Key:    aws.String(f.s3Path(filename)),
	})

	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, res.Body); err != nil {
		return err
	}

	f.Data = buf.Bytes()
	return res.Body.Close()
}

// PutS3 puts the file on S3 if it doesn't already exist
func (f *File) PutS3() error {
	if f.Data == nil || len(f.Data) == 0 {
		return fmt.Errorf("no data for saving url to s3: '%s'", f.Url)
	}

	filename, err := f.Filename()
	if err != nil {
		return err
	}

	svc := s3.New(session.New(&aws.Config{
		Region:      aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(AwsAccessKeyId, AwsSecretAccessKey, ""),
	}))

	// check to see if hash exists
	_, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(AwsS3BucketName),
		Key:    aws.String(f.s3Path(filename)),
	})

	// TODO - silly hack to check for 404 response
	// should open error up & look for ErrNotFound
	if err != nil {
		_, err = svc.PutObject(&s3.PutObjectInput{
			ACL:    aws.String(s3.BucketCannedACLPublicRead),
			Bucket: aws.String(AwsS3BucketName),
			Key:    aws.String(f.s3Path(filename)),
			Body:   bytes.NewReader(f.Data),
		})
	}

	return err
}

// calculated path for s3 placement
func (f *File) s3Path(filename string) string {
	return AwsS3BucketPath + "/" + filename
}

// Delete a file from S3
func (f *File) Delete() error {
	filename, err := f.Filename()
	if err != nil {
		return err
	}

	svc := s3.New(session.New(&aws.Config{
		Region:      aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(AwsAccessKeyId, AwsSecretAccessKey, ""),
	}))

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(AwsS3BucketName),
		Key:    aws.String(f.s3Path(filename)),
	})

	return nil
}

// calculate the sha2-256 hash of the file's data
func (f *File) calcHash() error {
	h := sha256.New()
	h.Write(f.Data)

	mhBuf, err := multihash.EncodeName(h.Sum(nil), "sha2-256")
	if err != nil {
		return err
	}

	f.Hash = hex.EncodeToString(mhBuf)
	return nil
}
