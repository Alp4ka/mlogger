package misc

type ParseMode string

const (
	MarkdownMode ParseMode = "markdown"
	HtmlMode     ParseMode = "html"

	DefaultMode = MarkdownMode
)
