package common

import (
	"fmt"

	"github.com/kpli0rn/jdwpgo/protocol/basetypes"
)

// ThreadID represents a threadID
type ThreadID basetypes.JWDPObjectID

func (t *ThreadID) String() string {
	return fmt.Sprintf("ThreadID: %s", ((*basetypes.JWDPObjectID)(t)).String())
}

// ThreadGroupID represents a thread Group ID
type ThreadGroupID basetypes.JWDPObjectID

func (t *ThreadGroupID) String() string {
	return fmt.Sprintf("ThreadGroupID: %s", ((*basetypes.JWDPObjectID)(t)).String())
}
