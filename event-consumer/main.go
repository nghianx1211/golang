package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

var ctx = context.Background()

func main() {
	// Connect Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Run 2 consumers song song
	go consumeTopic("team.activity", rdb)
	go consumeTopic("asset.changes", rdb)

	select {} // block main goroutine
}

func consumeTopic(topic string, rdb *redis.Client) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupID:     "cache-updater-group",
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
	})
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		var event map[string]interface{}
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Println("JSON parse error:", err)
			continue
		}

		log.Printf("[%s] Event: %+v\n", topic, event)
		handleEvent(event, rdb)
	}
}

func handleEvent(event map[string]interface{}, rdb *redis.Client) {
    switch event["eventType"] {
    // --- Team events ---
    case "MEMBER_ADDED":
        key := fmt.Sprintf("team:%s:members", event["teamId"].(string))
        rdb.SAdd(ctx, key, event["targetUserId"].(string))

    case "MEMBER_REMOVED":
        key := fmt.Sprintf("team:%s:members", event["teamId"].(string))
        rdb.SRem(ctx, key, event["targetUserId"].(string))

    // --- Asset metadata ---
    case "FOLDER_CREATED", "FOLDER_UPDATED":
        key := fmt.Sprintf("folder:%s", event["assetId"].(string))
        data, _ := json.Marshal(event)
        rdb.Set(ctx, key, data, 0)

    case "FOLDER_DELETED":
        key := fmt.Sprintf("folder:%s", event["assetId"].(string))
        rdb.Del(ctx, key)

    case "NOTE_CREATED", "NOTE_UPDATED":
        key := fmt.Sprintf("note:%s", event["assetId"].(string))
        data, _ := json.Marshal(event)
        rdb.Set(ctx, key, data, 0)

    case "NOTE_DELETED":
        key := fmt.Sprintf("note:%s", event["assetId"].(string))
        rdb.Del(ctx, key)

    // --- Asset ACL ---
    case "FOLDER_SHARED", "NOTE_SHARED":
        key := fmt.Sprintf("asset:%s:acl", event["assetId"].(string))
        rdb.HSet(ctx, key, event["targetUserId"].(string), event["permission"].(string))

    case "FOLDER_UNSHARED", "NOTE_UNSHARED":
        key := fmt.Sprintf("asset:%s:acl", event["assetId"].(string))
        rdb.HDel(ctx, key, event["targetUserId"].(string))
    }
}