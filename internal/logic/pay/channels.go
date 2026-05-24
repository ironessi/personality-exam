package pay

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type Channels struct{}

func NewChannels() *Channels {
	return &Channels{}
}

type ChannelInfo struct {
	Channel     string  `json:"channel"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Enabled     bool    `json:"enabled"`
	MinAmount   float64 `json:"minAmount"`
	MaxAmount   float64 `json:"maxAmount"`
	Description string  `json:"description"`
}

type GetChannelsOutput struct {
	Channels []*ChannelInfo `json:"channels"`
}

func (c *Channels) GetChannels(ctx context.Context) (*GetChannelsOutput, error) {
	channels := []*ChannelInfo{
		{
			Channel:     "alipay",
			Name:        "支付宝",
			Icon:        g.Cfg("payment").MustGet(ctx, "channels.alipay.icon").String(),
			Enabled:     g.Cfg("payment").MustGet(ctx, "channels.alipay.enabled", true).Bool(),
			MinAmount:   0.01,
			MaxAmount:   50000.00,
			Description: "推荐使用支付宝",
		},
		{
			Channel:     "wechat",
			Name:        "微信支付",
			Icon:        g.Cfg("payment").MustGet(ctx, "channels.wechat.icon").String(),
			Enabled:     g.Cfg("payment").MustGet(ctx, "channels.wechat.enabled", true).Bool(),
			MinAmount:   0.01,
			MaxAmount:   50000.00,
			Description: "微信扫码支付",
		},
	}

	return &GetChannelsOutput{
		Channels: channels,
	}, nil
}