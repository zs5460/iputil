// Package iputil 提供了一系列用于处理和合并 IP 地址的工具函数。
package iputil

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
)

type ip struct {
	origin string
	ip     net.IP
	net    *net.IPNet
	classB string
	classC string
	ipInt  uint32
}

// 检查IP是否包含另一个IP
func (i ip) contains(other ip) bool {
	if !i.isCIDR() || other.ip == nil {
		return false
	}
	return i.net.Contains(other.ip)
}

// 获取IP的B段
func (i ip) getClassB() string {
	if i.ip == nil {
		return ""
	}
	return i.classB
}

// 获取IP的C段
func (i ip) getClassC() string {
	if i.ip == nil {
		return ""
	}
	return i.classC
}

// 检查IP是否为CIDR块
func (i ip) isCIDR() bool {
	return i.net != nil
}

type ips []ip

// 添加IP到切片
func (s *ips) append(ip ip) {
	*s = append(*s, ip)
}

// 检查切片是否为空
func (s ips) isEmpty() bool {
	return len(s) == 0
}

func (s ips) Len() int           { return len(s) }
func (s ips) Less(i, j int) bool { return s[i].ipInt < s[j].ipInt }
func (s ips) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s ips) Output() []string {
	result := make([]string, 0, len(s))
	for _, ipAddr := range s {
		result = append(result, ipAddr.origin)
	}
	return result
}

// 单个IP段包含超MaxNumberOfIPs个IP则启动合并
var MaxNumberOfIPs = 10

// 第一步：检测如果有IP段已经包含部分IP,清除已经包含的IP
func step1(s ips) ips {
	if s.isEmpty() {
		return s
	}

	newips := newIPs(len(s))
	shouldKeep := make([]bool, len(s))
	for i := range s {
		shouldKeep[i] = true
	}

	for i, ipAddr := range s {
		if !shouldKeep[i] || !ipAddr.isCIDR() {
			continue
		}

		for j := i + 1; j < len(s) && shouldKeep[j]; j++ {
			if ipAddr.contains(s[j]) {
				shouldKeep[j] = false
			} else {
				break
			}
		}
	}

	for i, ipAddr := range s {
		if shouldKeep[i] {
			newips.append(ipAddr)
		}
	}

	return newips
}

// 将多个同段的IP合并成IPC段
func step2(s ips) ips {
	return mergeIPClass(s, false)
}

// 将多个同段的IP合并成IPB段
func step3(s ips) ips {
	return mergeIPClass(s, true)
}

// 将相邻的C段合并成更大的IP段
// 该函数实现了将连续的C段(/24)网段合并为更大的网段(/23, /22, /21, /20等)
// 合并规则：
// 1. 必须是连续的C段
// 2. 起始地址必须能被合并后的网段大小整除
// 3. 优先尝试最大范围的合并
func step4(s ips) ips {
	// 空切片或只有一个元素时直接返回
	if s.isEmpty() || len(s) == 1 {
		return s
	}

	newips := newIPs(len(s))
	i := 0
	for i < len(s) {
		// 检查是否存在可以合并的下一个网段
		if i+1 < len(s) && canMerge(s[i], s[i+1]) {
			// 计算从当前位置开始有多少个连续的C段
			// 例如：1.1.1.0/24, 1.1.2.0/24, 1.1.3.0/24 就是3个连续段
			count := 2
			for j := i + 2; j < len(s); j++ {
				if !canMerge(s[j-1], s[j]) {
					break
				}
				count++
			}

			// 计算可以合并的最大掩码长度
			// maxPower表示可以合并的位数，比如：
			// maxPower=1 表示可以合并2个C段为/23
			// maxPower=2 表示可以合并4个C段为/22
			// maxPower=3 表示可以合并8个C段为/21
			// maxPower=4 表示可以合并16个C段为/20
			maxPower := 0
			startC := s[i].net.IP[2] // 获取起始地址的C段值
			for power := 1; power <= 8; power++ {
				// 检查两个条件：
				// 1. 连续段数量必须大于等于2^power
				// 2. 起始地址必须能被2^power整除
				// 例如：要合并成/22，需要4个连续段，且起始地址必须是4的倍数
				if count >= (1<<power) && startC%(1<<power) == 0 {
					maxPower = power
				} else {
					break
				}
			}

			// 如果找到了可以合并的掩码长度
			if maxPower > 0 {
				// 计算新的掩码长度并创建合并后的网段
				// 24是C段的掩码长度，减去maxPower得到合并后的掩码长度
				ones, _ := s[i].net.Mask.Size()
				newMaskLen := ones - maxPower
				if mergedIP, err := NewIP(fmt.Sprintf("%s/%d", s[i].net.IP.String(), newMaskLen)); err == nil {
					newips.append(mergedIP)
					// 跳过已经合并的网段
					i += 1 << maxPower
					continue
				}
			}
		}

		// 如果不能合并，保持原样添加到结果中
		newips.append(s[i])
		i++
	}

	return newips
}

func NewIP(s string) (ip ip, err error) {
	if s == "" {
		return ip, errors.New("empty IP address")
	}

	if strings.Contains(s, "/") {
		ip.ip, ip.net, err = net.ParseCIDR(s)
		if err != nil {
			return ip, fmt.Errorf("invalid CIDR: %v", err)
		}
	} else {
		ip.ip = net.ParseIP(s)
		if ip.ip == nil {
			return ip, fmt.Errorf("invalid IP address: %s", s)
		}
		ip.net = nil
	}

	ip.ip = ip.ip.To4()
	if ip.ip == nil {
		return ip, errors.New("not an IPv4 address")
	}

	ip.origin = s
	ip.classB = fmt.Sprintf("%d.%d", ip.ip[0], ip.ip[1])
	ip.classC = fmt.Sprintf("%d.%d.%d", ip.ip[0], ip.ip[1], ip.ip[2])
	ip.ipInt = ipToInt(ip.ip)
	return ip, nil
}

func New(ss []string) (result ips) {
	if len(ss) == 0 {
		return
	}
	result = newIPs(len(ss))
	for _, s := range ss {
		ip, err := NewIP(s)
		if err != nil {
			continue
		}
		result = append(result, ip)
	}
	sort.Sort(result)
	return
}

func Merge(s []string) []string {
	if s == nil {
		return nil
	}
	if len(s) == 0 {
		return []string{}
	}

	ips := New(s)
	ips = step1(ips) // 先处理精确的包含关系
	ips = step4(ips) // 再处理精确的连续网段合并
	ips = step2(ips) // 然后是粗略的C段合并
	ips = step3(ips) // 最后是粗略的B段合并
	ips = step4(ips) // 最后再处理一次连续网段合并
	return ips.Output()
}

// 将IP地址转换为整数
func ipToInt(ip net.IP) uint32 {
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// canMerge 检查两个 CIDR 块是否可以合并
func canMerge(a, b ip) bool {
	if a.ip == nil || b.ip == nil {
		return false
	}
	if a.net == nil || b.net == nil {
		return false
	}
	// 检查两个块是否属于同一个父块
	if a.classB != b.classB {
		return false
	}
	if a.net.Mask.String() != b.net.Mask.String() {
		return false
	}
	if a.ipInt > b.ipInt {
		a, b = b, a
	}
	ones, _ := a.net.Mask.Size()
	return b.ipInt == a.ipInt+uint32(1<<(32-ones))
}

// 预分配指定容量的ips切片
func newIPs(capacity int) ips {
	return make(ips, 0, capacity)
}

// 合并IP段的通用函数
func mergeIPClass(s ips, isClassB bool) ips {
	if s.isEmpty() {
		return s
	}

	// 统计每个段的IP数量
	classCount := make(map[string]int, len(s))
	for _, ipAddr := range s {
		if isClassB {
			classCount[ipAddr.getClassB()]++
		} else {
			classCount[ipAddr.getClassC()]++
		}
	}

	// 预分配切片容量
	newips := newIPs(len(s))
	processedClass := make(map[string]bool, len(classCount))

	for _, ipAddr := range s {
		classKey := ipAddr.getClassC()
		if isClassB {
			classKey = ipAddr.getClassB()
		}

		if processedClass[classKey] || classKey == "" {
			continue
		}

		// 如果段的IP数量 >= MaxNumberOfIPs，合并整个段
		if classCount[classKey] >= MaxNumberOfIPs {
			cidr := classKey + ".0/24"
			if isClassB {
				cidr = classKey + ".0.0/16"
			}
			if mergedIP, err := NewIP(cidr); err == nil {
				newips.append(mergedIP)
				processedClass[classKey] = true
				continue
			}
		}
		newips.append(ipAddr)
	}
	return newips
}
