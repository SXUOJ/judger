package model

type Script struct {
	BaseDir string
	TimeOut uint

	CompileScript string
	RunScript     string
	AfterScript   string
}

var (
	Scripts = map[string]Script{
		"C": {
			BaseDir: "/tmp",
			TimeOut: 5000,

			CompileScript: "/usr/bin/gcc Main.c -o Main",
		},
		"Cpp": {
			BaseDir: "/tmp",
			TimeOut: 5000,

			CompileScript: "/usr/bin/g++ Main.cpp -o Main",
		},
		"Java": {
			BaseDir: "/tmp",
			TimeOut: 5000,

			CompileScript: "",
		},

		"Python2": {
			BaseDir: "/tmp",
			TimeOut: 5000,

			CompileScript: "",
		},

		"Python3": {
			BaseDir: "/tmp",
			TimeOut: 5000,

			CompileScript: "",
		},

		"Go": {
			BaseDir: "/tmp",
			TimeOut: 5000,

			CompileScript: "",
		},
	}
)
