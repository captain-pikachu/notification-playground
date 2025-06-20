package main

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"os"
	"sync"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"gonum.org/v1/gonum/stat/distuv"
)

type quoteGen struct {
	lock       sync.RWMutex
	normalDist distuv.Normal
	nextQuote  float64
}

func newQuoteGen() *quoteGen {
	return &quoteGen{
		nextQuote: 100,
		normalDist: distuv.Normal{
			Mu:    0,
			Sigma: 0.01,
			Src:   rand.NewPCG(0, 0),
		},
	}
}

func (qg *quoteGen) getQuote() float64 {
	qg.lock.Lock()
	defer qg.lock.Unlock()
	res := qg.nextQuote
	qg.nextQuote *= 1 + qg.normalDist.Rand()
	return res
}

func (qg *quoteGen) setNextQuote(nextQuote float64) {
	qg.lock.Lock()
	defer qg.lock.Unlock()
	qg.nextQuote = nextQuote
}

type quotesRes struct {
	Quote float64 `json:"quote"`
}

type setNextQuoteReq struct {
	NextQuote float64 `json:"nextQuote"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Timestamp().Logger().Level(zerolog.DebugLevel)

	quoteGen := newQuoteGen()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/quotes", func(w http.ResponseWriter, r *http.Request) {
		res := quotesRes{Quote: quoteGen.getQuote()}
		resJson, err := json.Marshal(res)
		if err != nil {
			logger.Error().Err(err).Send()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resJson)
	})

	mux.HandleFunc("POST /v1/setNextQuote", func(w http.ResponseWriter, r *http.Request) {
		var setNextQuoteReq setNextQuoteReq
		if err := json.NewDecoder(r.Body).Decode(&setNextQuoteReq); err != nil {
			logger.Error().Err(err).Send()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		quoteGen.setNextQuote(setNextQuoteReq.NextQuote)
	})

	handler := cors.New(cors.Options{AllowedHeaders: []string{"*"}}).Handler(mux)

	if err := http.ListenAndServe("0.0.0.0:3001", handler); err != nil {
		logger.Fatal().Err(err).Send()
	}
}
