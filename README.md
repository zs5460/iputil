# IPUtil

[![Go Report Card](https://goreportcard.com/badge/github.com/zs5460/iputil)](https://goreportcard.com/report/github.com/zs5460/iputil)
[![codecov](https://codecov.io/gh/zs5460/iputil/branch/main/graph/badge.svg?token=b7aeunEgyb)](https://codecov.io/gh/zs5460/iputil)
![license](https://img.shields.io/github/license/zs5460/iputil)

IPUtil 是一个高效的 IP 地址合并工具，可以将零散的 IP 地址和网段合并成最优的 CIDR 块。

## 功能特性

- 支持单个 IP 地址和 CIDR 网段的混合输入
- 智能识别并合并被包含的网段
- 自动合并连续的 C 段网段为更大的网段（如 /23, /22, /21, /20）
- 基于 IP 密度的智能 C 段和 B 段合并
- 保证合并后的网段不会包含过多的无效地址

## 使用方法

```go
package main

import (
    "fmt"
    "github.com/zs5460/misc/test/iputil"
)

func main() {
    ips := []string{
        "192.168.1.0/24",
        "192.168.1.2",
        "192.168.1.3",
        "192.168.2.0/24",
        "192.168.3.0/24",
    }
    
    result := iputil.Merge(ips)
    for _, ip := range result {
        fmt.Println(ip)
    }
}

//Output:
//192.168.1.0/24
//192.168.2.0/23
```

## License

See [LICENSE](LICENSE).
