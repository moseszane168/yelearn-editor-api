package redisClient

import (
	"crypto/tls"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

// RedisKeyPrefix defines prefix for iam analytics key.
const (
	RedisKeyPrefix      = "crf-mold-"
	defaultRedisAddress = "127.0.0.1:6379"
)

var redisClusterSingleton redis.UniversalClient

type RedisOptions struct {
	Host                  string
	Port                  int
	Addrs                 []string
	Username              string
	Password              string
	Database              int
	MaxIdle               int
	MaxActive             int
	Timeout               int
	UseSSL                bool
	SSLInsecureSkipVerify bool
}

// RedisClusterStorageManager is a storage manager that uses the redis database.
type RedisClusterStorageManager struct {
	db        redis.UniversalClient
	KeyPrefix string
	HashKeys  bool
	Config    RedisOptions
}

// NewRedisClusterPool returns a redis cluster client.
func NewRedisClusterPool(forceReconnect bool, config RedisOptions) redis.UniversalClient {
	if !forceReconnect {
		if redisClusterSingleton != nil {
			logrus.Debug("Redis pool already INITIALIZED")

			return redisClusterSingleton
		}
	} else {
		if redisClusterSingleton != nil {
			redisClusterSingleton.Close()
		}
	}

	logrus.Debug("Creating new Redis connection pool")

	maxActive := 500
	if config.MaxActive > 0 {
		maxActive = config.MaxActive
	}

	timeout := 5 * time.Second

	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	var tlsConfig *tls.Config
	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}

	var client redis.UniversalClient
	opts := &RedisOpts{
		Addrs:        getRedisAddrs(config),
		DB:           config.Database,
		Password:     config.Password,
		PoolSize:     maxActive,
		IdleTimeout:  240 * time.Second,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		DialTimeout:  timeout,
		TLSConfig:    tlsConfig,
	}

	logrus.Info("--> [REDIS] Creating single-node client")
	client = redis.NewClient(opts.simple())

	redisClusterSingleton = client

	return client
}

func getRedisAddrs(config RedisOptions) (addrs []string) {
	if len(config.Addrs) != 0 {
		addrs = config.Addrs
	}

	if len(addrs) == 0 && config.Port != 0 {
		addr := config.Host + ":" + strconv.Itoa(config.Port)
		addrs = append(addrs, addr)
	}

	return addrs
}

type RedisOpts redis.UniversalOptions

func (o *RedisOpts) simple() *redis.Options {
	addr := defaultRedisAddress
	if len(o.Addrs) > 0 {
		addr = o.Addrs[0]
	}

	return &redis.Options{
		Addr:      addr,
		OnConnect: o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

// Init initialize the redis cluster storage manager.
func (r *RedisClusterStorageManager) Init(config interface{}) error {
	r.Config = RedisOptions{}
	err := mapstructure.Decode(config, &r.Config)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	r.KeyPrefix = RedisKeyPrefix

	return nil
}

// Connect will establish a connection to the r.db.
func (r *RedisClusterStorageManager) Connect() bool {
	if r.db == nil {
		logrus.Debug("Connecting to redis cluster")
		r.db = NewRedisClusterPool(false, r.Config)

		return true
	}

	logrus.Debug("Storage Engine already initialized...")

	// Reset it just in case
	r.db = redisClusterSingleton

	return true
}

func (r *RedisClusterStorageManager) hashKey(in string) string {
	return in
}

func (r *RedisClusterStorageManager) fixKey(keyName string) string {
	setKeyName := r.KeyPrefix + r.hashKey(keyName)

	logrus.Debugf("Input key was: %s", setKeyName)

	return setKeyName
}

// GetAndDeleteSet get and delete key from redis.
func (r *RedisClusterStorageManager) GetAndDeleteSet(keyName string) []interface{} {
	logrus.Debugf("Getting raw key set: %s", keyName)

	if r.db == nil {
		logrus.Warn("Connection dropped, connecting..")
		r.Connect()

		return r.GetAndDeleteSet(keyName)
	}

	logrus.Debugf("keyName is: %s", keyName)

	fixedKey := r.fixKey(keyName)

	logrus.Debugf("Fixed keyname is: %s", fixedKey)

	var lrange *redis.StringSliceCmd
	_, err := r.db.TxPipelined(func(pipe redis.Pipeliner) error {
		lrange = pipe.LRange(fixedKey, 0, -1)
		pipe.Del(fixedKey)

		return nil
	})
	if err != nil {
		logrus.Errorf("Multi command failed: %s", err)
		r.Connect()
	}

	vals := lrange.Val()

	result := make([]interface{}, len(vals))
	for i, v := range vals {
		result[i] = v
	}

	logrus.Debugf("Unpacked vals: %d", len(result))

	return result
}

// SetKey will create (or update) a key value in the store.
func (r *RedisClusterStorageManager) SetKey(keyName, session string, timeout int64) error {

	r.ensureConnection()
	err := r.db.Set(r.fixKey(keyName), session, 0).Err()
	if timeout > 0 {
		if expErr := r.SetExp(keyName, timeout); expErr != nil {
			return expErr
		}
	}
	if err != nil {
		logrus.Errorf("Error trying to set value: %s", err.Error())

		return errors.Wrap(err, "failed to set key")
	}

	return nil
}

// GetKey will create (or update) a key value in the store.
func (r *RedisClusterStorageManager) GetKey(keyName string) (string, error) {

	r.ensureConnection()
	val, err := r.db.Get(r.fixKey(keyName)).Result()
	if err == redis.Nil {
		fmt.Printf("key %s does not exist", keyName)
		return "", err
	} else if err != nil {
		logrus.Errorf("Error trying to set value: %s", err.Error())
		return "", errors.Wrap(err, "failed to set key")
	} else {
		return val, nil
	}
}

// SetExp is used to set the expiry of a key.
func (r *RedisClusterStorageManager) SetExp(keyName string, timeout int64) error {
	err := r.db.Expire(r.fixKey(keyName), time.Duration(timeout)*time.Second).Err()
	if err != nil {
		logrus.Errorf("Could not EXPIRE key: %s", err.Error())
	}

	return errors.Wrap(err, "failed to set expire time for key")
}

func (r *RedisClusterStorageManager) ensureConnection() {
	if r.db != nil {
		// already connected
		return
	}
	logrus.Info("Connection dropped, reconnecting...")
	for {
		r.Connect()
		if r.db != nil {
			// reconnection worked
			fmt.Println("连接正常....")
			return
		}
		fmt.Println("重新连接。。。。。。")
		logrus.Info("Reconnecting again...")
	}
}

var RedisClient *RedisClusterStorageManager

func InitClient() {
	RedisClient = &RedisClusterStorageManager{}
	cfg := RedisOptions{
		Host:      viper.GetString("redis.host"),
		Port:      viper.GetInt("redis.port"),
		Password:  viper.GetString("redis.password"),
		Database:  viper.GetInt("redis.db"),
		MaxIdle:   viper.GetInt("redis.optimisation-max-idle"),
		MaxActive: viper.GetInt("redis.optimisation-max-active"),
	}

	if err := RedisClient.Init(cfg); err != nil {
		panic(err)
	}
}
