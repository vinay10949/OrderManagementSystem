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
	Data   map[string]int32

}
var s Store
type server struct{}
var mutex sync.RWMutex

func (*server) CheckProductAvailable(ctx context.Context, req *availablequantity.CheckProductAvailableRequest) (*availablequantity.CheckProductAvailableResponse, error) {
	fmt.Println("Checking for Product availablity : ", req.Pr.Product)
	fmt.Println("Current Quantity available :",qty)
	var productQuantity int32=0
	resp := availablequantity.CheckProductAvailableResponse{
	}
	mutex.RLock()
	value,ok:=s.Data[req.Pr.Product] ;if ok{
		productQuantity=value
	}
	if productQuantity<0{
		resp.Qty=0
		mutex.Unlock()
		return &resp,nil
	}
	resp.Qty=productQuantity
	mutex.RUnlock()
	return &resp, nil
}
func (*server) UpdateProductAvailable(ctx context.Context, req *availablequantity.UpdateProductQuantityRequest) (*availablequantity.UpdateProductQuantityResponse, error) {
	fmt.Println("Updating available items in stock  : ", req.Pr.Product)
	q,_:=s.Data[req.Pr.Product]
	fmt.Println("Current Quantity available :",q)
	res:=availablequantity.UpdateProductQuantityResponse{}
	mutex.Lock()
	if s.Data["X"]<req.Qty{
		excess:=req.Qty-q
		newQty:=req.Qty-excess
		s.Data["X"]=q-newQty

	}else{
		fmt.Println("Updating ",s.Data["X"]-req.Qty)
		s.Data["X"]=s.Data["X"]-req.Qty
	}
	mutex.Unlock()
	return &res, nil
}




func main() {

	mutex=sync.RWMutex{}
	s.Data=make(map[string]int32)
	//s.Data=sync.Map{}
	s.Data["X"]=qty
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
