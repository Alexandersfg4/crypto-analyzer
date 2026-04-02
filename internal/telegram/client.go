package telegram

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Alexandersfg4/crypto-analyzer/internal/cron"
	"github.com/Alexandersfg4/crypto-analyzer/internal/report"
	"github.com/Alexandersfg4/crypto-analyzer/internal/storage"
	"github.com/go-telegram/bot"
)

type Client struct {
	b             *bot.Bot
	r             *report.Report
	configStorage *storage.Config
	reportCron    *cron.Cron
}

func New(apiToken string, userID int64, r *report.Report) (*Client, error) {
	opts := []bot.Option{
		bot.WithMiddlewares(auth(userID)),
	}

	b, err := bot.New(apiToken, opts...)
	if nil != err {
		return nil, fmt.Errorf("failed init new bot: %w", err)
	}

	c := &Client{
		b:             b,
		r:             r,
		configStorage: storage.NewConfig("crypto-analyzer.json"),
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommandStartOnly, handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "report", bot.MatchTypeCommandStartOnly, c.handleReport)
	b.RegisterHandler(bot.HandlerTypeMessageText, "tokens", bot.MatchTypeCommandStartOnly, c.handleTokens)
	b.RegisterHandler(bot.HandlerTypeMessageText, "protocols", bot.MatchTypeCommandStartOnly, c.handleProtocols)
	b.RegisterHandler(bot.HandlerTypeMessageText, "config", bot.MatchTypeCommandStartOnly, c.handleConfig)
	b.RegisterHandler(bot.HandlerTypeMessageText, "cron", bot.MatchTypeCommandStartOnly, c.handleCron)
	b.RegisterHandler(bot.HandlerTypeMessageText, "debug", bot.MatchTypeCommandStartOnly, c.handleDebug)

	return c, nil
}

func (c *Client) Start(ctx context.Context) {
	c.b.Start(ctx)
}

// entityPattern describes one Telegram MarkdownV2 inline entity type.
type entityPattern struct {
	re    *regexp.Regexp
	open  string
	close string
}

// allEntityPatterns lists every recognized entity in priority order
// (longer / more-specific delimiters first so they win over shorter ones).
// Each entry's regex is anchored to find the first occurrence in a string.
var allEntityPatterns = []entityPattern{
	// fenced code block — content is completely verbatim
	{re: regexp.MustCompile("(?s)```(?:[\\w]*\\n)?[\\s\\S]*?```"), open: "```", close: "```"},
	// inline code — content is completely verbatim
	{re: regexp.MustCompile("`(?:[^`\\\n]|\\\\.)*`"), open: "`", close: "`"},
	// expandable block-quote opener  **>…
	{re: regexp.MustCompile(`(?m)\*\*>(?:[^\n]*)`), open: "**>", close: ""},
	// block-quote line  >…
	{re: regexp.MustCompile(`(?m)^>(?:[^\n]*)`), open: ">", close: ""},
	// custom emoji / timestamp  ![…](…)   — must come before plain link
	{re: regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`), open: "!", close: ""},
	// inline URL / user mention  […](…)
	{re: regexp.MustCompile(`\[[^\]]*\]\([^)]*\)`), open: "[", close: ""},
	// spoiler  ||…||  — before single | so it wins
	{re: regexp.MustCompile(`\|\|(?:[^|\\\n]|\\.)*\|\|`), open: "||", close: "||"},
	// underline  __…__  — before single _ so it wins
	{re: regexp.MustCompile(`__(?:[^_\\\n]|\\.)*__`), open: "__", close: "__"},
	// bold  *…*
	{re: regexp.MustCompile(`\*(?:[^*\\\n]|\\.)*\*`), open: "*", close: "*"},
	// italic  _…_
	{re: regexp.MustCompile(`_(?:[^_\\\n]|\\.)*_`), open: "_", close: "_"},
	// strikethrough  ~…~
	{re: regexp.MustCompile(`~(?:[^~\\\n]|\\.)*~`), open: "~", close: "~"},
}

// verbatimEntities are entity types whose inner content must never be
// touched (code blocks, URLs, quotes, custom emoji).
// Their content is passed through completely unchanged.
var verbatimEntities = map[string]bool{
	"```": true,
	"`":   true,
	"**>": true,
	">":   true,
	"!":   true,
	"[":   true,
}

// mdV2SpecialChars are all characters that must be escaped in Telegram MarkdownV2
var mdV2SpecialChars = map[rune]bool{
	'*':  true,
	'_':  true,
	'[':  true,
	']':  true,
	'(':  true,
	')':  true,
	'~':  true,
	'`':  true,
	'>':  true,
	'<':  true,
	'#':  true,
	'+':  true,
	'-':  true,
	'=':  true,
	'|':  true,
	'{':  true,
	'}':  true,
	'.':  true,
	'!':  true,
	'\\': true,
}

// processText walks `text`, finds the leftmost / longest matching entity,
// escapes the gap before it, processes the entity (recursing into its inner
// content when appropriate), then continues with the remainder.
func processText(text string) string {
	if text == "" {
		return ""
	}

	// Find the leftmost match among all entity patterns.
	bestStart := -1
	bestEnd := -1
	var bestPat *entityPattern

	for i := range allEntityPatterns {
		p := &allEntityPatterns[i]
		loc := p.re.FindStringIndex(text)
		if loc == nil {
			continue
		}
		if bestStart == -1 || loc[0] < bestStart ||
			(loc[0] == bestStart && (loc[1]-loc[0]) > (bestEnd-bestStart)) {
			bestStart = loc[0]
			bestEnd = loc[1]
			bestPat = p
		}
	}

	if bestPat == nil {
		// No entity found — escape everything.
		return escapeMarkdownV2(text)
	}

	var b strings.Builder

	// Plain text before the entity.
	if bestStart > 0 {
		b.WriteString(escapeMarkdownV2(text[:bestStart]))
	}

	// The matched entity span.
	matched := text[bestStart:bestEnd]

	if verbatimEntities[bestPat.open] {
		// Code blocks, URLs, quotes: pass through completely untouched.
		b.WriteString(matched)
	} else {
		// Inline formatting (bold, italic, underline, strikethrough, spoiler):
		// keep the delimiters and recursively process the inner content so that
		// nested entities survive but stray specials get escaped.
		openLen := len(bestPat.open)
		closeLen := len(bestPat.close)
		inner := matched[openLen : len(matched)-closeLen]

		b.WriteString(bestPat.open)
		b.WriteString(processText(inner))
		b.WriteString(bestPat.close)
	}

	// Continue with the remainder of the string.
	b.WriteString(processText(text[bestEnd:]))

	return b.String()
}

// escapeMarkdownV2 escapes every MarkdownV2 special character in a plain-text
// segment (i.e. a segment that is not part of any recognized entity).
// Already-escaped sequences (backslash + char) are forwarded verbatim to avoid
// double-escaping.
func escapeMarkdownV2(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 8)
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]
		// Forward an existing escape sequence verbatim.
		if ch == '\\' && i+1 < len(runes) {
			b.WriteRune(ch)
			b.WriteRune(runes[i+1])
			i++
			continue
		}
		if mdV2SpecialChars[ch] {
			b.WriteByte('\\')
		}
		b.WriteRune(ch)
	}

	return b.String()
}
