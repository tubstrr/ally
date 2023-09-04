package ally_redis

import (
	"context"
	"crypto/tls"

	"github.com/redis/go-redis/v9"
	"github.com/tubstrr/ally/environment"
	ally_strings "github.com/tubstrr/ally/utilities/strings"
)

func MakeClient() *redis.Client {
	host := environment.Get_environment_variable("ALLY_REDIS_HOST", "redis")
	port := environment.Get_environment_variable("ALLY_REDIS_PORT", "6379")
	db_name := environment.Get_environment_variable("ALLY_REDIS_DB_NAME", "")
	user := environment.Get_environment_variable("ALLY_REDIS_USERNAME", "default")
	password := environment.Get_environment_variable("ALLY_REDIS_PASSWORD", "")
	tls_enabled := environment.Get_environment_variable("ALLY_REDIS_TLS", "false")

	options := &redis.Options{
		Addr:	  host + ":" + port,
		Password: password, // no password set
		DB:		  0,  // use default DB
	}
	if (db_name != "") {
		options.DB = ally_strings.StringToNumber(db_name)
	}
	if (user != "") {
		options.Username = user
	}
	if (tls_enabled == "true") {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		}
	}

	client := redis.NewClient(options)

	return client
}

func SetKey(key string, value string) {
	ctx := context.Background()
	
	client := MakeClient()
	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func GetKey(key string) (string, error)  {
	ctx := context.Background()
	
	client := MakeClient()

	val, err := client.Get(ctx, key).Result()
	return val, err
}

func GetAllKeys() []string {
	ctx := context.Background()
	
	client := MakeClient()
	val, err := client.Keys(ctx, "*").Result()
	if err != nil {
		panic(err)
	}
	return val
}

func GetKeysByPattern(pattern string) []string {
	ctx := context.Background()
	
	client := MakeClient()
	val, err := client.Keys(ctx, pattern).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func DeleteKey(key string) {
	ctx := context.Background()
	
	client := MakeClient()
	err := client.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}

func DeleteAllKeys() {
	ctx := context.Background()
	
	client := MakeClient()
	err := client.FlushAll(ctx).Err()
	if err != nil {
		panic(err)
	}
}

func DeleteKeysByPattern(pattern string) {
	ctx := context.Background()
	
	client := MakeClient()
	val, err := client.Keys(ctx, pattern).Result()
	if err != nil {
		panic(err)
	}
	for _, key := range val {
		err := client.Del(ctx, key).Err()
		if err != nil {
			panic(err)
		}
	}
}
