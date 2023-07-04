package loader

import (
	"context"
	"whale/pkg/gqlient/hoopoe"

	"github.com/letjoy-club/mida-tool/midacontext"
)

type TopicOptionConfig = hoopoe.GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig
type TopicOptionConfigProperty = hoopoe.TopicOptionConfigFieldsPropertiesTopicOptionProperty
type TopicOptionConfigOptions = hoopoe.TopicOptionConfigFieldsPropertiesTopicOptionPropertyOptionsTopicOption

type TopicOptionConfigLoader struct {
	TopicID string
	m       map[string]*TopicOptionConfig
}

func (l *TopicOptionConfigLoader) LoadAll(ctx context.Context) {
	resp, err := hoopoe.GetTopicConfigOptions(ctx, midacontext.GetServices(ctx).Hoopoe, &hoopoe.GraphQLPaginator{Size: 9999})
	if err != nil {
		return
	}
	m := make(map[string]*TopicOptionConfig)
	for _, v := range resp.TopicOptionConfigs {
		m[v.TopicId] = v
	}
	l.m = m
}

func (l *TopicOptionConfigLoader) Load(ctx context.Context, topicID string) *hoopoe.TopicOptionConfigFields {
	resp, err := hoopoe.GetTopicConfigOption(ctx, midacontext.GetServices(ctx).Hoopoe, topicID)
	if err != nil {
		return nil
	}
	return &resp.TopicOptionConfig.TopicOptionConfigFields
}

func (l *TopicOptionConfigLoader) GetTopicOptionConfig(topicID string) *TopicOptionConfig {
	config, ok := l.m[topicID]
	if !ok {
		return nil
	}
	return config
}

func NewTopicOptionConfigLoader() *TopicOptionConfigLoader {
	return &TopicOptionConfigLoader{}
}
