package judge

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/*
#include <unistd.h>
*/
import "C"

var (
	sc_clk_tck int64
	TICK       time.Duration
	PAGESIZE   int64
)

func init() {
	//timer click number per second
	sc_clk_tck = int64(C.sysconf(C._SC_CLK_TCK))
	TICK = time.Second / time.Duration(sc_clk_tck)
	PAGESIZE = int64(syscall.Getpagesize())
}

func getResourceUsage(pid int) (ok bool, rss int64,
	rt int64, ct int64) {
	stat, err := os.Open("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			return
		}
	}
	bs, err := ioutil.ReadAll(stat)

	if err != nil {
		return
	}
	rss, err = strconv.ParseInt(strings.Split(string(bs), " ")[23], 10, 64)
	rss = rss * PAGESIZE

	// 14 stime 15 utime  TODO: consider cstime cutime
	stime, err := strconv.ParseInt(strings.Split(string(bs), " ")[13], 10, 64)
	if err != nil {
		return
	}
	utime, err := strconv.ParseInt(strings.Split(string(bs), " ")[14], 10, 64)
	if err != nil {
		return
	}
	ct = int64(float64(utime+stime) * 1000 / float64(sc_clk_tck))

	startTime, err := strconv.ParseInt(strings.Split(string(bs), " ")[21], 10, 64)
	if err != nil {
		return
	}
	upTimeFile, err := os.Open("/proc/uptime")
	if err != nil {
		return
	}
	defer upTimeFile.Close()
	bs, err = ioutil.ReadAll(upTimeFile)
	if err != nil {
		return
	}
	//uptime is first paramater in uptime file
	upTime, err := strconv.ParseFloat(strings.Split(string(bs), " ")[0], 64)
	if err != nil {
		return
	}
	rt = int64(upTime*1000) - int64(startTime*1000)/sc_clk_tck
	ok = true
	return

}

// VirtualMemory returns process virtual memory
func virtualMemory(pid int) int64 {
	stat, err := os.Open("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		//panic(err)
	}
	bs, err := ioutil.ReadAll(stat)
	if err != nil {
		//panic(err)
	}
	//virtual memory size is 23nd paramater in the stat file,in bytes
	vmSize, err := strconv.ParseInt(strings.Split(string(bs), " ")[22], 10, 64)

	if err != nil {
		//panic(err)
	}
	return vmSize
}

// RssSize returns process resident memory, but doesn't include swapped out memory
func rssSize(pid int) int64 {
	stat, err := os.Open("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		//panic(err)
	}
	bs, err := ioutil.ReadAll(stat)

	// rss size is 24nd paramater in the stat file,in bytes
	rssSize, err := strconv.ParseInt(strings.Split(string(bs), " ")[23], 10, 64)
	if err != nil {
		//panic(err)
	}
	return rssSize * PAGESIZE

}

// Running returns process total running time from the start
func runningTime(pid int) int64 {
	upTimeFile, err := os.Open("/proc/uptime")
	if err != nil {
		//panic(err)
	}
	defer upTimeFile.Close()
	bs, err := ioutil.ReadAll(upTimeFile)
	if err != nil {
		//panic(err)
	}
	//uptime is first paramater in uptime file
	upTime, err := strconv.ParseFloat(strings.Split(string(bs), " ")[0], 64)
	if err != nil {
		//panic(err)
	}
	stat, err := os.Open("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		//panic(err)
	}
	defer stat.Close()
	bs, err = ioutil.ReadAll(stat)
	if err != nil {
		//panic(err)
	}
	//startTime is 22nd paramater in the stat file
	startTime, err := strconv.ParseInt(strings.Split(string(bs), " ")[21], 10, 64)
	if err != nil {
		//panic(err)
	}
	return int64(upTime*1000) - int64(startTime*1000)/sc_clk_tck
}

// CpuTime returns cpu usage of process in second.
// TODO: dont //panic error
func cpuTime(pid int) int64 {
	stat, err := os.Open("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		//panic(err)
	}
	defer stat.Close()
	bs, err := ioutil.ReadAll(stat)
	if err != nil {
		//panic(err)
	}
	// 14 stime 15 utime  TODO: consider cstime cutime
	stime, err := strconv.ParseInt(strings.Split(string(bs), " ")[13], 10, 64)
	if err != nil {
		//panic(err)
	}
	utime, err := strconv.ParseInt(strings.Split(string(bs), " ")[14], 10, 64)
	if err != nil {
		//panic(err)
	}
	return int64(float64(utime+stime) * 1000 / float64(sc_clk_tck))
}
