package sweeper

import "github.com/bwmarrin/snowflake"

var node, _ = snowflake.NewNode(0)

type Snowflake string

func (s Snowflake) String() string {
	return string(s)
}

func NextSnowflake() Snowflake {
	return Snowflake(node.Generate().String())
}
