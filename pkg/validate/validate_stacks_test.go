package validate

import (
	e "github.com/cloudposse/terraform-provider-utils/internal/exec"
	u "github.com/cloudposse/terraform-provider-utils/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateStacksCommand(t *testing.T) {
	err := e.ExecuteValidateStacks(nil, nil)
	u.PrintError(err)
	assert.NotNil(t, err)
}
