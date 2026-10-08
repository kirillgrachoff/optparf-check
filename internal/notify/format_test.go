package notify

import (
	"errors"
	"testing"

	"github.com/kirillgrachoff/optparf-check/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestFormatPrice(t *testing.T) {
	assert.Equal(t, "21 265.34 ₽", formatPrice(21265.34))
	assert.Equal(t, "382 725.00 ₽", formatPrice(382725))
	assert.Equal(t, "999.90 ₽", formatPrice(999.9))
}

func TestFormatTgMessage(t *testing.T) {
	msg := formatTgMessage([]types.QueryResult{
		{Pattern: "Roja Oceania", Items: []types.Item{
			{Name: "Roja Dove Oceania <100ml>", Price: 21265.34, URL: "https://optparf.ru/detail/index.php?ELEMENT_ID=1&a=b"},
		}},
		{Pattern: "nothing"},
		{Pattern: "broken", Err: errors.New("boom")},
	})
	t.Log("\n" + msg)
	assert.Equal(t, "🔎 <b>Roja Oceania</b>\n"+
		"• <a href=\"https://optparf.ru/detail/index.php?ELEMENT_ID=1&amp;a=b\">Roja Dove Oceania &lt;100ml&gt;</a> — <b>21 265.34 ₽</b>\n"+
		"\n🔎 <b>nothing</b>\n— <i>ничего не найдено</i>\n"+
		"\n🔎 <b>broken</b>\n⚠️ <i>ошибка запроса</i>\n", msg)
}
