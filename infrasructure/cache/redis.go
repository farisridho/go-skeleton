package redis

import (
	"context"
	"fmt"
	"time"

	v8 "github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go/log"
	"github.com/segmentio/encoding/json"
)

type redis struct {
	debug  bool
	mode   Mode
	prefix string
	client client
}

func New(config Config, opts ...Option) (cache.Cache, error) {
	rdb := &redis{
		debug:  config.Debug,
		prefix: config.Prefix,
		mode:   config.Mode,
	}

	for i := range opts {
		if err := opts[i](rdb); err != nil {
			return nil, err
		}
	}

	if rdb.client == nil {
		if config.Mode == Single {
			if len(config.Addresses) == 0 {
				return nil, ErrorAddressesCannotBeEmpty
			}

			rdb.client = v8.NewClient(&v8.Options{
				Network:            config.Network,
				Addr:               config.Addresses[0],
				Username:           config.Username,
				Password:           config.Password,
				DB:                 config.DB,
				MaxRetries:         config.MaxRetries,
				MinRetryBackoff:    config.MinRetryBackoff,
				MaxRetryBackoff:    config.MaxRetryBackoff,
				DialTimeout:        config.DialTimeout,
				ReadTimeout:        config.ReadTimeout,
				WriteTimeout:       config.WriteTimeout,
				PoolSize:           config.PoolSize,
				MinIdleConns:       config.MinIdleConnections,
				MaxConnAge:         config.MaxConnectionAge,
				PoolTimeout:        config.PoolTimeout,
				IdleTimeout:        config.IdleTimeout,
				IdleCheckFrequency: config.IdleCheckFrequency,
			})
		}

		if config.Mode == Sentinel {
			if len(config.Addresses) == 0 {
				return nil, ErrorAddressesCannotBeEmpty
			}

			rdb.client = v8.NewFailoverClient(&v8.FailoverOptions{
				MasterName:         config.MasterName,
				SentinelAddrs:      config.Addresses,
				SentinelPassword:   config.SentinelPassword,
				RouteByLatency:     config.RouteByLatency,
				RouteRandomly:      config.RouteRandomly,
				SlaveOnly:          config.SlaveOnly,
				Username:           config.Username,
				Password:           config.Password,
				DB:                 config.DB,
				MaxRetries:         config.MaxRetries,
				MinRetryBackoff:    config.MinRetryBackoff,
				MaxRetryBackoff:    config.MaxRetryBackoff,
				DialTimeout:        config.DialTimeout,
				ReadTimeout:        config.ReadTimeout,
				WriteTimeout:       config.WriteTimeout,
				PoolSize:           config.PoolSize,
				MinIdleConns:       config.MinIdleConnections,
				MaxConnAge:         config.MaxConnectionAge,
				PoolTimeout:        config.PoolTimeout,
				IdleTimeout:        config.IdleTimeout,
				IdleCheckFrequency: config.IdleCheckFrequency,
			})
		}

		if config.Mode == Cluster {
			if len(config.Addresses) == 0 {
				return nil, ErrorAddressesCannotBeEmpty
			}

			rdb.client = v8.NewClusterClient(&v8.ClusterOptions{
				Addrs:              config.Addresses,
				MaxRedirects:       config.MaxRedirects,
				ReadOnly:           config.ReadOnly,
				RouteByLatency:     config.RouteByLatency,
				RouteRandomly:      config.RouteRandomly,
				Username:           config.Username,
				Password:           config.Password,
				MaxRetries:         config.MaxRetries,
				MinRetryBackoff:    config.MinRetryBackoff,
				MaxRetryBackoff:    config.MaxRetryBackoff,
				DialTimeout:        config.DialTimeout,
				ReadTimeout:        config.ReadTimeout,
				WriteTimeout:       config.WriteTimeout,
				PoolSize:           config.PoolSize,
				MinIdleConns:       config.MinIdleConnections,
				MaxConnAge:         config.MaxConnectionAge,
				PoolTimeout:        config.PoolTimeout,
				IdleTimeout:        config.IdleTimeout,
				IdleCheckFrequency: config.IdleCheckFrequency,
			})
		}
	}

	if rdb.client == nil {
		return nil, ErrorNoClient
	}

	if rdb.client != nil {
		if _, err := rdb.client.Ping(context.Background()).Result(); err != nil {
			return nil, ErrorFailedPing
		}
	}

	return rdb, nil
}

func (r redis) Get(ctx context.Context, key string, value interface{}) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache.Redis.Get")
	defer func() {
		span.SetTag("error", err != nil)

		if err != nil {
			span.LogFields(log.String("error", err.Error()))
		}

		span.LogFields(log.Object("value", value))
		span.Finish()
	}()

	span.LogFields(log.String("key", r.key(key)))
	span.LogFields(log.String("mode", r.mode.String()))

	get := r.client.Get(ctx, r.key(key))

	if r.debug {
		fmt.Println(get.Args())
	}

	if get.Err() != nil {
		err = get.Err()
		return
	}

	data, _ := get.Bytes()

	err = json.Unmarshal(data, &value)

	return
}

func (r redis) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache.Redis.Set")
	defer func() {
		span.SetTag("error", err != nil)

		if err != nil {
			span.LogFields(log.String("error", err.Error()))
		}

		span.Finish()
	}()

	span.LogFields(log.Object("value", value))
	span.LogFields(log.String("key", r.key(key)))
	span.LogFields(log.String("mode", r.mode.String()))

	var valueBytes []byte
	if valueBytes, err = json.Marshal(value); err != nil {
		return
	}

	set := r.client.Set(ctx, r.key(key), valueBytes, expiration)

	if r.debug {
		fmt.Println(set.Args())
	}

	if set.Err() != nil {
		err = set.Err()
		return
	}

	return
}

func (r redis) SetIfAbsent(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache.Redis.SetIfAbsent")
	defer func() {
		span.SetTag("error", err != nil)

		if err != nil {
			span.LogFields(log.String("error", err.Error()))
		}

		span.Finish()
	}()

	span.LogFields(log.Object("value", value))
	span.LogFields(log.String("key", r.key(key)))
	span.LogFields(log.String("mode", r.mode.String()))

	var valueBytes []byte
	if valueBytes, err = json.Marshal(value); err != nil {
		return
	}

	add := r.client.SetNX(ctx, r.key(key), valueBytes, expiration)

	if r.debug {
		fmt.Println(add.Args())
	}

	if add.Err() != nil {
		err = add.Err()
		return
	}

	if add.Val() {
		return
	}

	return ErrorKeyAlreadySet
}

func (r redis) Delete(ctx context.Context, key string) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache.Redis.Delete")
	defer func() {
		span.SetTag("error", err != nil)

		if err != nil {
			span.LogFields(log.String("error", err.Error()))
		}

		span.Finish()
	}()

	span.LogFields(log.String("key", r.key(key)))
	span.LogFields(log.String("mode", r.mode.String()))

	del := r.client.Del(ctx, r.key(key))

	if r.debug {
		fmt.Println(del.Args())
	}

	if del.Err() != nil {
		err = del.Err()
		return
	}

	if del.Val() == 0 {
		err = ErrorDelete
		return
	}

	return
}

func (r redis) FlushDB(ctx context.Context) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache.Redis.FlushDB")
	defer func() {
		span.SetTag("error", err != nil)

		if err != nil {
			span.LogFields(log.String("error", err.Error()))
		}

		span.Finish()
	}()

	span.LogFields(log.String("mode", r.mode.String()))

	flushDB := r.client.FlushDB(ctx)

	if r.debug {
		fmt.Println(flushDB.Args())
	}

	if flushDB.Err() != nil {
		err = flushDB.Err()
		return
	}

	return
}

func (r redis) Close(ctx context.Context) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache.Redis.Close")
	defer func() {
		span.SetTag("error", err != nil)
		span.Finish()
	}()

	return r.client.Close()
}

func (r redis) key(key string) string {
	return r.prefix + key
}
