package main

import (
	"OrderManagementSystem/availablequantity"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

var qty int32 = 10

type Store struct {
	Data map[string]int32
}

var s Store

type server struct{}

var mutex sync.RWMutex

func (*server) CheckProductAvailable(ctx context.Context, req *availablequantity.CheckProductAvailableRequest) (*availablequantity.CheckProductAvailableResponse, error) {
	fmt.Println("Checking for Product availablity : ", req.Pr.Product)
	var productQuantity int32 = 0
	resp := availablequantity.CheckProductAvailableResponse{}
	mutex.RLock()
	value, ok := s.Data[req.Pr.Product]
	if ok {
		productQuantity = value
	}
	if productQuantity < 0 {
		resp.Qty = 0
		mutex.Unlock()
		return &resp, nil
	}
	resp.Qty = productQuantity
	mutex.RUnlock()
	return &resp, nil
}
func (*server) UpdateProductAvailable(ctx context.Context, req *availablequantity.UpdateProductQuantityRequest) (*availablequantity.UpdateProductQuantityResponse, error) {
	fmt.Println("Updating available items in stock  : ", req.Pr.Product)
	mutex.Lock()
	stock, _ := s.Data[req.Pr.Product]
	res := availablequantity.UpdateProductQuantityResponse{}
	if stock == 0 {
		res.Success = false
		fmt.Println("Item stock over")
	} else if stock < req.Qty {
		res.Success = false
		res.Message = "More quantity demanding as compared to Stock"

	} else {
		fmt.Println("Updating ", s.Data["X"]-req.Qty)
		s.Data["X"] = s.Data[req.Pr.Product] - req.Qty
		res.Success = true
	}
	mutex.Unlock()
	fmt.Println("Final Quantity Left :", s.Data[req.Pr.Product])
	return &res, nil
}

func main() {

	mutex = sync.RWMutex{}
	s.Data = make(map[string]int32)
	//s.Data=sync.Map{}
	s.Data["X"] = qty
	//s.Data.Store("X",qty)
	fmt.Println("Check Availablity MicroService Running 8081")
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)

	}
	s := grpc.NewServer()
	availablequantity.RegisterCheckProductAvailableServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v", err)
	}
}
