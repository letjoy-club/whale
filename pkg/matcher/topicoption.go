package matcher

import (
	"context"
	"fmt"
	"whale/pkg/gqlient/hoopoe"

	"github.com/letjoy-club/mida-tool/midacontext"
)

func getTopicOptions(ctx context.Context) map[string]*hoopoe.TopicOptionConfigFields {
	config, err := hoopoe.GetTopicConfigOptions(ctx, midacontext.GetServices(ctx).Hoopoe, &hoopoe.GraphQLPaginator{Size: 9999})
	if err != nil {
		fmt.Println(err)
	}
	ret := map[string]*hoopoe.TopicOptionConfigFields{}
	for _, config := range config.TopicOptionConfigs {
		ret[config.TopicId] = &config.TopicOptionConfigFields
	}
	return ret
}
