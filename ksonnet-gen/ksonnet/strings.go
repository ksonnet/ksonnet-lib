package ksonnet

import (
	"bytes"
	"strings"
	"unicode"
)

func camelCase(in string) string {
	out := ""

	for i, r := range in {
		if i == 0 {
			out += strings.ToLower(string(r))
			continue
		}

		out += string(r)

	}

	return out
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var (
	capTransforms = [][]string{
		{"ISCSI", "Iscsi"},
		{"CIDR", "Cidr"},
		{"HTTP", "Http"},
		{"UUID", "Uuid"},
		{"API", "Api"},
		{"AWS", "Aws"},
		{"CPU", "Cpu"},
		{"GCE", "Gce"},
		{"IPC", "Ipc"},
		{"NFS", "Nfs"},
		{"PID", "Pid"},
		{"RBD", "Rbd"},
		{"TCP", "Tcp"},
		{"TLS", "Tls"},
		{"URI", "Uri"},
		{"URL", "Url"},
		{"WWN", "Wwn"},
		{"FC", "Fc"},
		{"FS", "Fs"},
		{"ID", "Id"},
		{"IO", "Io"},
		{"IP", "Ip"},
		{"SE", "Se"},
	}
)

// capitalizer adjusts the case of terms found in a string.
func toLower(b byte) byte {
	return byte(unicode.ToLower(rune(b)))
}

func isUpper(b byte) bool {
	return unicode.IsUpper(rune(b))
}

// capitalizer adjusts the case of terms found in a string. It will convert `HTTPHeader` into
// `HttpHeader`.
func capitalize(in string) string {
	l := len(in) - 1

	if l == 0 {
		// nothing to do when there is a one character strings
		in = strings.Replace(in, k, v, -1)
		return in
	}

	var b bytes.Buffer
	b.WriteByte(in[0])

	for i := 1; i <= l; i++ {
		if isUpper(in[i-1]) {
			if i < l {
				if isUpper(in[i+1]) || (isUpper(in[i]) && i+1 == l) {
					b.WriteByte(toLower(in[i]))
				} else {
					b.WriteByte(in[i])
				}
			} else if i == l && isUpper(in[i]) {
				b.WriteByte(toLower(in[i]))
			} else {
				b.WriteByte(in[i])
			}
		} else {
			b.WriteByte(in[i])
		}
	}

	return b.String()
}
