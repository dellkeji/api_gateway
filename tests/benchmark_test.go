package tests

import "testing"

func Add(a int, b int) int {
	return a + b
}

func Benchmark_Division(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	b.StartTimer()             // 重新开始计时
	for i := 0; i < b.N; i++ { //use b.N for looping
		Add(4, 5)
	}
}
