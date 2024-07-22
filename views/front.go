package views

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func FrontPage() g.Node {
	return Page(
		"Deeler",
		"/",
		H1(g.Text(`Your local deeler 🍃`)),
		P(g.Text(`Choose the right high for your occasion`)),
		P(g.Raw(`Use <em>deeler</em> find strains that suits your mood`)),
	)
}
