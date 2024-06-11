package schemas_test

import (
	// embed used for loading request cases
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed testdata/fail/create.json
var createFailScenarios []byte

//go:embed testdata/success/create.json
var createSuccessScenario []byte

//go:embed create.json
var createRequestSchema []byte

func TestSchemaValidation_Success(t *testing.T) {
	t.Run("Given a valid request", func(t *testing.T) {
		var testcases []testCase
		err := json.Unmarshal(createSuccessScenario, &testcases)
		assert.Nil(t, err)

		loader := gojsonschema.NewBytesLoader(createRequestSchema)
		schema, err := gojsonschema.NewSchema(loader)
		assert.Nil(t, err)
		for _, tc := range testcases {
			t.Run(fmt.Sprintf("Should return valid for scenario: %s", tc.Scenario), func(t *testing.T) {
				requestJSON := gojsonschema.NewBytesLoader(tc.Payload)
				result, err := schema.Validate(requestJSON)
				assert.Nil(t, err)
				assert.True(t, result.Valid())
			})
		}
	})
}

func TestSchemaValidation_Fail(t *testing.T) {
	t.Run("Given an invalid request", func(t *testing.T) {
		var testcases []testCase
		err := json.Unmarshal(createFailScenarios, &testcases)
		assert.Nil(t, err)

		loader := gojsonschema.NewBytesLoader(createRequestSchema)
		schema, err := gojsonschema.NewSchema(loader)
		assert.Nil(t, err)
		for _, tc := range testcases {
			t.Run(fmt.Sprintf("Should return valid for scenario: %s", tc.Scenario), func(t *testing.T) {
				requestJSON := gojsonschema.NewBytesLoader(tc.Payload)
				result, err := schema.Validate(requestJSON)
				assert.Nil(t, err)
				assert.False(t, result.Valid())
			})
		}
	})
}

type testCase struct {
	Scenario string          `json:"scenario"`
	Payload  json.RawMessage `json:"payload"`
}
