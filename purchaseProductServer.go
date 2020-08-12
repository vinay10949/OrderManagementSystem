package main

import (
	addtocart "OrderManagementSystem/addToCart"
	"OrderManagementSystem/availablequantity"
	"OrderManagementSystem/purchase"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (*server) PurchaseProduct(ctx context.Context, req *purchase.PurchaseRequest) (*purchase.PurchaseResponse, error) {
	c1:=addtocart.NewAddToCartServiceClient(addToCartService)
	getItem := &addtocart.GetProductFromUserCartRequest{UserNo: req.UserNo}
	res, _ := c1.GetProductFromUserCart(context.Background(), getItem)
	purchaseResponse := &purchase.PurchaseResponse{}
	if res.Name!="" && res.Qty>0{
		fmt.Println(fmt.Sprintf(" Proceeding UserNo : %d ProductName : %s Qty :  %d ",req.UserNo,res.Name,res.Qty))

		c := availablequantity.NewCheckProductAvailableServiceClient(productAvailablityService)
		updateQty := &availablequantity.UpdateProductQuantityRequest{Pr: &availablequantity.Product{Product: res.Name },Qty: res.Qty}
		res, _ := c.UpdateProductAvailable(context.Background(), updateQty)

		if res.Success{
			purchaseResponse.Success = true
			purchaseResponse.ErrMessage = ""
		} else {
			purchaseResponse.Success = false
			purchaseResponse.ErrMessage = "Quantity Exhausted"
		}
	}
	return purchaseResponse, nil
}

var productAvailablityService *grpc.ClientConn

var addToCartService *grpc.ClientConn

func main() {
	var err error
	productAvailablityService, err = grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	addToCartService, err = grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer productAvailablityService.Close()
	defer addToCartService.Close()
	fmt.Println("Purchase Server Running 8083 ")

	lis, err := net.Listen("tcp", "localhost:8083")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)

	}
	s := grpc.NewServer()
	purchase.RegisterPurchaseServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v", err)
	}

}
