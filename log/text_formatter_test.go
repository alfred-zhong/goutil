package log

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_TextFormatter(t *testing.T) {
	tf := &TextFormatter{
		MaxMessageSize: 10,
	}

	entry := logrus.NewEntry(logrus.New())
	entry.Message = "foobar_foobar"
	bb, err := tf.Format(entry)
	require.NoError(t, err)
	t.Log(string(bb))
}

func Benchmark_TextFormatter(b *testing.B) {
	tf := &TextFormatter{
		MaxMessageSize: 10,
	}

	entry := logrus.NewEntry(logrus.New())
	entry.Message = strings.Repeat("foobar", 100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := tf.Format(entry); err != nil {
			b.Fatal(err)
		}
	}
}
