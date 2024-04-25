package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/imagekit-developer/imagekit-go"
)

func main() {
	privateKey, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		log.Fatalf("PRIVATE_KEY not found")
	}
	publicKey, ok := os.LookupEnv("PUBLIC_KEY")
	if !ok {
		log.Fatalf("PUBLIC_KEY not found")
	}
	urlEndpoint, ok := os.LookupEnv("URL_ENDPOINT")
	if !ok {
		log.Fatalf("URL_ENDPOINT not found")
	}

	// // Create a new ImageKit object
	imgkit := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		UrlEndpoint: urlEndpoint,
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		url := imgkit.SignToken(imagekit.SignTokenParam{})

		res := make(map[string]interface{})
		res["status"] = "success"
		res["data"] = url

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Printf("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
