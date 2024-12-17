package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pb "go-api-scraper/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type ScraperServer struct {
	pb.UnimplementedScraperServer
}

func NewScraperServer() *ScraperServer {
	return &ScraperServer{}
}

func (s *ScraperServer) ScrapeJson(ctx context.Context, req *pb.ScrapeJsonRequest) (*pb.ScrapeJsonResponse, error) {
	// Fetch the JSON data
	jsonBytes, err := fetchJSON(req.Url, req.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON: %v", err)
	}

	// Parse the JSON into a map[string]interface{}
	var jsonData interface{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Convert to google.protobuf.Struct
	structData, err := structpb.NewStruct(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to Struct: %v", err)
	}

	// Return the response
	return &pb.ScrapeJsonResponse{
		Url:  req.Url,
		Data: structData,
	}, nil
}

// fetchJSON makes an HTTP GET request to fetch JSON data
func fetchJSON(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
