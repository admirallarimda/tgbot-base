package botbase

import "log"
import "fmt"
import "github.com/go-redis/redis"

type RedisPropertyStorage struct {
	client *redis.Client
}

func NewRedisPropertyStorage(pool RedisPool) *RedisPropertyStorage {
	r := &RedisPropertyStorage{client: pool.GetConnByName("property")}
	return r
}

func redisPropertyKey(name string, user UserID, chat ChatID) string {
	return fmt.Sprintf("tg:property:%s:%d:%d", name, user, chat)
}

func (r *RedisPropertyStorage) SetPropertyForUserInChat(name string, user UserID, chat ChatID, value interface{}) error {
	log.Printf("Setting property '%s' for user %d chat %d with value: %v", name, user, chat, value)
	key := redisPropertyKey(name, user, chat)
	return r.client.Set(key, value, 0).Err()
}

func (r *RedisPropertyStorage) SetPropertyForUser(name string, user UserID, value interface{}) error {
	log.Printf("Setting property '%s' for user %d with value: %v", name, user, value)
	return r.SetPropertyForUserInChat(name, user, 0, value)
}

func (r *RedisPropertyStorage) SetPropertyForChat(name string, chat ChatID, value interface{}) error {
	log.Printf("Setting property '%s' for chat %d with value: %v", name, chat, value)
	return r.SetPropertyForUserInChat(name, 0, chat, value)
}

func (r *RedisPropertyStorage) GetProperty(name string, user UserID, chat ChatID) (string, error) {
	log.Printf("Getting property '%s' for user %d chat %d", name, user, chat)
	return "", nil
}

var _ PropertyStorage = &RedisPropertyStorage{}
