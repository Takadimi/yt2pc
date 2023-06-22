package youtube_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Takadimi/yt2pc/core/youtube"
)

func TestGetVideoData(t *testing.T) {
	ctx := context.Background()
	vd, err := youtube.GetVideoData(ctx, "BaW_jenozKc")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", vd)
}
