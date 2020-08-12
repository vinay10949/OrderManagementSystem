package main

import (
	addtocart "OrderManagementSystem/addToCart"
	"OrderManagementSystem/availablequantity"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type server struct{}

var userQty sync.Map

func (*server) AddToCart(ctx context.Context, req *addtocart.AddToCartRequest) (*addtocart.AddToCartResponse, error) {
	fmt.Println("Checking for Product availablity : ", req.Pr.Name)

	c := availablequantity.NewCheckProductAvailableServiceClient(productAvailablityService)
	availableReq := &availablequantity.CheckProductAvailableRequest{Pr: &availablequantity.Product{Product: req.Pr.Name}}
	res, _ := c.CheckProductAvailable(context.Background(), availableReq)
	addToCartResponse := &addtocart.AddToCartResponse{}
	if res.Qty > 0 {
		userQty.Store(req.UserNo, map[string]int32{req.Pr.Name: req.Pr.Qty})
		addToCartResponse.Success = true
		addToCartResponse.ErrMessage = ""
		fmt.Println(fmt.Sprintf("Added Product : %s for User : %d Qty : %d", req.Pr.Name, req.UserNo, req.Pr.Qty))

	} else {
		addToCartResponse.Success = false
		addToCartResponse.ErrMessage = "Quantity Exhausted"
	}
	return addToCartResponse, nil
}

func (*server) GetProductFromUserCart(ctx context.Context, req *addtocart.GetProductFromUserCartRequest) (*addtocart.GetProductFromUserCartResponse, error) {
	val, ok := userQty.Load(req.UserNo)
	res := addtocart.GetProductFromUserCartResponse{}

	if val == nil {
		fmt.Println(" NIL ")
		return &res, nil
	}
	v := val.(map[string]int32)
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	if ok {
		res.Name = keys[0]
		res.Qty = v[keys[0]]

	}

	return &res, nil
}

var productAvailablityService *grpc.ClientConn

func main() {
	var err error
	userQty = sync.Map{}
	productAvailablityService, err = grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer productAvailablityService.Close()
	fmt.Println("Add To Cart Server Running 8082 ")

	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)

	}
	s := grpc.NewServer()
	addtocart.RegisterAddToCartServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v", err)
	}

}
