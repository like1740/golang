package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	/*
		指针：
		1.编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
	*/
	num := 10
	fmt.Println("修改前的值:", num)
	addTen(&num)
	fmt.Println("修改后的值:", num)

	/*
		指针：
		2.实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
	*/
	numbers := []int{1, 2, 3, 4, 5}
	doubleSlice(&numbers)
	fmt.Println("修改后的切片:", numbers)

	/*
		Goroutine:
		1.编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	*/
	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Printf("奇数: %d\n", i)
		}
	}()
	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Printf("偶数: %d\n", i)
		}
	}()
	time.Sleep(1 * time.Second) // 主协程等待足够时间让子协程完成

	/*
		Goroutine:
		2.设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
	*/
	scheduler := NewTaskScheduler()

	scheduler.AddTask(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Task 1 completed")
	})

	scheduler.AddTask(func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Task 2 completed")
	})

	scheduler.AddTask(func() {
		time.Sleep(1500 * time.Millisecond)
		fmt.Println("Task 3 completed")
	})

	scheduler.AddTask(func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Task 4 completed")
	})

	results := scheduler.Run()

	fmt.Println("\nTask Execution Times:")
	for _, result := range results {
		fmt.Printf("Task %d 执行时间 %v\n", result.ID+1, result.Duration.Round(time.Millisecond))
	}

	totalTime := results[len(results)-1].EndTime.Sub(results[0].StartTime)
	fmt.Printf("\n任务总执行时间: %v\n", totalTime.Round(time.Millisecond))

	/*
		面向对象:
		1.定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
	*/
	rect := Rectangle{Width: 5, Height: 3} // 创建 Rectangle 实例
	circle := Circle{Radius: 4}            // 创建 Circle 实例

	shapes := []Shape{rect, circle}
	for _, shape := range shapes {
		printShapeInfo(shape)
	}

	/*
		面向对象:
		2.使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
	*/
	emp := Employee{
		Person: Person{
			Name: "张三",
			Age:  28,
		},
		EmployeeID: "1",
	}
	emp.PrintInfo()

	/*
		Channel:
		1.编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
	*/
	ch1 := make(chan int)

	var wg1 sync.WaitGroup
	wg1.Add(2)

	go producer1(ch1, &wg1)
	go consumer1(ch1, &wg1)

	wg1.Wait()
	fmt.Println("所有数字处理完毕！")

	/*
		Channel:
		2.实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
	*/
	ch2 := make(chan int, BUFFER_SIZE)

	var wg2 sync.WaitGroup
	wg2.Add(2)

	go producer2(ch2, &wg2)
	go consumer2(ch2, &wg2)

	wg2.Wait()
	fmt.Println("\n所有数字处理完毕！")

	/*
		锁机制:
		1.编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	*/
	counter := Counter{}

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.mu.Lock()
				counter.count++
				counter.mu.Unlock()
			}
		}()
	}

	wg.Wait()

	fmt.Println("Final counter:", counter.count)

	/*
		锁机制:
		2.使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	*/
	var counter3 int64

	var wg3 sync.WaitGroup

	numRoutines := 10
	incrementsPerRoutine := 1000

	wg3.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(id int) {
			defer wg3.Done()

			for j := 0; j < incrementsPerRoutine; j++ {
				atomic.AddInt64(&counter3, 1)
			}
		}(i)
	}
	wg3.Wait()

	finalValue := atomic.LoadInt64(&counter3)
	fmt.Printf("\nFinal counter value: %d\n", finalValue)

}

func addTen(ptr *int) {
	*ptr += 10
}

func doubleSlice(slicePtr *[]int) {
	slice := *slicePtr
	for i := range slice {
		slice[i] *= 2
	}
}

type Task func()

type TaskResult struct {
	ID        int
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

type TaskScheduler struct {
	tasks     []Task
	results   []TaskResult
	wg        sync.WaitGroup
	mu        sync.Mutex
	startTime time.Time
	endTime   time.Time
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks:   make([]Task, 0),
		results: make([]TaskResult, 0),
	}
}

func (ts *TaskScheduler) AddTask(task Task) {
	ts.tasks = append(ts.tasks, task)
}

func (ts *TaskScheduler) Run() []TaskResult {
	ts.startTime = time.Now()
	ts.results = make([]TaskResult, len(ts.tasks))

	ts.wg.Add(len(ts.tasks))

	for i, task := range ts.tasks {
		go ts.executeTask(i, task)
	}

	ts.wg.Wait()
	ts.endTime = time.Now()

	return ts.results
}

func (ts *TaskScheduler) executeTask(id int, task Task) {
	defer ts.wg.Done()

	result := TaskResult{
		ID:        id,
		StartTime: time.Now(),
	}

	task()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	ts.mu.Lock()
	ts.results[id] = result
	ts.mu.Unlock()
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func printShapeInfo(s Shape) {
	var shapeName string
	switch s.(type) {
	case Rectangle:
		shapeName = "Rectangle"
	case Circle:
		shapeName = "Circle"
	default:
		shapeName = "Unknown Shape"
	}

	fmt.Printf("=== %s ===\n", shapeName)
	fmt.Printf("Area:      %.2f\n", s.Area())
	fmt.Printf("Perimeter: %.2f\n\n", s.Perimeter())
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person     Person
	EmployeeID string
	Name       string
}

func (e Employee) PrintInfo() {
	fmt.Println("=== 员工信息 ===")
	fmt.Printf("姓名:      %s\n", e.Person.Name)
}

func producer1(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)

	fmt.Println("生产者开始生成数字...")

	for i := 1; i <= 10; i++ {
		time.Sleep(100 * time.Millisecond)

		ch <- i
		fmt.Printf("生产者发送: %d\n", i)
	}

	fmt.Println("生产者已完成所有数字生成")
}

func consumer1(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("消费者准备接收数字...")

	for num := range ch {
		time.Sleep(200 * time.Millisecond)

		fmt.Printf("消费者接收: %d\n", num)
	}

	fmt.Println("消费者已处理所有数字")
}

const (
	BUFFER_SIZE  = 10
	TOTAL_NUMS   = 100
	PRODUCE_TIME = 50
	CONSUME_TIME = 100
)

func producer2(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)

	fmt.Printf("生产者开始生成 %d 个数字 (缓冲区大小: %d)...\n", TOTAL_NUMS, BUFFER_SIZE)

	for i := 1; i <= TOTAL_NUMS; i++ {
		time.Sleep(time.Duration(PRODUCE_TIME) * time.Millisecond)

		ch <- i
	}
}

func consumer2(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("消费者准备接收 %d 个数字...\n", TOTAL_NUMS)

	count := 0
	for num := range ch {
		time.Sleep(time.Duration(CONSUME_TIME) * time.Millisecond)

		count++
		fmt.Printf("消费: %3d | 缓冲区: %2d/%d\n", num, len(ch), cap(ch))
	}
}

type Counter struct {
	mu    sync.Mutex
	count int
}
