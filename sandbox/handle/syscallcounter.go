package handle

type SyscallCounter map[string]int

// NewSyscallCounter creates a new SyscallCounter
func NewSyscallCounter() SyscallCounter {
	return SyscallCounter(make(map[string]int))
}

// Check return inside, allow
func (s SyscallCounter) Check(name string) (bool, bool) {
	n, o := s[name]
	if o {
		s[name] = n - 1
		if n <= 1 {
			return true, false
		}
		return true, true
	}
	return false, true
}
