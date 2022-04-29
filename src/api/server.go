package api

import (
	"fmt"
	"log"
	"os"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/controllers"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

var (
	MLatencyMs = stats.Float64("latency", "The latency in milliseconds", "ms")
)
var (
	KeyMethod, _ = tag.NewKey("method")
)

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	// seed.Load(server.DB)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)

	// view1 := &view.View{
	// 	Name:        "dist",
	// 	Measure:     MLatencyMs,
	// 	Description: "The dist of the latencies",
	// 	TagKeys:     []tag.Key{KeyMethod},
	// 	Aggregation: view.Distribution(0, 10, 100, 1000, 10000, 100000),
	// }

	// view2 := &view.View{
	// 	Name:        "last",
	// 	Measure:     MLatencyMs,
	// 	Description: "The last of the latencies",
	// 	TagKeys:     []tag.Key{KeyMethod},
	// 	Aggregation: view.LastValue(),
	// }

	// if err := view.Register(view1, view2); err != nil {
	// 	log.Fatalf("Failed to register the views: %v", err)
	// }

	// pe, err := prometheus.NewExporter(prometheus.Options{
	// 	Namespace: "distlast",
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to create the Prometheus stats exporter: %v", err)
	// }

	// go func() {
	// 	mux := http.NewServeMux()
	// 	mux.Handle("/metrics", pe)
	// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
	// }()

	// rand.Seed(time.Now().UnixNano())
	// ctx := context.Background()

	// for {
	// 	n := rand.Intn(100)
	// 	log.Printf("[loop] n=%d\n", n)
	// 	stats.Record(ctx, MLatencyMs.M(float64(time.Duration(n))))
	// 	time.Sleep(1 * time.Second)
	// }
}
