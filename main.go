package main

import (
	"flag"
	"fmt"
	"math"
	"sort"
)

var (
	// 此处数据可以接入第三方接口，获取实时数据
	maoTai  = target{Name: "贵州茅台", CurrentPrice: 2094.44, Weight: 20}
	meiTuan = target{Name: "美团", CurrentPrice: 261.988, Weight: 40}
	tengXun = target{Name: "腾讯控股", CurrentPrice: 521.303, Weight: 40}
)

var (
	position float64
)

func init() {
	flag.Float64Var(&position, "position", 0, "头寸量")
}

func main() {
	fmt.Println("Start")
	flag.Parse()
	if position == 0 {
		fmt.Println("头寸不可为空")
		return
	}

	arr := []target{maoTai, meiTuan, tengXun}

	resArr := F(arr, position)

	var t float64
	fmt.Println("最优结果为：")
	for k, v := range resArr {
		fmt.Printf("\t%8s\t买入 %d手\n", k.Name, v)
		t += k.CurrentPrice * float64(100*v)
	}
	fmt.Printf("剩余金额: %.3f\nend", position-t)
}

func F(arr []target, position float64) map[target]int {
	resultArr := make(map[target]int, len(arr))
	// 初始化结果
	for _, v := range arr {
		resultArr[v] = 0
	}

	totalWeight := 0
	for _, v := range arr {
		totalWeight += v.Weight
	}

	var t float64
	var success bool //是否分配成功标记
	i := 0
	for {
		i++
		a := float64(i) / float64(arr[0].Weight) / float64(totalWeight)
		t = 0
		var r bool = true // 整数 公倍数标记
		ra := make(map[target]int, len(arr))
		for _, v := range arr {
			b := a * float64(v.Weight) / float64(totalWeight)
			// 非整数倍
			if b != math.Trunc(b) {
				r = false
				break
			}
			t += v.CurrentPrice * float64(100*b)
			ra[v] = int(b)
		}

		if r {
			// fmt.Printf("[DEBUG] t=%.3f\t%v\n", t, ra)
			// fmt.Printf("[DEBUG] position=%.3f\n", position)
			if t > position {
				break
			}
			resultArr = ra
			success = true
		}
	}

	// 如果分配失败，则移除权重最小的标的
	if !success {
		sort.SliceStable(arr, func(i, j int) bool {
			return arr[i].Weight <= arr[j].Weight && arr[i].CurrentPrice >= arr[j].CurrentPrice
		})
		if len(arr) > 1 {
			resA := F(arr[1:], position)
			for k, v := range resA {
				resultArr[k] = v
			}

		}
		return resultArr
	}

	// 如果分配成功，则根据剩余金额，重新选定标的
	t = 0
	for k, v := range resultArr {
		t += k.CurrentPrice * float64(100*v)
	}

	remain := position - t
	if remain > 0 {
		arrB := make([]target, 0)
		for j, _ := range arr {
			if remain >= arr[j].CurrentPrice*float64(100) {
				arrB = append(arrB, arr[j])
			}
		}
		if len(arrB) > 0 {
			resB := F(arrB, remain)
			for k, v := range resB {
				resultArr[k] += v
			}
		}
	}

	return resultArr
}
