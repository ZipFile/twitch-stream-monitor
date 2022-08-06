package helix

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

func errorFromResponse(response *helix.ResponseCommon) error {
	if response == nil {
		return nil
	}

	return fmt.Errorf(
		"%s (%d): %s",
		response.Error,
		response.ErrorStatus,
		response.ErrorMessage,
	)
}
