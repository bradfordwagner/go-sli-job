package template

//go:generate mockgen -destination=mocks/pkg/mock_sarama/module.go -package=mock_sarama github.com/Shopify/sarama ConsumerGroup,SyncProducer,ConsumerGroupSession,ConsumerGroupClaim,ClusterAdmin,Client

