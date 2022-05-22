package admin

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Noctember/gocto"
	"os/exec"
	"strings"
)

func AdminExec(ctx *gocto.CommandContext) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", ctx.JoinedArgs())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	stdoutstr := stdout.String()
	var out string

	if err != nil {
		out = fmt.Sprintf("**`ERROR:`** ```\n%s```", err)
	}
	p := gocto.NewPaginatorForContext(ctx)
	p.Delete()
	if stdoutstr != "" {
		p.SetTemplate(func() *gocto.Embed {
			return gocto.NewEmbed().SetAuthor("Nocto Exec")
		})
	}

	scanner := bufio.NewScanner(strings.NewReader(stdoutstr))
	for scanner.Scan() {
		if len(out) > 600 {
			p.AddPage(func(em *gocto.Embed) *gocto.Embed {
				return em.AddField("📤 Out", "```bash\n"+out+"```")
			})
			out = ""
		}
		out += scanner.Text() + "\n"
	}
	if len(out) != 0 {
		p.AddPage(func(em *gocto.Embed) *gocto.Embed {
			return em.AddField("📤 Out", "```bash\n"+out+"```")
		})
	}
	go p.Run()
}
