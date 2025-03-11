package iputil

import (
	"net"
	"reflect"
	"testing"
)

func TestNewIP(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantIp  ip
		wantErr bool
	}{
		{"testCIDR", args{"1.1.1.0/24"}, ip{"1.1.1.0/24", net.IP{1, 1, 1, 0}, &net.IPNet{}, "1.1", "1.1.1", 16843008}, false},
		{"testip", args{"1.1.1.1"}, ip{"1.1.1.1", net.IP{1, 1, 1, 1}, nil, "1.1", "1.1.1", 16843009}, false},
		{"testerror", args{"1.1.1.1111/24"}, ip{}, true},
		{"testnull", args{"1.1.1.1111"}, ip{}, true},
		{"testnull", args{""}, ip{}, true},
		{"testipv6", args{"::1"}, ip{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIp, err := NewIP(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIp.origin, tt.wantIp.origin) || !reflect.DeepEqual(gotIp.classB, tt.wantIp.classB) || !reflect.DeepEqual(gotIp.ipInt, tt.wantIp.ipInt) {
				t.Errorf("NewIP() = %v, want %v", gotIp, tt.wantIp)
			}
		})
	}
}

func Test_ips_Output(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		want []string
	}{
		{"test1", []string{"1.1.1.2", "1.1.1.1", "1.1.1.19"}, []string{"1.1.1.1", "1.1.1.2", "1.1.1.19"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.s).Output(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ips.Output() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_step1(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"test_step1_1", []string{"1.1.1.3", "1.1.1.0/24", "1.1.1.2"}, []string{"1.1.1.0/24"}},
		{"test_step1_2", []string{"1.1.1.3", "1.1.1.2"}, []string{"1.1.1.2", "1.1.1.3"}},
		{"test_step1_3", []string{"1.1.0.0/16", "1.1.1.2"}, []string{"1.1.0.0/16"}},
		{"test_step1_4", []string{"1.1.0.0/16", "1.1.1.0/24"}, []string{"1.1.0.0/16"}},
		{"test_step1_5", []string{"1.1.1.0/24", "1.1.2.0/24"}, []string{"1.1.1.0/24", "1.1.2.0/24"}},
		{"test_step1_6", []string{"1.1.1.0/24", "1.1.1.0/24"}, []string{"1.1.1.0/24"}},
		{"test_step1_7", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := step1(New(tt.args)).Output(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("step1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_step2(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"test_step2_1", []string{
			"1.1.1.1",
			"1.1.1.2",
			"1.1.1.3",
			"1.1.1.4",
			"1.1.1.5",
			"1.1.1.6",
			"1.1.1.7",
			"1.1.1.8",
			"1.1.1.9",
			"1.1.1.10",
		}, []string{"1.1.1.0/24"}},
		{"test_step2_2", []string{
			"1.1.1.1",
			"1.1.1.2",
			"1.1.1.3",
			"1.1.1.4",
			"1.1.1.5",
			"1.1.1.6",
			"1.1.1.7",
			"1.1.1.8",
			"1.1.1.9",
			"1.1.1.10",
			"1.1.1.12",
			"1.1.1.14",
			"1.1.1.15",
			"1.1.1.16",
			"1.1.1.26",
			"1.1.1.36",
			"1.1.1.46",
			"1.1.1.56",
			"1.1.1.66",
			"1.1.1.76",
			"1.1.1.86",
		}, []string{"1.1.1.0/24"}},
		{"test_step2_3", []string{
			"1.1.1.1",
			"1.1.1.2",
			"1.1.1.4",
			"1.1.1.6",
			"1.1.1.7",
			"1.1.1.8",
			"1.1.1.12",
			"1.1.1.14",
			"1.1.1.15",
			"1.1.1.16",
			"1.1.2.1",
			"1.1.2.2",
			"1.1.2.4",
			"1.1.2.6",
			"1.1.2.7",
			"1.1.2.8",
			"1.1.2.12",
			"1.1.2.14",
			"1.1.2.15",
			"1.1.2.16",
		}, []string{"1.1.1.0/24", "1.1.2.0/24"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := step2(New(tt.args)).Output(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("step2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_step3(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"test_step3_1", []string{
			"1.1.1.1",
			"1.1.1.2",
			"1.1.1.3",
			"1.1.1.4",
			"1.1.1.5",
			"1.1.1.6",
			"1.1.1.7",
			"1.1.1.8",
			"1.1.1.9",
			"1.1.1.10",
		}, []string{"1.1.1.0/24"}},
		{"test_step3_2", []string{
			"1.1.1.1",
			"1.1.1.2",
			"1.1.1.3",
			"1.1.1.4",
			"1.1.1.5",
			"1.1.1.6",
			"1.1.1.7",
			"1.1.1.8",
			"1.1.1.9",
			"1.1.1.10",
			"1.1.1.12",
			"1.1.1.14",
			"1.1.1.15",
			"1.1.1.16",
			"1.1.1.26",
			"1.1.1.36",
			"1.1.1.46",
			"1.1.1.56",
			"1.1.1.66",
			"1.1.1.76",
			"1.1.1.86",
		}, []string{"1.1.1.0/24"}},
		{"test_step3_3", []string{
			"1.1.1.1",
			"1.1.1.2",
			"1.1.1.3",
			"1.1.1.4",
			"1.1.1.5",
			"1.1.1.6",
			"1.1.1.7",
			"1.1.1.8",
			"1.1.1.9",
			"1.1.1.10",
			"1.1.2.1",
			"1.1.2.2",
			"1.1.2.4",
			"1.1.2.6",
			"1.1.2.7",
			"1.1.2.8",
			"1.1.2.12",
			"1.1.2.14",
			"1.1.2.15",
			"1.1.2.16",
		}, []string{"1.1.1.0/24", "1.1.2.0/24"}},
		{"test_step3_4", []string{
			"1.1.1.1",
			"1.1.1.2",
			"1.1.1.3",
			"1.1.1.4",
			"1.1.1.5",
			"1.1.1.6",
			"1.1.1.7",
			"1.1.1.8",
			"1.1.1.9",
			"1.1.1.10",
			"1.1.2.2",
			"1.1.3.3",
			"1.1.4.4",
			"1.1.5.5",
			"1.1.6.6",
			"1.1.7.7",
			"1.1.8.8",
			"1.1.9.9",
			"1.1.10.0/24",
			"1.1.11.0/24",
			"1.1.12.10",
			"1.1.13.10",
			"1.1.14.10",
			"1.1.15.10",
		}, []string{"1.1.0.0/16"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := step3(step2(New(tt.args))).Output(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("step3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_step4(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"test_step4_1", []string{"1.1.0.0/24", "1.1.1.0/24"}, []string{"1.1.0.0/23"}},
		{"test_step4_2", []string{"1.1.0.0/24", "1.1.1.0/24", "1.1.2.0/24", "1.1.3.0/24"}, []string{"1.1.0.0/22"}},
		{"test_step4_3", []string{"1.1.0.0/24", "1.1.1.0/24", "1.1.2.0/24", "1.1.3.0/24", "1.1.4.0/24", "1.1.5.0/24"}, []string{"1.1.0.0/22", "1.1.4.0/23"}},
		{"test_step4_4", []string{"1.1.1.0/24", "1.1.2.0/24", "1.1.3.0/24", "1.1.4.0/24"}, []string{"1.1.1.0/24", "1.1.2.0/23", "1.1.4.0/24"}},
		{"test_step4_5", []string{"1.1.0.0/24", "1.1.1.0/24", "1.1.2.0/24", "1.1.3.0/24", "1.1.4.0/24", "1.1.5.0/24", "1.1.6.0/24", "1.1.7.0/24"}, []string{"1.1.0.0/21"}},
		{"test_step4_6", []string{"1.1.7.0/24", "1.1.8.0/24", "1.1.9.0/24", "1.1.10.0/24", "1.1.11.0/24", "1.1.12.0/24", "1.1.13.0/24", "1.1.14.0/24", "1.1.15.0/24"}, []string{"1.1.7.0/24", "1.1.8.0/21"}},
		{"test_step4_7", []string{
			"1.1.0.0/24", "1.1.1.0/24", "1.1.2.0/24", "1.1.3.0/24",
			"1.1.4.0/24", "1.1.5.0/24", "1.1.6.0/24", "1.1.7.0/24",
			"1.1.8.0/24", "1.1.9.0/24", "1.1.10.0/24", "1.1.11.0/24",
			"1.1.12.0/24", "1.1.13.0/24", "1.1.14.0/24", "1.1.15.0/24",
		}, []string{"1.1.0.0/20"}},
		{"test_step4_8", []string{
			"1.1.3.0/24", "1.1.4.0/24", "1.1.5.0/24", "1.1.6.0/24",
			"1.1.7.0/24", "1.1.8.0/24", "1.1.9.0/24", "1.1.10.0/24",
		}, []string{"1.1.3.0/24", "1.1.4.0/22", "1.1.8.0/23", "1.1.10.0/24"}},
		{"test_step4_9", []string{
			"1.1.16.0/24", "1.1.17.0/24", "1.1.18.0/24", "1.1.19.0/24",
			"1.1.20.0/24", "1.1.21.0/24", "1.1.22.0/24", "1.1.23.0/24",
			"1.1.24.0/24", "1.1.25.0/24", "1.1.26.0/24", "1.1.27.0/24",
			"1.1.28.0/24", "1.1.29.0/24", "1.1.30.0/24", "1.1.31.0/24",
		}, []string{"1.1.16.0/20"}},
		{"test_step4_10", []string{
			"1.1.15.0/24",
			"1.1.16.0/24",
			"1.1.17.0/24",
			"1.1.18.0/24",
			"1.1.19.0/24",
			"1.1.20.0/24",
			"1.1.21.0/24",
			"1.1.22.0/24",
			"1.1.23.0/24",
		}, []string{"1.1.15.0/24", "1.1.16.0/21"}},
		{"test_step4_11", []string{"1.1.0.0/24", "1.1.1.0/24", "1.1.2.0/24", "1.1.3.0/24", "1.2.0.0/24", "1.3.0.0/24", "1.3.1.0/24", "1.3.2.0/24", "1.3.3.0/24"}, []string{"1.1.0.0/22", "1.2.0.0/24", "1.3.0.0/22"}},
		{"test_step4_12", []string{"1.1.0.0/23", "1.1.2.0/23"}, []string{"1.1.0.0/22"}},
		{"test_step4_13", []string{"1.1.0.0/22", "1.1.4.0/22"}, []string{"1.1.0.0/21"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := step4(New(tt.args)).Output(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("step4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_canMerge(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test_canMerge_1", args{"1.1.0.0/24", "1.1.1.0/24"}, true},
		{"test_canMerge_2", args{"1.1.0.0/23", "1.1.2.0/23"}, true},
		{"test_canMerge_3", args{"1.1.0.0/24", "1.1.2.0/23"}, false},
		{"test_canMerge_4", args{"1.1.0.0/23", "1.1.2.0/24"}, false},
		{"test_canMerge_5", args{"1.1.0.0/24", "1.1.2.0/24"}, false},
		{"test_canMerge_6", args{"1.1.0.0/24", "1.1.1.1"}, false},
		{"test_canMerge_7", args{"1.1.2.0/24", "1.1.0.0/24"}, false},
		{"test_canMerge_6", args{"1.1.0.0/24", "1.1.1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip1, _ := NewIP(tt.args.a)
			ip2, _ := NewIP(tt.args.b)
			if got := canMerge(ip1, ip2); got != tt.want {
				t.Errorf("canMerge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ip_contains(t *testing.T) {
	tests := []struct {
		name  string
		i     string
		other string
		want  bool
	}{
		{"Test_ip_contains_1", "1.1.1.1", "1.1.1.2", false},
		{"Test_ip_contains_2", "1.1.1.0/24", "1.1.1.2", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip1, _ := NewIP(tt.i)
			ip2, _ := NewIP(tt.other)
			if got := ip1.contains(ip2); got != tt.want {
				t.Errorf("ip.contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ipToInt(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want uint32
	}{
		{"Test_ipToInt_1", nil, 0},
		{"Test_ipToInt_2", net.IP{}, 0},
		{"Test_ipToInt_3", net.IP{1, 1, 1, 1}, 16843009},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ipToInt(tt.ip); got != tt.want {
				t.Errorf("ipToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"TestMerge_1", nil, nil},
		{"TestMerge_2", []string{}, []string{}},
		{"TestMerge_3", []string{"1.1.1.1"}, []string{"1.1.1.1"}},
		{"TestMerge_4", []string{"1.1.1.0/24", "1.1.1.2"}, []string{"1.1.1.0/24"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
