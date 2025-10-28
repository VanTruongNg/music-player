package producer

import (
	"notification/configs"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type ProducerProfile string

const (
	// ProfileBalanced is recommended for most use cases (default)
	// Throughput: ~100k-500k msg/s, Good durability, Fast enough
	ProfileBalanced ProducerProfile = "balanced"

	// ProfileSafe is for critical data that cannot be lost
	// Throughput: ~5k-10k msg/s, Maximum durability, Slower
	ProfileSafe ProducerProfile = "safe"

	// ProfileFast is for high-traffic endpoints with acceptable data loss
	// Throughput: ~500k-1M+ msg/s, Minimum durability, Very fast
	ProfileFast ProducerProfile = "fast"
)

func GetProducerOpts(cfg *configs.KafkaConfig, profile ProducerProfile) []kgo.Opt {
	baseOpts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),
		kgo.ConnIdleTimeout(90 * time.Second),
		kgo.MetadataMaxAge(5 * time.Minute),
		kgo.MetadataMinAge(10 * time.Second),
	}

	if cfg.Debug {
		baseOpts = append(baseOpts, kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, nil)))
	}

	switch profile {
	case ProfileSafe:
		return append(baseOpts, getSafeOpts()...)
	case ProfileFast:
		return append(baseOpts, getFastOpts()...)
	default:
		return append(baseOpts, getBalancedOpts()...)
	}
}

func getBalancedOpts() []kgo.Opt {
	return []kgo.Opt{
		// Wait for leader ACK only (3-5x faster than AllISRAcks)
		kgo.RequiredAcks(kgo.LeaderAck()),

		kgo.DisableIdempotentWrite(),

		// Large batching for throughput
		kgo.ProducerBatchMaxBytes(16_777_216),
		kgo.ProducerLinger(50 * time.Millisecond),
		kgo.MaxBufferedRecords(100_000),

		// Fast compression
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),

		// Retry settings
		kgo.RequestRetries(5),
		kgo.RetryBackoffFn(func(tries int) time.Duration {
			return time.Duration(50*(1<<tries)) * time.Millisecond
		}),

		// Timeouts
		kgo.RequestTimeoutOverhead(20 * time.Second),
		kgo.ProduceRequestTimeout(60 * time.Second),
	}
}

func getSafeOpts() []kgo.Opt {
	return []kgo.Opt{
		// Wait for ALL replicas (maximum durability)
		kgo.RequiredAcks(kgo.AllISRAcks()),

		// Moderate batching (balance speed vs latency)
		kgo.ProducerBatchMaxBytes(4_000_000),
		kgo.ProducerLinger(20 * time.Millisecond),
		kgo.MaxBufferedRecords(50_000),

		// Balanced compression
		kgo.ProducerBatchCompression(kgo.Lz4Compression()),

		// More retries for critical data
		kgo.RequestRetries(10),
		kgo.RetryBackoffFn(func(tries int) time.Duration {
			return time.Duration(100*(1<<tries)) * time.Millisecond
		}),

		// Longer timeouts
		kgo.RequestTimeoutOverhead(30 * time.Second),
		kgo.ProduceRequestTimeout(90 * time.Second),
	}
}

func getFastOpts() []kgo.Opt {
	return []kgo.Opt{
		// Leader ACK only (or NoAck for maximum speed)
		kgo.RequiredAcks(kgo.LeaderAck()),

		kgo.DisableIdempotentWrite(),

		// Maximum batching
		kgo.ProducerBatchMaxBytes(16_777_216),
		kgo.ProducerLinger(50 * time.Millisecond),
		kgo.MaxBufferedRecords(200_000),

		// No compression for maximum CPU efficiency
		// Use Snappy if network is bottleneck
		kgo.ProducerBatchCompression(kgo.NoCompression()),

		// Fast-fail retry
		kgo.RequestRetries(3),
		kgo.RetryBackoffFn(func(tries int) time.Duration {
			return time.Duration(tries*50) * time.Millisecond
		}),

		// Short timeouts
		kgo.RequestTimeoutOverhead(10 * time.Second),
		kgo.ProduceRequestTimeout(30 * time.Second),
	}
}
