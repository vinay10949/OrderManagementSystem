package main

import (
	addtocart "OrderManagementSystem/addToCart"
	"OrderManagementSystem/availablequantity"
	"OrderManagementSystem/purchase"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	cc, err := grpc.Dial("localhost:8081", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect : %v", err)

	}
	defer cc.Close()

	c1 := availablequantity.NewCheckProductAvailableServiceClient(cc)

	cc1, err := grpc.Dial("localhost:8082", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect : %v", err)

	}
	defer cc1.Close()
	c2 := addtocart.NewAddToCartServiceClient(cc1)

	cc2, err := grpc.Dial("localhost:8083", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect : %v", err)

	}
	defer cc2.Close()
	c3 := purchase.NewPurchaseServiceClient(cc2)
	var i int32
	for i = 0; i < 10	; i++ {
		wg.Add(1)
		go orderProcess(c1, c2,c3, &wg, i)
	}

	wg.Wait()

	
	fmt.Println("Checking for Availablity of item ")
	availableReq := &availablequantity.CheckProductAvailableRequest{Pr: &availablequantity.Product{Product: "X"}}
	res, _ := c1.CheckProductAvailable(context.Background(), availableReq)
	fmt.Println("Item available " ,res.Qty)
	if res.Qty==0{
		fmt.Println("Item exhausted ")
	}

}

func orderProcess(c availablequantity.CheckProductAvailableServiceClient,
	c1 addtocart.AddToCartServiceClient,
	c2 purchase.PurchaseServiceClient,
	wg *sync.WaitGroup, i int32) {
	availableReq := &availablequantity.CheckProductAvailableRequest{Pr: &availablequantity.Product{Product: "X"}}
	res, err := c.CheckProductAvailable(context.Background(), availableReq)
	if err != nil {
		log.Fatalf("Error while calling availablity server %v", err)
	}
	log.Printf("Available Quantity %d USer : %d", res.Qty, i)

	addToCartReq := &addtocart.AddToCartRequest{Pr: &addtocart.Product{Name: "X", Qty: randomQuantity()},UserNo: i}
	res1, err := c1.AddToCart(context.Background(), addToCartReq)
	if err != nil {
		log.Fatalf("Error while calling availablity server %v", err)
	}
	log.Printf("Add to Cart Status %v for User %d ", res1.Success, i)

	purchaseReq := &purchase.PurchaseRequest{UserNo: i}
	res2, err := c2.PurchaseProduct(context.Background(), purchaseReq)
	if err != nil {
		log.Fatalf("Error while calling availablity server %v", err)
	}
	log.Printf("Purchase Status for User %d  Status  %v ", i,res2.Success)

	wg.Done()
}

func randomQuantity() int32 {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 6
	return int32(rand.Intn(max-min+1) + min)
}
