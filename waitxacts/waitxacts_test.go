package waitxacts

import (
	"context"
	"github.com/lesovsky/noisia"
	"github.com/lesovsky/noisia/db"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfig_validate(t *testing.T) {
	testcases := []struct {
		valid  bool
		config Config
	}{
		{valid: true, config: Config{Jobs: 2, WaitXactsLocktimeMin: 5, WaitXactsLocktimeMax: 10}},
		{valid: false, config: Config{Jobs: 1}},
		{valid: false, config: Config{Jobs: 2, WaitXactsLocktimeMin: 5, WaitXactsLocktimeMax: 4}},
		{valid: false, config: Config{Jobs: 2, WaitXactsLocktimeMin: 0, WaitXactsLocktimeMax: 0}},
	}

	for _, tc := range testcases {
		if tc.valid {
			assert.NoError(t, tc.config.validate())
		} else {
			assert.Error(t, tc.config.validate())
		}
	}
}

func TestWorkload_Run(t *testing.T) {
	config := Config{
		PostgresConninfo:     db.TestConninfo,
		Jobs:                 2,
		WaitXactsLocktimeMin: 1,
		WaitXactsLocktimeMax: 2,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	w, err := NewWorkload(config)
	assert.NoError(t, err)
	err = w.Run(ctx)
	assert.Nil(t, err)

	assert.NoError(t, noisia.Cleanup(context.Background(), config.PostgresConninfo))
}

func Test_selectRandomTable(t *testing.T) {
	testcases := []struct {
		tables []string
		want   int
	}{
		{tables: []string{"test.test1", "test.test2", "test.test3"}, want: 10},
		{tables: []string{}, want: 0},
	}

	for _, tc := range testcases {
		assert.Equal(t, tc.want, len(selectRandomTable(tc.tables)))
	}
}
