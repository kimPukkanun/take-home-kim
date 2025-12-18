package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	baleriondocs "baleriontakehome/docs"
	"baleriontakehome/internal/httpapi"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}

	srv := httpapi.NewServer()

	// swagger UI
	srv.Mux().Handle(
		"/docs/",
		httpSwagger.Handler(
			httpSwagger.URL("/docs/doc.json"),
		),
	)
	srv.Mux().HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		http.Redirect(w, r, "/docs/index.html", http.StatusFound)
	})

	fmt.Println("Thai Baht Text API")
	fmt.Println("- API:     http://localhost" + addr + "/v1/baht-text")
	fmt.Println("- Swagger: http://localhost" + addr + "/docs")

	// set swagger Host
	host := addr
	if strings.HasPrefix(addr, ":") {
		host = "localhost" + addr
	}
	baleriondocs.SwaggerInfo.Host = host

	log.Fatal(http.ListenAndServe(addr, srv.Handler()))
}
