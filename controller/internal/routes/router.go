package routes

import (
	"context"
	"controller-service-grpc/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/pablogolobaro/pdfcomposer"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
)

func Start(port string) {
	r := chi.NewRouter()
	httpClient := &http.Client{Timeout: time.Second * 5}
	cwt, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(cwt, "localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	pdfComposeClient := pdfcomposer.NewPdfComposeClient(conn)
	h := &handlers.Handler{HttpClient: httpClient, PdfComposeCLient: pdfComposeClient}

	r.Get("/", h.Web)
	r.Post("/", h.Send)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
