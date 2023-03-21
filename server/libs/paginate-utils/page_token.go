package paginateutils

import (
	"encoding/base64"
	"encoding/json"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.New("Libs.Pagination")

type PageToken struct {
	PageSize int   `json:"pageSize"`
	Offset   int64 `json:"offset"`
}

func (p *PageToken) Encode() (string, error) {
	logger.Info("Start Encode PageToken")

	tokenJSON, err := json.Marshal(p)
	if err != nil {
		logger.ErrorData("Finish Encode PageToken: FAILED", logging.Data{"error": err.Error()})
		return "", err
	}
	encodedToken := base64.StdEncoding.EncodeToString(tokenJSON)

	logger.Info("Finish Encode PageToken: SUCESSFUL")
	return encodedToken, nil
}

func (p *PageToken) GetEncodedNextPageToken() (string, error) {
	nextPageToken := PageToken{
		PageSize: p.PageSize,
		Offset:   p.Offset + int64(p.PageSize),
	}

	return nextPageToken.Encode()
}

func ValidateAndGetPageToken(tokenString string, pageSize int32) (*PageToken, error) {
	logger.Info("Start ValidateAndGetPageToken")

	if tokenString != "" {
		return DecodePageToken(tokenString)
	}

	logger.Info("Finish ValidateAndGetPageToken: SUCESSFUL")
	return &PageToken{
		PageSize: int(pageSize),
		Offset:   0,
	}, nil
}

func DecodePageToken(tokenStr string) (*PageToken, error) {
	logger.Info("Start DecodePageToken")

	var token PageToken
	decodedTokenStr, err := base64.StdEncoding.DecodeString(tokenStr)
	if err != nil {
		logger.ErrorData("Finish DecodePageToken: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}
	err = json.Unmarshal([]byte(decodedTokenStr), &token)

	logger.Info("Finish DecodePageToken: SUCCESSFUL")
	return &token, err
}
