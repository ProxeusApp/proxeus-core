package ethglue

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// Connection dialler. To start an Ethereum connection you pass from here.
type Dialler struct {
	connectionTimeout time.Duration
	ethDialler        ETHDiallerIF
}

func NewDefaultDialler() *Dialler {
	return &Dialler{
		connectionTimeout: 10,
		ethDialler:        defaultETHDialler{},
	}
}

func NewCustomDialler(ethDialler ETHDiallerIF, connectionTimeout time.Duration) *Dialler {
	return &Dialler{ethDialler: ethDialler, connectionTimeout: connectionTimeout}
}

func (me Dialler) Dial(rawUrl string) (ethClient ETHClientIF, err error) {
	return me.DialContext(context.Background(), rawUrl)
}

func (me Dialler) DialContext(ctx context.Context, rawUrl string) (ethClient ETHClientIF, err error) {
	ctx, cancel := context.WithTimeout(ctx, me.connectionTimeout*time.Second)
	defer cancel()

	sTime := time.Now()
	// we don't want to return error immediately as we usually retry this func on err
	defer func() {
		if err == nil {
			return
		}

		dur := time.Now().Sub(sTime)
		const shouldBlockOnErrTimeSec = 3
		if dur.Seconds() < shouldBlockOnErrTimeSec {
			time.Sleep(shouldBlockOnErrTimeSec*time.Second - dur)
		}
	}()

	ethClient, err = me.ethDialler.DialContext(ctx, rawUrl)
	if err != nil {
		return nil, err
	}

	h, err := ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	if h == nil || h.Number == nil || h.Number == big.NewInt(0) {
		return nil, errors.New("unreasonable response from eth node")
	}

	return ethClient, nil
}

// Default's ETHDialler. Makes a real connection to Ethereum
type defaultETHDialler struct{}

func (me defaultETHDialler) Dial(rawUrl string) (ethClient ETHClientIF, err error) {
	return me.DialContext(context.Background(), rawUrl)
}

func (me defaultETHDialler) DialContext(ctx context.Context, rawUrl string) (ETHClientIF, error) {
	return ethclient.DialContext(ctx, rawUrl)
}
