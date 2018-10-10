package tgbotbase

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
	return r.SetPropertyForUserInChat(name, user, ChatID(user), value)
}

func (r *RedisPropertyStorage) SetPropertyForChat(name string, chat ChatID, value interface{}) error {
	log.Printf("Setting property '%s' for chat %d with value: %v", name, chat, value)
	return r.SetPropertyForUserInChat(name, 0, chat, value)
}

func (r *RedisPropertyStorage) GetProperty(name string, user UserID, chat ChatID) (string, error) {
	log.Printf("Getting property '%s' for user %d chat %d", name, user, chat)

	// checking specific property value for this user in this chat
	res := r.client.Get(redisPropertyKey(name, user, chat))
	err := res.Err()
	if err != nil {
		if err == redis.Nil {
			log.Printf("No property '%s' for user %d chat %d, checking next", name, user, chat)
		} else {
			return "", err
		}
	} else {
		return res.Val(), nil
	}

	// checking user-defined property (for any chat, set via direct msg)
	res = r.client.Get(redisPropertyKey(name, user, ChatID(user)))
	err = res.Err()
	if err != nil {
		if err == redis.Nil {
			log.Printf("No property '%s' for user %d, checking next", name, user)
		} else {
			return "", err
		}
	} else {
		return res.Val(), nil
	}

	// checking chat-defined property (default property for this chat)
	res = r.client.Get(redisPropertyKey(name, 0, chat))
	err = res.Err()
	if err != nil {
		if err == redis.Nil {
			log.Printf("No property '%s' for chat %d", name, chat)
		} else {
			return "", err
		}
	} else {
		return res.Val(), nil
	}

	log.Printf("No property '%s' for user %d chat %d, returning null", name, user, chat)
	return "", nil
}

var _ PropertyStorage = &RedisPropertyStorage{}
