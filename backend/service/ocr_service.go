package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"racha-historico/domain"
)

type OCRService struct{}

func NewOCRService() *OCRService {
	return &OCRService{}
}

type visionRequestBody struct {
	Requests []visionItem `json:"requests"`
}

type visionItem struct {
	Image    visionImage     `json:"image"`
	Features []visionFeature `json:"features"`
}

type visionImage struct {
	Content string `json:"content"`
}

type visionFeature struct {
	Type       string `json:"type"`
	MaxResults int    `json:"maxResults"`
}

type visionResponse struct {
	Responses []struct {
		FullTextAnnotation struct {
			Text string `json:"text"`
		} `json:"fullTextAnnotation"`
	} `json:"responses"`
}

func (s *OCRService) ExtractFromImage(imageBytes []byte) (*domain.OCRResult, error) {
	apiKey := os.Getenv("GOOGLE_VISION_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GOOGLE_VISION_API_KEY não configurada")
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	requestBody := visionRequestBody{
		Requests: []visionItem{
			{
				Image: visionImage{Content: imageBase64},
				Features: []visionFeature{
					{Type: "TEXT_DETECTION", MaxResults: 1},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://vision.googleapis.com/v1/images:annotate?key=%s", apiKey)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var visionResp visionResponse
	err = json.Unmarshal(responseBytes, &visionResp)
	if err != nil {
		return nil, err
	}

	if len(visionResp.Responses) == 0 || visionResp.Responses[0].FullTextAnnotation.Text == "" {
		return nil, errors.New("nenhum texto encontrado na imagem")
	}

	rawText := visionResp.Responses[0].FullTextAnnotation.Text

	amount, _ := extractAmount(rawText)
	date := extractDate(rawText)
	description := extractDescription(rawText)

	return &domain.OCRResult{
		Amount:      amount,
		Description: description,
		Date:        date,
		RawText:     rawText,
	}, nil
}

func extractAmount(text string) (float64, error) {
	re := regexp.MustCompile(`R\$\s*(\d{1,3}(?:\.\d{3})*(?:,\d{2})?)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return 0, errors.New("valor não encontrado")
	}
	valueStr := strings.ReplaceAll(matches[1], ".", "")
	valueStr = strings.ReplaceAll(valueStr, ",", ".")
	return strconv.ParseFloat(valueStr, 64)
}

func extractDate(text string) string {
	re := regexp.MustCompile(`(\d{2}/\d{2}/\d{4})`)
	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return time.Now().Format("2006-01-02")
	}
	parts := strings.Split(matches[1], "/")
	return fmt.Sprintf("%s-%s-%s", parts[2], parts[1], parts[0])
}

func extractDescription(text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) > 0 && lines[0] != "" {
		return strings.ToUpper(strings.TrimSpace(lines[0]))
	}
	return ""
}
