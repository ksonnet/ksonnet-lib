package printer

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/pkg/errors"
)

const (
	space   = byte(' ')
	tab     = byte('\t')
	newline = byte('\n')
	comma   = byte(',')

	syntaxSugar = '+'
)

// Fprint prints a node to the supplied writer using the default
// configuration.
func Fprint(output io.Writer, node ast.Node) error {
	return DefaultConfig.Fprint(output, node)
}

// DefaultConfig is a default configuration.
var DefaultConfig = Config{
	IndentSize: 2,
}

// IndentMode is the indent mode for Config.
type IndentMode int

const (
	// IndentModeSpace indents with spaces.
	IndentModeSpace IndentMode = iota
	// IndentModeTab indents with tabs.
	IndentModeTab
)

// Config is a configuration for the printer.
type Config struct {
	IndentSize int
	IndentMode IndentMode
}

// Fprint prints a node to the supplied writer.
func (c *Config) Fprint(output io.Writer, node ast.Node) error {
	p := printer{cfg: *c}

	p.print(node)

	if p.err != nil {
		return errors.Wrap(p.err, "output")
	}

	_, err := output.Write(p.output)
	return err
}

type printer struct {
	cfg Config

	output      []byte
	indentLevel int

	err error
}

func (p *printer) indent() {
	if len(p.output) == 0 {
		return
	}

	r := p.indentLevel
	var ch byte
	if p.cfg.IndentMode == IndentModeTab {
		ch = tab
	} else {
		ch = space
		r = r * p.cfg.IndentSize
	}

	last := p.output[len(p.output)-1]
	if last == newline {
		pre := bytes.Repeat([]byte{ch}, r)
		p.output = append(p.output, pre...)
	}
}

func (p *printer) writeByte(ch byte, n int) {
	if p.err != nil {
		return
	}

	for i := 0; i < n; i++ {
		p.output = append(p.output, ch)
	}

	p.indent()
}

func (p *printer) writeString(s string) {
	for _, b := range []byte(s) {
		p.writeByte(b, 1)
	}
}

func (p *printer) print(n interface{}) {
	if p.err != nil {
		return
	}

	if n == nil {
		p.err = errors.New("node is nil")
		return
	}

	switch t := n.(type) {
	default:
		p.err = errors.Errorf("unknown node type: %v", n)
		return
	case *ast.Apply:
		p.handleApply(t)
	case ast.Arguments:
		// NOTE: only supporting positional arguments
		for _, arg := range t.Positional {
			p.print(arg)
		}
	case *ast.Array:
		p.writeString("[")
		for i := 0; i < len(t.Elements); i++ {
			p.print(t.Elements[i])

			if i < len(t.Elements)-1 {
				p.writeString(",")
			}
		}
		p.writeString("]")
	case *ast.Binary:
		p.print(t.Left)
		p.writeByte(space, 1)

		p.writeString(t.Op.String())
		p.writeByte(space, 1)

		p.print(t.Right)
	case *ast.Conditional:
		p.writeString("if ")
		p.print(t.Cond)

		p.writeString(" then ")
		p.print(t.BranchTrue)

		if t.BranchFalse != nil {
			p.writeString(" else ")
			p.print(t.BranchFalse)
		}
	case *ast.Index:
		id, err := indexID(t)
		if err != nil {
			p.err = err
			return
		}

		p.writeString(id)
		p.writeString(".")
		p.print(t.Target)
	case *ast.Object:
		p.writeString("{")

		for i, field := range t.Fields {
			if !t.Oneline {
				p.indentLevel++
				p.writeByte(newline, 1)
			}

			p.print(field)
			if i < len(t.Fields)-1 {
				if t.Oneline {
					p.writeByte(comma, 1)
					p.writeByte(space, 1)
				}
			}

			if !t.Oneline {
				p.indentLevel--
				p.writeByte(comma, 1)
			}
		}

		// write an extra newline at the end
		if !t.Oneline {
			p.writeByte(newline, 1)
		}

		p.writeString("}")
	case ast.ObjectField:
		p.handleObjectField(t)
	case *ast.LiteralString:
		switch t.Kind {
		default:
			p.err = errors.Errorf("unknown string literal kind %#v", t.Kind)
			return
		case ast.StringDouble:
			p.writeString(strconv.Quote(t.Value))
		case ast.StringSingle:
			p.writeString(fmt.Sprintf("'%s'", t.Value))
		}

	case *ast.LiteralNumber:
		p.writeString(t.OriginalString)
	case *ast.Self:
		p.writeString("self")
	case *ast.Var:
		p.writeString(string(t.Id))
	}
}

func (p *printer) handleApply(t *ast.Apply) {
	s, err := extractApply(t.Target)
	if err != nil {
		p.err = err
		return
	}

	p.writeString(s)
	p.writeString("(")

	p.print(t.Arguments)
	p.writeString(")")
}

func (p *printer) writeComment(c *ast.Comment) {
	if c == nil {
		return
	}

	lines := strings.Split(c.Text, "\n")
	for _, line := range lines {
		p.writeString("//")
		if len(line) > 0 {
			p.writeByte(space, 1)
		}
		p.writeString(strings.TrimSpace(line))
		p.writeByte(newline, 1)
	}
}

func (p *printer) handleObjectField(t ast.ObjectField) {
	p.writeComment(t.Comment)

	var fieldType string

	switch hide := t.Hide; hide {
	default:
		p.err = errors.Errorf("unknown Hide type %#v", hide)
		return
	case ast.ObjectFieldHidden:
		fieldType = "::"
	case ast.ObjectFieldVisible:
		fieldType = ":::"
	case ast.ObjectFieldInherit:
		fieldType = ":"
	}

	switch kind := t.Kind; kind {
	default:
		p.err = errors.Errorf("unknown Kind type %#v", kind)
		return
	case ast.ObjectFieldID:
		p.writeString(string(*t.Id))
		if t.Method != nil {
			p.addMethodSignature(t)
		}

		if t.SuperSugar {
			p.writeByte(syntaxSugar, 1)
		}

		p.writeString(fieldType)
		p.writeByte(space, 1)
		p.print(t.Expr2)
	case ast.ObjectLocal:
		p.writeString("local ")
		p.writeString(string(*t.Id))
		p.addMethodSignature(t)
		p.writeString(" = ")
		p.print(t.Expr2)
	}
}

func (p *printer) addMethodSignature(of ast.ObjectField) {
	if of.Method == nil {
		return
	}
	params := of.Method.Parameters

	p.writeString("(")
	var args []string
	for _, arg := range params.Required {
		args = append(args, string(arg))
	}

	for _, opt := range params.Optional {
		if opt.DefaultArg != nil {
			var arg string
			arg += string(opt.Name)
			arg += "="

			child := printer{cfg: p.cfg}
			child.print(opt.DefaultArg)
			if child.err != nil {
				p.err = errors.Wrapf(child.err, "invalid argument for %s", string(opt.Name))
				return
			}

			arg += string(child.output)

			args = append(args, arg)
		}
	}

	p.writeString(strings.Join(args, ", "))
	p.writeString(")")
}

func extractApply(n ast.Node) (string, error) {
	switch t := n.(type) {
	default:
		return "", errors.Errorf("invalid type %T when extracting apply", t)
	case *ast.Index:
		var s string
		if t.Target != nil {
			var err error
			s, err = extractApply(t.Target)
			if err != nil {
				return "", errors.Wrap(err, "extract apply")
			}
		}

		id, err := indexID(t)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s.%s", id, s), nil
	case *ast.Var:
		return string(t.Id), nil
	}
}

func literalStringValue(ls *ast.LiteralString) (string, error) {
	if ls == nil {
		return "", errors.New("literal string is nil")
	}

	return ls.Value, nil
}

func indexID(i *ast.Index) (string, error) {
	if i == nil {
		return "", errors.New("index is nil")
	}

	if i.Index != nil {
		ls, ok := i.Index.(*ast.LiteralString)
		if !ok {
			return "", errors.New("index is not a literal string")
		}

		return literalStringValue(ls)
	} else if i.Id != nil {
		return string(*i.Id), nil
	} else {
		return "", errors.New("index and id can't both be blank")
	}
}
