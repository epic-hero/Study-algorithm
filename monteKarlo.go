package main

import "fmt"
import "math"
import "sync"
import "math/rand"

func main() {
	mc := input()
	integral := mc.startAlgorithm()
	fmt.Println(integral)
}

func input() MonteCarlo {
	return newMonteCarlo( 0, 4, 0, 1, 5000)
}

func (mc *MonteCarlo) startAlgorithm() float64 {
	wg := new(sync.WaitGroup)
	for i := 1; i <= mc.N; i++ {
		wg.Add(1)
		go mc.genPoints(i, wg)
	}
	wg.Wait()
	return mc.getIntegral()
}

type Point struct {
	X      float64
	Y      float64
	Fx     float64
	isDown bool
}

type MonteCarlo struct {
	Xmin float64
	Xmax float64
	Ymin float64
	Ymax float64
	N    int

	PointsMutex *sync.Mutex
	Points      []Point
}

func newMonteCarlo(Xmin float64, Xmax float64, Ymin float64, Ymax float64, N int) MonteCarlo {
	return MonteCarlo{
		Xmin:        Xmin,
		Xmax:        Xmax,
		Ymin:        Ymin,
		Ymax:        Ymax,
		N:           N,
		Points:      make([]Point, 0),
		PointsMutex: new(sync.Mutex),
	}
}

func (mc *MonteCarlo) genPoints(value int, wg *sync.WaitGroup) {
	defer wg.Done()
	mc.PointsMutex.Lock()
	defer mc.PointsMutex.Unlock()

	Xcoord := mc.Xcoord()
	Ycoord := mc.Ycoord()
	Fx := f(Xcoord)
	isDown := false
	if Ycoord <= Fx {
		isDown = true
	}

	point := Point{
		Y:      Ycoord,
		X:      Xcoord,
		Fx:     Fx,
		isDown: isDown,
	}
	mc.Points = append(mc.Points, point)
}

func random() float64 {
	return rand.Float64()
}

func (mc *MonteCarlo) S() float64 {
	return (mc.Xmax - mc.Xmin) * (mc.Ymax - mc.Ymin)
}

func (mc *MonteCarlo) Xcoord() float64 {
	return mc.Xmin + random()*(mc.Xmax-mc.Xmin)
}

func (mc *MonteCarlo) Ycoord() float64 {
	return mc.Ymin + random()*(mc.Ymax-mc.Ymin)
}

func f(x float64) float64 {
	return math.Sin(x) / x
}

func (mc MonteCarlo) getIntegral() float64 {
	sumUp := 0.00
	sumDown := 0.00
	for _, value := range mc.Points {
		if value.isDown == false {
			sumUp++
		} else if value.isDown == true {
			sumDown++
		}
	}
	return (sumUp / (sumUp + sumDown)) * mc.S()
}
