package youtube_test

import (
	"context"
	"testing"

	"github.com/Takadimi/yt2pc/core/youtube"
)

func TestGetVideoDataVideoFound(t *testing.T) {
	ctx := context.Background()
	_, err := youtube.GetVideoData(ctx, "BaW_jenozKc")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetVideoDataNotFound(t *testing.T) {
	ctx := context.Background()
	_, err := youtube.GetVideoData(ctx, "garbage")
	if err != youtube.ErrVideoNotAvailable {
		t.Fatal("expected ErrVideoNotAvailable, but error is nil")
	}
}
