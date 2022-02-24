package main

import (
	"fmt"

	"github.com/Sxu-Online-Judge/judge/model"
	"github.com/Sxu-Online-Judge/judge/sandbox"
)

func main() {
	submit := model.Submit{
		SubmitId:    "1",
		ProblemId:   "1",
		ProblemType: 1,
		CodeType:    "C",
		CodeSource: `#include<stdio.h>
		int main(){
			printf("hello");
			return 0;
		}`,
		Limit: model.Limit{
			TimeLimit:   2000,
			MemoryLimit: 256,
		},
	}

	sandbox := sandbox.StdSandbox{}
	sandbox.TimeLimit = submit.TimeLimit
	sandbox.MemoryLimit = submit.MemoryLimit

	result, _ := sandbox.Run(submit)

	fmt.Println(result)
}
