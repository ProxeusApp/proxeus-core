package ethglue

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Dial(rawurl string) (*ethclient.Client, error) {
	return DialContext(context.TODO(), rawurl)
}

func DialContext(ctx context.Context, rawurl string) (c *ethclient.Client, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(20*time.Second))
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
			time.Sleep(time.Duration(shouldBlockOnErrTimeSec*time.Second) - dur)
		}
	}()
	c, err = ethclient.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	h, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	if h == nil || h.Number == nil || h.Number == big.NewInt(0) {
		return nil, errors.New("unreasonable response from eth node")
	}
	return c, nil
}
