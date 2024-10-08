package main

import (
	"context"
	"fmt"
	"os"
	usecase "test/lambda/handler"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambdaV2 // TODO documentar sobre a integração V2 do API Gateway

func init() {
	r := gin.Default()
	r.POST("/dev/test", usecase.HandlePostRequest)
	ginLambda = ginadapter.NewV2(r)
}

func HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {

	if os.Getenv("STAGE") == "dev" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Erro ao carregar variáveis de ambiente", err)
			return
		}
		if os.Getenv("ENVIRONMENT") == "dev" {
			r := gin.Default()
			r.POST("/test", usecase.HandlePostRequest)
			address := fmt.Sprintf(":%s", os.Getenv("PORT"))
			err := r.Run(address)
			if err != nil {
				return
			}
		}
	}

	lambda.Start(HandleRequest)

}
