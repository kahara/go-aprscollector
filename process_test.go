package main

import (
	"testing"
)

func Test_isSamecall(t *testing.T) {
	type args struct {
		packet []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Not samecall, has SSID",
			args: args{
				packet: []byte("OH3RDA-3>APZMDR,WIDE1-1,qAR,OH5RBG:!6059.29N/02424.74E#PHG2560 R,Wn,Tn Hml"),
			},
			want: false,
		},
		{
			name: "Not samecall, no SSID",
			args: args{
				packet: []byte("OE9ASJ>APBM1D,DMR*,qAS,S52DMR-10:@191209h4725.85N/00945.25Ev255/000OE9ASJ 73"),
			},
			want: false,
		},
		{
			name: "Samecall, has SSID",
			args: args{
				packet: []byte("TA3ILI-7>APBM1D,TA3ILI,DMR*,qAR,TA3ILI:@015857h2300.00N/11300.00E[000/000Op. Ali Mutlu 73"),
			},
			want: true,
		},
		{
			name: "Samecall, no SSID",
			args: args{
				packet: []byte("VK2XB>APBM1D,VK2XB,DMR*,qAR,VK2XB:@191333h3402.23S/15050.44Eu018/000vk2xb"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSamecall(tt.args.packet); got != tt.want {
				t.Errorf("isSamecall() = %v, want %v", got, tt.want)
			}
		})
	}
}
