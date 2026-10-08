package query

import (
	"os"
	"testing"

	"github.com/kirillgrachoff/optparf-check/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	b, err := os.ReadFile("testdata/oceania.html")
	require.NoError(t, err)

	items := NewParser("https://optparf.ru/index.php?q=").Parse(b)

	require.Len(t, items, 3)
	assert.Equal(t, types.Item{
		Id:      "1038790",
		Article: "63998",
		Name:    "Roja Dove Oceania (унисекс) 100ml парфюмерная вода",
		Price:   21265.34,
		URL:     "https://optparf.ru/detail/index.php?ELEMENT_ID=1038790",
	}, items[0])
	assert.Equal(t, 16244.55, items[2].Price)
}

func TestParseEmpty(t *testing.T) {
	assert.Empty(t, NewParser("").Parse([]byte("<div>nothing</div>")))
}
