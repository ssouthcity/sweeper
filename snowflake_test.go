package sweeper

import (
	"testing"

	"github.com/bwmarrin/snowflake"
)

func TestSnowflake(t *testing.T) {
	s := NextSnowflake()

	if _, err := snowflake.ParseString(s.String()); err != nil {
		t.Errorf("expected snowflake to parse correctly, got error %s", err)
	}
}
