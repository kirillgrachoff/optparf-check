package notify

import (
	"fmt"
	"html"
	"strings"

	"github.com/kirillgrachoff/optparf-check/internal/types"
)

const nbsp = " "

// formatPrice renders 21265.34 as "21 265.34 ₽" with non-breaking spaces.
func formatPrice(price float64) string {
	s := fmt.Sprintf("%.2f", price)
	intPart, frac, _ := strings.Cut(s, ".")
	var b strings.Builder
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteString(nbsp)
		}
		b.WriteRune(c)
	}
	return b.String() + "." + frac + nbsp + "₽"
}

// formatTgMessage renders results as a Telegram HTML message.
func formatTgMessage(results []types.QueryResult) string {
	var b strings.Builder
	for i, r := range results {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "🔎 <b>%s</b>\n", html.EscapeString(r.Pattern))
		switch {
		case r.Err != nil:
			b.WriteString("⚠️ <i>ошибка запроса</i>\n")
		case len(r.Items) == 0:
			b.WriteString("— <i>ничего не найдено</i>\n")
		}
		for _, it := range r.Items {
			name := html.EscapeString(it.Name)
			if it.URL != "" {
				name = fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(it.URL), name)
			}
			fmt.Fprintf(&b, "• %s — <b>%s</b>\n", name, formatPrice(it.Price))
		}
	}
	return b.String()
}
