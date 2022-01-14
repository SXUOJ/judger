package cgroup

type CgroupType int

const (
	_ = iota
	CgroupTypeV1
	CgroupTypeV2
)

type Builder struct {
	Type CgroupType

	// CPU：使用调度程序向 cgroup 任务提供对 CPU 的访问
	// CPUAcct：生成有关 cgroup 中任务所使用的 CPU 资源的自动报告
	// CPUSet：为 cgroup 中任务分配单个 CPU（在多核系统中）和内存节点
	// Memory：设置 cgroup 中任务对内存使用的限制，并自动生成有关这些任务使用的内存资源的报告
	CPU     bool
	CPUAcct bool
	CPUSet  bool
	Memory  bool

	// Pids用来限制一个进程可以派生出的进程数量
	Pids bool
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) AddCPU() *Builder     { builder.CPU = true; return builder }
func (builder *Builder) AddCPUAcct() *Builder { builder.CPUAcct = true; return builder }
func (builder *Builder) AddCPUSet() *Builder  { builder.CPUSet = true; return builder }
func (builder *Builder) AddMemory() *Builder  { builder.Memory = true; return builder }
func (builder *Builder) AddPids() *Builder    { builder.Pids = true; return builder }

func (builder *Builder) Build() (Cgroup, error) {
	return nil, nil
}
