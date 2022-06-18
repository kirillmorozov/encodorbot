package helpers

import (
	"testing"

	"gopkg.in/telebot.v3"
)

func Test_EscapeMarkdownV2(t *testing.T) {
	type args struct {
		text       string
		entityType telebot.EntityType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty text",
			args: args{text: ``},
			want: ``,
		},
		{
			name: "Escape markdown v2",
			args: args{text: `a_b*c[d]e (fg) h~I` + "`" + `>JK#L+MN -O=|p{qr}s.t!\ \u`},
			want: `a\_b\*c\[d\]e \(fg\) h\~I\` + "`" + `\>JK\#L\+MN \-O\=\|p\{qr\}s\.t\!\\ \\u`,
		},
		{
			name: "Monospaced code block",
			args: args{
				text:       `mono/pre: ` + "`" + `abc` + "`" + ` \int (` + "`" + `\some \` + "`" + `stuff)`,
				entityType: telebot.EntityCodeBlock,
			},
			want: `mono/pre: \` + "`" + `abc\` + "`" + ` \\int (\` + "`" + `\\some \\\` + "`" + `stuff)`,
		},
		{
			name: "Monospaced code",
			args: args{
				text:       `mono/pre: ` + "`" + `abc` + "`" + ` \int (` + "`" + `\some \` + "`" + `stuff)`,
				entityType: telebot.EntityCode,
			},
			want: `mono/pre: \` + "`" + `abc\` + "`" + ` \\int (\` + "`" + `\\some \\\` + "`" + `stuff)`,
		},
		{
			name: "Text link",
			args: args{
				text:       `https://url.containing/funny)cha)\\ra\\)cter\\s`,
				entityType: telebot.EntityTextLink,
			},
			want: `https://url.containing/funny\)cha\)\\\\ra\\\\\)cter\\\\s`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeMarkdownV2(tt.args.text, tt.args.entityType); got != tt.want {
				t.Errorf("escapeMarkdownV2() = %v, want %v", got, tt.want)
			}
		})
	}
}
