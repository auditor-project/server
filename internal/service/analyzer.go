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
	"strings"
	"time"

	"auditor.z9fr.xyz/server/internal/db"
	"auditor.z9fr.xyz/server/internal/lib"
	"auditor.z9fr.xyz/server/internal/proto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Results struct {
	ProjectID string        `json:"projectId"`
	Results   []MatchResult `json:"results"`
}

type MatchResult struct {
	ID       int             `json:"id"`
	File     string          `json:"file"`
	Filetype string          `json:"filetype"`
	Search   string          `json:"search"`
	MatchStr string          `json:"match_str"`
	Hits     string          `json:"hits"`
	Line     int             `json:"line"`
	Code     [][]interface{} `json:"code"`
}

type AnalyzerService struct {
	log  lib.Logger
	env  *lib.Env
	db   *db.Database
	sess *session.Session
	s3   *s3.S3
}

const (
	TEMP_DIRERCTORY_PREFIX = "auditor-"
)

func NewAnalyzerServiceImpl(
	log lib.Logger,
	env *lib.Env,
) *AnalyzerService {
	log.Debug("Init analyzier service ")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(env.AWS_REGION)},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &AnalyzerService{
		log:  log,
		env:  env,
		sess: sess,
		s3:   s3.New(sess),
	}
}

func (s *AnalyzerService) GenerateWorkerDir() string {
	dir, err := ioutil.TempDir("", TEMP_DIRERCTORY_PREFIX)

	if err != nil {
		s.log.Fatal(err)
	}

	return dir
}

func (s *AnalyzerService) DownloadAndSetupSignatureFilesForAnalysis(signatureName string, workingDir string) (bool, error) {
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

func (s *AnalyzerService) DownloadFileFromS3(fileName string) (*http.Response, error) {
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

func (s *AnalyzerService) DownloadAndSetupArchiveForAnalysis(fileName string, workingDir string) (bool, error) {
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
		filePath := filepath.Join(workingDir, zipFile.Name)

		if zipFile.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, zipFile.Mode())
			if err != nil {
				return false, err
			}
			continue
		}

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

func (s *AnalyzerService) InitiateAnalyzer(req *proto.AuditStartRequest) (string, error) {
	s.log.Info("Starting to process", "request", req)
	//if err := s.db.Debug().Table("project").Where("id = ?", req.ProjectId).Update("currentStatus", "processing").Error; err != nil {
	//		return "", err
	// 	}

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

	url := fmt.Sprintf("%s/save-results", s.env.NEXT_API_URL)
	method := "POST"

	resultsApiFormat := Results{
		ProjectID: req.ProjectId,
		Results:   results,
	}

	jsonData, err := json.Marshal(resultsApiFormat)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	httpreq, err := http.NewRequest(method, url, strings.NewReader(string(jsonData)))

	if err != nil {
		return "", err
	}

	httpreq.Header.Add("x-api-token", s.env.BULK_SAVE_API_KEY)
	httpreq.Header.Add("Content-Type", "application/json")

	res, err := client.Do(httpreq)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	s.log.Info(body)
	return "ok", err
}
