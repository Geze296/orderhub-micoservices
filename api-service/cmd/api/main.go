package main

import (
	"context"
	"fmt"

	"github.com/Geze296/orderhub/api-service/internal/config"
	"github.com/Geze296/orderhub/api-service/internal/db"
)

func main(){
	cfg := config.LoadConfig()

	psgrs, err:= db.NewPostgres(context.Background(),cfg.PostgresURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(psgrs)
}