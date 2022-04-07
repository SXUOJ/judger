package seccomp

import (
	"syscall"

	libseccomp "github.com/elastic/go-seccomp-bpf"
	"golang.org/x/net/bpf"
)

type Builder struct {
	Allow, Trace []string
	Default      Action
}

// Build builds the filter
func (b *Builder) Build() (Filter, error) {
	policy := libseccomp.Policy{
		DefaultAction: ToSeccompAction(b.Default),
		Syscalls: []libseccomp.SyscallGroup{
			{
				Action: libseccomp.ActionAllow,
				Names:  b.Allow,
			},
			{
				Action: libseccomp.ActionTrace,
				Names:  b.Trace,
			},
		},
	}
	program, err := policy.Assemble()
	if err != nil {
		return nil, err
	}
	return ExportBPF(program)
}

// ExportBPF convert libseccomp filter to kernel readable BPF content
func ExportBPF(filter []bpf.Instruction) (Filter, error) {
	raw, err := bpf.Assemble(filter)
	if err != nil {
		return nil, err
	}
	return sockFilter(raw), nil
}

func sockFilter(raw []bpf.RawInstruction) []syscall.SockFilter {
	filter := make([]syscall.SockFilter, 0, len(raw))
	for _, instruction := range raw {
		filter = append(filter, syscall.SockFilter{
			Code: instruction.Op,
			Jt:   instruction.Jt,
			Jf:   instruction.Jf,
			K:    instruction.K,
		})
	}
	return filter
}
