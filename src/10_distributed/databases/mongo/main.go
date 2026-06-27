// MongoDB Advanced Patterns: Aggregation pipelines, Change Streams, Indexes.
//
// LEARNING NOTES:
// - Aggregation pipelines: multi-stage data transformations (like Unix pipes)
// - Change Streams: real-time notifications when documents change (requires replica set)
// - Indexes: B-tree indexes, compound indexes, TTL indexes for auto-expiry
// - vs Postgres: schema-flexible, great for hierarchical data, horizontal scaling via sharding
//
// For C++ devs: Think of aggregation as std::transform | std::accumulate chained together.
// Change streams are like inotify for your database.
//
// Run: go run main.go
// Requires: MongoDB on localhost:27017
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Category  string             `bson:"category" json:"category"`
	Price     float64            `bson:"price" json:"price"`
	Stock     int                `bson:"stock" json:"stock"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type CategoryStats struct {
	Category     string  `bson:"_id" json:"category"`
	TotalProducts int    `bson:"total_products" json:"total_products"`
	AvgPrice     float64 `bson:"avg_price" json:"avg_price"`
	TotalStock   int     `bson:"total_stock" json:"total_stock"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("learning_db")
	products := db.Collection("products")

	// --- Create indexes ---
	indexModels := []mongo.IndexModel{
		{Keys: bson.D{{Key: "category", Value: 1}}},                    // Single field index
		{Keys: bson.D{{Key: "category", Value: 1}, {Key: "price", Value: -1}}}, // Compound index
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(86400 * 30), // TTL: auto-delete after 30 days
		},
	}
	_, err = products.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Printf("Index creation: %v", err)
	}
	fmt.Println("Indexes created")

	// --- Seed data ---
	seedProducts := []interface{}{
		Product{Name: "Go in Action", Category: "books", Price: 39.99, Stock: 100, CreatedAt: time.Now()},
		Product{Name: "Kubernetes Up", Category: "books", Price: 49.99, Stock: 50, CreatedAt: time.Now()},
		Product{Name: "Mechanical KB", Category: "electronics", Price: 149.99, Stock: 25, CreatedAt: time.Now()},
		Product{Name: "USB-C Hub", Category: "electronics", Price: 29.99, Stock: 200, CreatedAt: time.Now()},
		Product{Name: "Standing Desk", Category: "furniture", Price: 599.99, Stock: 10, CreatedAt: time.Now()},
	}
	products.InsertMany(ctx, seedProducts)

	// --- Aggregation Pipeline ---
	// Group by category, compute stats
	pipeline := mongo.Pipeline{
		// Stage 1: Group by category
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category"},
			{Key: "total_products", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "avg_price", Value: bson.D{{Key: "$avg", Value: "$price"}}},
			{Key: "total_stock", Value: bson.D{{Key: "$sum", Value: "$stock"}}},
		}}},
		// Stage 2: Sort by total products descending
		{{Key: "$sort", Value: bson.D{{Key: "total_products", Value: -1}}}},
	}

	cursor, err := products.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	fmt.Println("\n--- Category Stats (Aggregation Pipeline) ---")
	for cursor.Next(ctx) {
		var stat CategoryStats
		cursor.Decode(&stat)
		fmt.Printf("  %s: %d products, avg $%.2f, %d total stock\n",
			stat.Category, stat.TotalProducts, stat.AvgPrice, stat.TotalStock)
	}

	// --- Find with filter + projection ---
	fmt.Println("\n--- Electronics under $100 ---")
	filter := bson.M{"category": "electronics", "price": bson.M{"$lt": 100}}
	opts := options.Find().SetProjection(bson.M{"name": 1, "price": 1, "_id": 0})

	cur, _ := products.Find(ctx, filter, opts)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		cur.Decode(&result)
		fmt.Printf("  %s: $%.2f\n", result["name"], result["price"])
	}

	// --- Cleanup ---
	products.Drop(ctx)
	fmt.Println("\nDemo complete! Collection dropped.")
}
