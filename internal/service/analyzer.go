package service

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"auditor.z9fr.xyz/server/internal/lib"
	"auditor.z9fr.xyz/server/internal/proto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type MatchResult struct {
	ID       int
	File     string
	FileType string
	Search   string
	MatchStr string
	Hits     string
	Line     int
}

type AnalyzerServiceImpl struct {
	log  lib.Logger
	env  *lib.Env
	sess *session.Session
	s3   *s3.S3
}

const (
	TEMP_DIRERCTORY_PREFIX = "auditor-"
)

func NewAnalyzerServiceImpl(
	log lib.Logger,
	env *lib.Env,
) *AnalyzerServiceImpl {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(env.AWS_REGION)},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &AnalyzerServiceImpl{
		log:  log,
		env:  env,
		sess: sess,
		s3:   s3.New(sess),
	}
}

func (s *AnalyzerServiceImpl) GenerateWorkerDir() string {
	dir, err := ioutil.TempDir("", TEMP_DIRERCTORY_PREFIX)

	if err != nil {
		s.log.Fatal(err)
	}

	return dir
}

func (s *AnalyzerServiceImpl) DownloadAndSetupSignatureFilesForAnalysis(signatureName string, workingDir string) (bool, error) {
	resp, err := s.DownloadFileFromS3(signatureName)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	fileDetails, err := os.Create(filepath.Join(workingDir, "signature.yaml"))
	if err != nil {
		return false, err
	}
	defer fileDetails.Close()

	_, err = io.Copy(fileDetails, resp.Body)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *AnalyzerServiceImpl) DownloadFileFromS3(fileName string) (*http.Response, error) {
	req, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.env.S3_BUCKET_NAME),
		Key:    aws.String(fileName),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return nil, err
	}

	resp, err := http.Get(urlStr)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AnalyzerServiceImpl) DownloadAndSetupArchiveForAnalysis(fileName string, workingDir string) (bool, error) {
	resp, err := s.DownloadFileFromS3(fileName)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	buffer, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	readerAt := bytes.NewReader(buffer)
	zipReader, err := zip.NewReader(readerAt, int64(len(buffer)))

	for _, zipFile := range zipReader.File {
		f, err := os.OpenFile(filepath.Join(workingDir, zipFile.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())

		if err != nil {
			return false, err
		}
		defer f.Close()

		rc, err := zipFile.Open()

		if err != nil {
			return false, err
		}

		defer rc.Close()

		_, err = io.Copy(f, rc)

		if err != nil {
			s.log.Error(err)
			return false, err
		}
	}

	if err != nil {
		s.log.Error(err)
		return false, err
	}

	return true, err

}

func (s *AnalyzerServiceImpl) InitiateAnalyzer(req *proto.AuditStartRequest) (string, error) {
	fmt.Print("Analyze start ")

	tmpDir := s.GenerateWorkerDir()

	s.log.Info(tmpDir)
	ok, err := s.DownloadAndSetupSignatureFilesForAnalysis(req.Signature, tmpDir)

	if !ok {
		return "", err
	}

	ok, err = s.DownloadAndSetupArchiveForAnalysis(req.FileName, tmpDir)

	if !ok {
		return "", err
	}

	cmd := exec.Command(s.env.AUDITOR_INSTALL_NAME, "-p", tmpDir, "-s", fmt.Sprintf("%s/signature.yaml", tmpDir))
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var results []MatchResult
	err = json.Unmarshal(out, &results)

	if err != nil {
		panic(err)
	}

	s.log.Info(results)

	return "ok", nil
}
