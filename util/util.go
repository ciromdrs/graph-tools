package util

import(
    "fmt"
    "strings"
    "syscall"
)

func GetFileName(path string) string {
    names := strings.Split(path,"/")
    return names[len(names)-1]
}

func GetTime() (int64, int64) {
    var r syscall.Rusage
    err := syscall.Getrusage(syscall.RUSAGE_SELF, &r)
    if err != nil {
        panic(err)
    }
    return r.Utime.Nano(), r.Stime.Nano()
}

func Assert(b bool, msg string){
    if !b {
        fmt.Println("Assertion error:", msg)
    }
}

func ENTER() {
    fmt.Println("Press ENTER")
    fmt.Scanln()
}
