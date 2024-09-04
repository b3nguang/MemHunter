package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	ps "github.com/mitchellh/go-ps"
)

const (
	PROCESS_VM_READ           = 0x0010
	PROCESS_QUERY_INFORMATION = 0x0400
	MEM_COMMIT                = 0x1000
	PAGE_READWRITE            = 0x04
	PAGE_READONLY             = 0x02
	MEM_PRIVATE               = 0x20000
)

var (
	modkernel32              = syscall.NewLazyDLL("kernel32.dll")
	procVirtualQueryEx       = modkernel32.NewProc("VirtualQueryEx")
	procReadProcessMemory    = modkernel32.NewProc("ReadProcessMemory")
	procGetModuleFileNameExW = modkernel32.NewProc("K32GetModuleFileNameExW")
	wg                       sync.WaitGroup
	mu                       sync.Mutex
)

type MEMORY_BASIC_INFORMATION struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

// VirtualQueryEx function wrapper
func VirtualQueryEx(hProcess syscall.Handle, lpAddress uintptr, lpBuffer unsafe.Pointer, dwLength uintptr) (int, error) {
	ret, _, err := procVirtualQueryEx.Call(uintptr(hProcess), lpAddress, uintptr(lpBuffer), dwLength)
	if ret == 0 {
		return 0, err
	}
	return int(ret), nil
}

// ReadProcessMemory function wrapper
func ReadProcessMemory(hProcess syscall.Handle, lpBaseAddress uintptr, lpBuffer unsafe.Pointer, nSize uintptr, lpNumberOfBytesRead *uintptr) (bool, error) {
	ret, _, err := procReadProcessMemory.Call(uintptr(hProcess), lpBaseAddress, uintptr(lpBuffer), nSize, uintptr(unsafe.Pointer(lpNumberOfBytesRead)))
	if ret == 0 {
		return false, err
	}
	return true, nil
}

// GetProcessFilePath function wrapper
func GetProcessFilePath(hProcess syscall.Handle) (string, error) {
	var filePath [syscall.MAX_PATH]uint16
	ret, _, err := procGetModuleFileNameExW.Call(uintptr(hProcess), 0, uintptr(unsafe.Pointer(&filePath[0])), uintptr(len(filePath)))
	if ret == 0 {
		return "", err
	}
	return syscall.UTF16ToString(filePath[:]), nil
}

func scanProcessMemory(process ps.Process, target string) {
	defer wg.Done()

	pid := process.Pid()

	// 打开进程句柄
	handle, err := syscall.OpenProcess(PROCESS_VM_READ|PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return
	}
	defer syscall.CloseHandle(handle)

	// 获取进程的可执行文件路径
	processFilePath, err := GetProcessFilePath(handle)
	if err != nil {
		return
	}

	var addr uintptr
	var memInfo MEMORY_BASIC_INFORMATION

	// 遍历进程内存块
	for {
		bytesReturned, err := VirtualQueryEx(handle, addr, unsafe.Pointer(&memInfo), unsafe.Sizeof(memInfo))
		if err != nil || memInfo.RegionSize == 0 || bytesReturned == 0 {
			break
		}

		// 匹配内存字符串
		if memInfo.State == MEM_COMMIT && memInfo.Protect&(PAGE_READWRITE|PAGE_READONLY) != 0 {
			buffer := make([]byte, memInfo.RegionSize)
			var bytesRead uintptr
			success, err := ReadProcessMemory(handle, addr, unsafe.Pointer(&buffer[0]), uintptr(len(buffer)), &bytesRead)
			if err == nil && success {
				if idx := strings.Index(string(buffer[:bytesRead]), target); idx != -1 {
					mu.Lock()
					fmt.Printf("[+] 在 PID %d 的地址 0x%x 处找到字符串: %s\n", pid, addr+uintptr(idx), target)
					fmt.Printf("[+] 进程名称: %s\n进程文件路径: %s\n", process.Executable(), processFilePath)
					mu.Unlock()
					break
				}
			}
		}
		addr += memInfo.RegionSize
	}
}

func main() {
	// 获取用户输入字符串
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("[*]请输入检索的字符串: ")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)

	// 获取进程列表
	processes, err := ps.Processes()
	if err != nil {
		fmt.Printf("无法获取进程列表: %v\n", err)
		return
	}

	for _, process := range processes {
		wg.Add(1)
		go scanProcessMemory(process, target)
	}

	wg.Wait()
}
