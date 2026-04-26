package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// chain creates a chain of errors of the given length, starting with the given error.
	chain := func(length int, init error) error {
		t.Helper()
		e := init
		for i := range length - 1 {
			e = fmt.Errorf("%d: %w", i, e)
		}
		return e
	}

	// unwrapChain unwraps the error chain up to the given depth.
	unwrapChain := func(err error, depth int) error {
		t.Helper()
		for range depth + 1 {
			err = errors.Unwrap(err)
		}
		return err
	}

	defaultChainDepth := 5
	defaultMetadata := []map[string]any{
		{"k1": "v1", "k2": "v2"},
	}

	tests := []struct {
		name             string
		err              error
		public           string
		code             ErrCode
		metadata         []map[string]any
		expectedMetadata map[string]any
		chainDepth       int
	}{
		{
			name:             "fully-populated",
			err:              errors.New("Failed command: 'rm -rf --force /'"),
			public:           "no biggie, try again later",
			code:             ErrUnknown,
			metadata:         defaultMetadata,
			expectedMetadata: defaultMetadata[0],
			chainDepth:       defaultChainDepth,
		},
		{
			name:             "nil-error",
			err:              nil,
			public:           "an ouchie occured :'(",
			code:             ErrUnknown,
			metadata:         defaultMetadata,
			expectedMetadata: defaultMetadata[0],
			chainDepth:       defaultChainDepth,
		},
		{
			name:             "no-metadata",
			err:              errors.New("boom! :: boom2!"),
			public:           "das message 4 da peeps",
			code:             ErrUnknown,
			metadata:         nil,
			expectedMetadata: nil,
			chainDepth:       defaultChainDepth,
		},
		{
			name:             "invalid-error-code",
			err:              errors.New("whoopsies!"),
			public:           "bonjour, welcome to my error",
			code:             99999,
			metadata:         defaultMetadata,
			expectedMetadata: defaultMetadata[0],
			chainDepth:       defaultChainDepth,
		},
		{
			name:   "metadata-override",
			err:    errors.New("who put that there?!"),
			public: "oh well, it's gone now",
			code:   ErrUnknown,
			metadata: []map[string]any{
				defaultMetadata[0],
				{"k3": "v3", "k1": "v100"},
			},
			expectedMetadata: map[string]any{"k1": "v100", "k2": "v2", "k3": "v3"},
			chainDepth:       defaultChainDepth,
		},
		{
			name:             "very-very-deep",
			err:              errors.New("look out down below!"),
			public:           "we're on an adventure, boys",
			code:             ErrUnknown,
			metadata:         nil,
			expectedMetadata: nil,
			chainDepth:       10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := chain(tt.chainDepth, New(tt.err, tt.public, tt.code, tt.metadata...))

			var err *Error
			require.ErrorAs(t, e, &err)
			assert.Equal(t, tt.public, err.Public())
			assert.Equal(t, tt.expectedMetadata, err.Metadata())

			// check that the status code is correct
			status, ok := statusMap[err.Code()]
			require.True(t, ok, "status code not in statusMap: %v", tt.code)
			assert.Equal(t, status, err.HTTPStatus())

			// check that the error code is correct
			if _, ok = statusMap[tt.code]; !ok {
				assert.Equal(t, ErrUnknown, err.Code())
			}

			// check that the error is unwrapped correctly
			assert.Equal(t, tt.err, unwrapChain(e, tt.chainDepth), "not unwrapped correctly")
		})
	}
}

func TestWrap(t *testing.T) {
	t.Parallel()

	// test With builder and overwriting
	err := New(nil, "test", ErrUnexpected).
		With("key1", "value1").
		With("key2", "value2").
		With("key1", "overwritten")

	expected := map[string]any{
		"key1": "overwritten",
		"key2": "value2",
	}

	assert.Equal(t, expected, err.Metadata())
}

func TestShortcuts(t *testing.T) {
	t.Parallel()

	err := BadCredentials(nil, "no, you cannot be Admin...")
	assert.Equal(t, ErrBadCredentials, err.Code())

	err = NotAuthorized(nil, "this is a holy place, you are not invited!")
	assert.Equal(t, ErrNotAuthorized, err.Code())

	err = Unexpected(nil, "hmm, didn't see that coming, did yah? (neither did i)")
	assert.Equal(t, ErrUnexpected, err.Code())

	err = BadInput(nil, "i'm sorry, but your password has to be less than 1GB in size")
	assert.Equal(t, ErrBadInput, err.Code())

	err = Conflict(nil, "MINE! find your own, pesant")
	assert.Equal(t, ErrConflict, err.Code())
}
