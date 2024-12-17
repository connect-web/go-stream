package main

import (
	"context"
	"log"
	"net"

	pb "github.com/connect-web/go-stream/api" // Update with your import path
	"github.com/connect-web/go-stream/internal/queue"
	"github.com/connect-web/go-stream/internal/scraper"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedScraperServiceServer
	taskQueue *queue.TaskQueue
}

func (s *Server) SubmitTask(ctx context.Context, req *pb.ScrapeRequest) (*pb.ScrapeResponse, error) {
	log.Printf("Received task: %s", req.Url)

	// Enqueue task (simulate processing in the worker)
	s.taskQueue.Enqueue(req)

	// Simulated response
	response := &pb.ScrapeResponse{
		ResponseBody:    []byte("Task accepted for " + req.Url),
		StatusCode:      202,
		ResponseHeaders: map[string]string{"Content-Type": "application/json"},
		ResponseCookies: map[string]string{"session_id": "placeholder_session"},
	}

	// Process immediately in the placeholder
	scraper.Process(req.Url, req.Headers, req.Cookies)

	return response, nil
}

func main() {
	// Initialize the task queue
	taskQueue := queue.NewTaskQueue(10)
	taskQueue.StartWorker()

	// Setup GRPC server
	server := grpc.NewServer()
	pb.RegisterScraperServiceServer(server, &Server{taskQueue: taskQueue})

	// Start listening
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Server is running on port 50051...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// Wait for all tasks before exiting
	taskQueue.Wait()
}
