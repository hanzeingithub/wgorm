package log

import (
	"context"
	"testing"
)

func TestInfoOf(t *testing.T) {
	logger := WgormLogger{}
	ctx := context.Background()
	logger.InfoOf(ctx, "%s hello word %d", "fuck", 10086)
	logger.WarnOf(ctx, "%s hello word %d", "fuck", 10086)
	logger.ErrorOf(ctx, "%s hello word %d", "fuck", 10086)
}
