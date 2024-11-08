package su2

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	maxTokens    = 2
	refillPeriod = time.Second * 5
)

type Bucket struct {
	tokens       int
	maxTokens    int
	lastRefill   time.Time
	refillPeriod time.Duration
	mu           sync.Mutex
}

func (bucket *Bucket) InitBucket() {
	bucket.tokens = 0
	bucket.maxTokens = maxTokens
	bucket.lastRefill = time.Now()
	bucket.refillPeriod = refillPeriod
	bucket.mu = sync.Mutex{}
}

func (bucket *Bucket) Refill() {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	if time.Since(bucket.lastRefill) < bucket.refillPeriod {
		return
	}

	bucket.tokens = bucket.maxTokens
	bucket.lastRefill = time.Now()
}

func (bucket *Bucket) Allow() bool {
	bucket.Refill()

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

func TokenBucketRLMiddleware(bucket *Bucket) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !bucket.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ChiRateLimiter() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	var rlBucket Bucket
	rlBucket.InitBucket()

	r.Use(TokenBucketRLMiddleware(&rlBucket))

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("tokens left: " + strconv.Itoa(rlBucket.tokens)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	r.Get("/", handlerFunc)

	http.ListenAndServe(":8081", r)
}
