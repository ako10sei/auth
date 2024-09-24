package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	desc "github.com/ako10sei/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserV1Server
	mu    sync.Mutex
	users map[int64]*desc.User
}

func (s *server) List(_ context.Context, req *desc.ListRequest) (*desc.ListResponse, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	l := req.GetLimit()
	if l == 0 {
		return nil, fmt.Errorf("required parameter is missing: limit")
	}

	o := req.GetOffset()
	if o == 0 {
		return nil, fmt.Errorf("required parameter is missing: offset")
	}

	if _, exists := s.users[l]; !exists {
		return nil, fmt.Errorf("user with id %v does not exist", l)
	}
	if _, exists := s.users[o]; !exists {
		return nil, fmt.Errorf("user with id %v does not exist", o)
	}

	users := make([]*desc.User, 0, l)

	for i := l; i <= o; i++ {
		users = append(users, s.users[i])
	}

	return &desc.ListResponse{
		Users: users,
	}, nil
}

// Create method
func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Генерация нового ID для пользователя
	newID := int64(len(s.users) + 1)
	user := &desc.User{
		Id: newID,
		Info: &desc.UserInfo{
			Name:            req.GetInfo().GetName(),
			Email:           req.GetInfo().GetEmail(),
			Password:        req.GetInfo().GetPassword(),
			PasswordConfirm: req.GetInfo().GetPasswordConfirm(),
			Enum:            req.GetInfo().GetEnum(),
		},
	}

	s.users[newID] = user

	log.Printf("Created user: %v", user)
	return &desc.CreateResponse{
		User: user,
	}, nil
}

// Get method (already implemented)
func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.GetId()]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	log.Printf("Retrieved user with ID: %d", req.GetId())
	return &desc.GetResponse{
		User: user,
	}, nil
}

// Update method
func (s *server) Update(_ context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.GetId()]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	if req.GetName() != nil {
		user.Info.Name = req.GetName().GetValue()
	}
	if req.GetEmail() != nil {
		user.Info.Email = req.GetEmail().GetValue()
	}

	log.Printf("Updated user with ID: %d", req.GetId())
	return &emptypb.Empty{}, nil
}

// Delete method
func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.users[req.GetId()]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	delete(s.users, req.GetId())
	log.Printf("Deleted user with ID: %d", req.GetId())
	return &emptypb.Empty{}, nil
}

// Main function
func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Инициализируем сервер и сохраняем пользователей в памяти
	s := grpc.NewServer()
	reflection.Register(s)

	srv := &server{
		users: make(map[int64]*desc.User),
	}
	desc.RegisterUserV1Server(s, srv)

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
