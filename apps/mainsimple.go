package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"math"
	"errors"
	"strconv"
	"sync/atomic"
	"runtime"
	"sync"
	"math/rand"
	"sort"
	"regexp"
	"encoding/json"
	"bytes"
	"bufio"
	"io/ioutil"
	"flag"
)
const (
	win            = 100 // The winning score in a game of Pig
	gamesPerSeries = 10  // The number of games per series to simulate
)

// A score includes scores accumulated in previous turns for each player,
// as well as the points scored by the current player in this turn.
type score struct {
	player, opponent, thisTurn int
}

// An action transitions stochastically to a resulting score.
type action func(current score) (result score, turnIsOver bool)

// roll returns the (result, turnIsOver) outcome of simulating a die roll.
// If the roll value is 1, then thisTurn score is abandoned, and the players'
// roles swap.  Otherwise, the roll value is added to thisTurn.
func roll(s score) (score, bool) {
	outcome := rand.Intn(6) + 1 // A random int in [1, 6]
	if outcome == 1 {
		return score{s.opponent, s.player, 0}, true
	}
	return score{s.player, s.opponent, outcome + s.thisTurn}, false
}

// stay returns the (result, turnIsOver) outcome of staying.
// thisTurn score is added to the player's score, and the players' roles swap.
func stay(s score) (score, bool) {
	return score{s.opponent, s.player + s.thisTurn, 0}, true
}

// A strategy chooses an action for any given score.
type strategy func(score) action

// stayAtK returns a strategy that rolls until thisTurn is at least k, then stays.
func stayAtK(k int) strategy {
	return func(s score) action {
		if s.thisTurn >= k {
			return stay
		}
		return roll
	}
}

// play simulates a Pig game and returns the winner (0 or 1).
func play(strategy0, strategy1 strategy) int {
	strategies := []strategy{strategy0, strategy1}
	var s score
	var turnIsOver bool
	currentPlayer := rand.Intn(2) // Randomly decide who plays first
	for s.player+s.thisTurn < win {
		action := strategies[currentPlayer](s)
		s, turnIsOver = action(s)
		if turnIsOver {
			currentPlayer = (currentPlayer + 1) % 2
		}
	}
	return currentPlayer
}

// roundRobin simulates a series of games between every pair of strategies.
func roundRobin(strategies []strategy) ([]int, int) {
	wins := make([]int, len(strategies))
	for i := 0; i < len(strategies); i++ {
		for j := i + 1; j < len(strategies); j++ {
			for k := 0; k < gamesPerSeries; k++ {
				winner := play(strategies[i], strategies[j])
				if winner == 0 {
					wins[i]++
				} else {
					wins[j]++
				}
			}
		}
	}
	gamesPerStrategy := gamesPerSeries * (len(strategies) - 1) // no self play
	return wins, gamesPerStrategy
}

// ratioString takes a list of integer values and returns a string that lists
// each value and its percentage of the sum of all values.
// e.g., ratios(1, 2, 3) = "1/6 (16.7%), 2/6 (33.3%), 3/6 (50.0%)"
func ratioString(vals ...int) string {
	total := 0
	for _, val := range vals {
		total += val
	}
	s := ""
	for _, val := range vals {
		if s != "" {
			s += ", "
		}
		pct := 100 * float64(val) / float64(total)
		s += fmt.Sprintf("%d/%d (%0.1f%%)", val, total, pct)
	}
	return s
}



// play ping game
func playPigGame(){
	strategies := make([]strategy, win)
	for k := range strategies {
		strategies[k] = stayAtK(k + 1)
	}
	wins, games := roundRobin(strategies)

	for k := range strategies {
		fmt.Printf("Wins, losses staying at k =% 4d: %s\n",
			k+1, ratioString(wins[k], games-wins[k]))
	}
}
//Go’s mechanism for grouping and naming related sets of methods: interfaces.

type user struct {
	name string

	age int                     //Omitted fields will be zero-valued.
}

type geometry interface {
	area() float64
	perim() float64
}

type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) perim()float64 {
	return 2 * math.Pi * c.radius
}

type rect struct {
	width, height float64
}

func (r rect) area() float64 {
	return r.width * r.height
}
/*
Go automatically handles conversion between values and pointers for method calls.
You may want to use a pointer receiver type to avoid copying on method calls
or to allow the method to mutate the receiving struct.
 */
func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func doInterface(){
	r := rect{width: 20, height: 40}
	c := circle{radius : 15}

	measure(r)
	measure(c)
}


func ferror(n int)(int, error){
	if n== 99 {
		return -1, errors.New("can't reach the point 99")
	}
	return n + 3, nil
}

type argError struct {
	arg  int
	prob string
}

func (e *argError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.prob)
}

func f2(arg int) (int, error) {
	if arg == 99 {
		return -1, &argError{arg, "can't work with it"}
	}
	return arg + 3, nil
}

func doError(){

	for _, i := range []int{7, 99} {
		if c, e := ferror(i); e != nil {
			fmt.Println("f1 failed:", e)
		}else {
			fmt.Println("value i ", c )
		}
	}

	for _, i := range []int{99, 42} {
		if r, e := f2(i); e != nil {
			fmt.Println("f2 failed:", e)
		} else {
			fmt.Println("f2 worked:", r)
		}
	}
	_, e := f2(99)
	if ae, ok := e.(*argError); ok {
		fmt.Println(ae.arg)
		fmt.Println(ae.prob)
	}
}

// define the const string variable.
const GOOD_FAQ string= "How do you know that？"

func plus(x, y int) int{
	return x + y
}

func plusPlus(x, y, z int) int {
	return x + y + z
}

func vals()(int, int ){
	return 100, 200
}

func increaseInt()func() int{
	i := 0
	return func() int{
		i += 1
		return i
	}
}

func initNplus2() func()int {
	i := 0
	return  func() int{
		i += 2
		return i
	}
}

//Go supports recursive functions. Here’s a classic factorial example.
func multipleN(n int) int{
	if n== 1 {
		return 1
	}
	return n * multipleN(n -1)
}

func zeroVal(val int){
	val = 1000
}

func zeroOpr(opr *int){
	*opr = 1000
}

func doPoints(){
	i := 0
	zeroVal(i)
	fmt.Println("current zereVal value i= ", i)
	zeroOpr(&i)
	fmt.Println("current zeroptr value i = ", i)

	//struct definition
	fmt.Println(user{name: "John.Lau", age: 29})
	fmt.Println(user{name: "John"})
	s := user{name: "Sean", age: 50}

	u := &s
	fmt.Println("u dereference, name ", u.name)

	//Methods
	r := rect{width: 10, height: 5}
	fmt.Println("area: ", r.area())
	fmt.Println("perim:", r.perim())

	pr := &r
	fmt.Println("pr area:", pr.area())
	fmt.Println("pr perim:", pr.perim())

	//interface
	doInterface()
	doError()
}

func doCloure(){
	nextInt := increaseInt()
	fmt.Println("nextInt", nextInt())
	fmt.Println("nextInt", nextInt())
	fmt.Println("nextInt", nextInt())

	nextPlus2 := initNplus2()
	fmt.Println("nextPlus2: ", nextPlus2())
	fmt.Println("nextPlus2: ", nextPlus2())
	//recursive func
	fmt.Println("recursive func 10*9*8...1 ", multipleN( 10))

}

// Variadic functions can be called with any number of trailing arguments.
// For example, fmt.Println is a common variadic function.\
func sum(nums ...int)int {
	total := 0
	fmt.Println("input", nums)
	for _, v := range nums {
		total += v
	}
	return total
}

func doArrange(){
	num := []int{1 ,3, 5, 7}
	total := 0
	for _, v := range num {
		total += v;
	}
	fmt.Println("sume: ", total)

	//index value
	for i, v := range num {
		if v == 7 {
			fmt.Println("value 7 , index ", i)
		}
	}

	//iterate the map kv
	jobs := map[string]string{"C++": "8k", "Java": "10"}
	for k, v := range jobs {
		fmt.Println("k:", k, ",v:", v)
	}

	//range on strings iterates over Unicode code points. The first value is
	// the starting byte index of the rune and the second the rune itself.
	for i, c := range "go" {
		fmt.Println(i, c)
	}
	s := plus(20, 40)
	z := plusPlus(20, 40, 80)
	fmt.Println("plus(20, 40)= ", s, "plusPlus(20, 40 ,80)=", z)
	a, b := vals()
	fmt.Println("multiple return case:",a, b )
	fmt.Println("sum(10, 100, 1000)", sum(10, 100, 1000))
	nums := []int{1, 10, 100, 1000, 10000}
	//If you already have multiple args in a slice,
	// apply them to a variadic function using func(slice...) like this.
	fmt.Println("sum(nums): ", sum(nums...))

}

func doMaps(){
	m := make(map[string]int)
	m["ok"] = 1
	m["failed"] = 0
	fmt.Println("map[string]int", m)

	m["success"] = 2
	fmt.Println("len:", len(m))

	v := m["ok"]
	fmt.Println("key:ok, value:", v)

	//delete
	delete(m, "ok")
	fmt.Println("Delete kv [\"ok\"]", m)

	//judge it exits or not.
	_, pros := m["success"]
	fmt.Println("result search for a key 'success':", pros)

	//init values
	x := map[string]int{"return": 1, "rssult":0}
	fmt.Println("map init:", x)


}
func doSlice(){
	s:= make([]string, 3)
	fmt.Println("string slice", s)

	s[0] = "0"
	s[2] = "2"
	s[1] = "1"
	fmt.Println("set:", s)
	fmt.Println("get:", s[2])

	fmt.Print("len s", len(s))
	s = append(s, "c")
	fmt.Println("s appended,", s)

	s = append(s, "d", "e")
	fmt.Println("s appended many string", s)

	//copy
	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("copy c :", c)

	//cut arrange
	x := s[2:5]
	fmt.Println("arrange :", x)

	x = s[2:]
	fmt.Println("s[2:]: ", x)

	x = s[:2]
	fmt.Println("s[2:]:", x)

}

func doArray(){
	var a [5]int
	fmt.Println("emp",a )

	//assign value
	a[4] = 100
	fmt.Println("set ", a, ", value: a[4]:", a[4] )

	b := [5]int{12, 13, 14, 15,17}
	fmt.Println("b  value:", b)
	doArrange()
}

func doStr(){
	who := "World!"
	if len(os.Args) > 1 {
		/* os.Args[0] is "hello" or "hello.exe" */
		who = strings.Join(os.Args[1:], " ")
	}
	fmt.Println(who)
	fmt.Println(fmt.Sprintf("print statement: %s", GOOD_FAQ))
	const n = 50000000
	const d = 3e20 / n
	fmt.Println(d)

	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("it's the weekend")
	default:
		fmt.Println("it's a weekday")
	}
	compareStringAppend4Buffer()
	compareStringAppend4ArrayStr()
}

func compareStringAppend4Buffer(){
	var buffer bytes.Buffer
	for i:=0 ; i < 1000; i++ {
		if i % 2 == 0 {
			buffer.WriteString("0")
		} else {
			buffer.WriteString("1")
		}
	}
	fmt.Println("\n"+ buffer.String())
}

func compareStringAppend4ArrayStr(){
	s := []string{}
	for i := 0; i < 1000; i++ {
		if i % 2 == 0 {
			s = append(s, "0")
		} else {
			s = append(s, "1")
		}
	}
	fmt.Println("\n"+ strings.Join(s, ""))
}
func doSimple(){
	doStr()
	doArray()
	doSlice()
	doMaps()
	doCloure()
	doPoints()
}

// A goroutine is a lightweight thread of execution.

func pIndex(form string){
	for i:=0 ; i < 3; i++ {
		fmt.Println(form, ":",  i)
	}
}

func doGo(){
	go func(msg string){
		fmt.Println("msg:", msg)
	}("golang")

	go pIndex("let's go")
	pIndex("show me ")

	//var input string
	//fmt.Scanln(&input)
	//fmt.Println("done")
}
//Channels are the pipes that connect concurrent goroutines.
// You can send values into channels from one goroutine and
// receive those values into another goroutine.
func doChan(){
	message := make(chan string)

	go func(){message <- "ping!"}()

	//This property allowed us to wait at the end of our program for the "ping"
	// message without having to use any other synchronization.
	msg := <- message
	fmt.Println("received msg: ", msg)

	/*
	By default channels are unbuffered, meaning that they will only accept
	 sends (chan <-) if there is a corresponding receive (<- chan) ready
	 to receive the sent value. Buffered channels accept a limited number
	 of values without a corresponding receiver for those values.
	 */
	message2 := make(chan string, 2)
	message2 <- "ping2"
	message2 <- "pong2"
	fmt.Println("received one ", <- message2)
	fmt.Println("received two ", <- message2)

	doChanSync()
	doChanDirection()
}
//This is the function we’ll run in a goroutine. The done channel will be used to
// notify another goroutine that this function’s work is done.
func doChanSync(){
	done := make(chan bool, 1)
	go doAsWorker(done)

	//sync until it has been done.
	<- done
}

func doAsWorker(done chan bool){
	fmt.Println("start to work ...")
	time.Sleep(time.Second)
	fmt.Println("done!")
	done <- true
}

func p(ping chan string, msg string){
	ping <- msg
}

func q(ping chan string, pong chan string){
	msg := <- ping
	pong <- msg
}

//The pong function accepts one channel for receives (pings) and a second for sends (pongs).
func doChanDirection(){
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	p(pings, "go channel direction")
	q(pings, pongs)
	fmt.Println(<-pongs)
}

// Go’s select lets you wait on multiple channel operations.
//Combining goroutines and channels with select is a powerful feature of Go.
func doSelect(){
	t1 := time.Now() // get current time

	x1 := make(chan string )
	x2 := make(chan string)

	go func(){
		time.Sleep(time.Second * 1)
		x1 <- "ping"
	}()

	go func(){
		time.Sleep(time.Second * 2)
		x2 <- "pong"
	}()
	for i:= 0; i < 2; i++ {
		select {
		case msg := <- x1:
			fmt.Println("msg x1:", msg)
		case msg := <- x2:
			fmt.Println("msg x2:", msg)
		}
	}
	elapsed := time.Since(t1)
	//Note that the total execution time is only ~2 seconds since both the 1 and 2 second Sleeps execute concurrently.
	fmt.Println("App elapsed: ", elapsed)
}
//Using this select timeout pattern requires communicating results over channels. This is a good idea in general because
// other important Go features are based on channels and select.
func doTimeout(){
	x1 := make(chan string, 1)
	go func(){
		time.Sleep(time.Second * 1)
		x1 <- "get resource x1"
	}()
	select {
	case msg := <- x1:
		fmt.Println("get result ", msg)
	case <- time.After(time.Second * 1):
		fmt.Println("get timeout after 1 second")
	}

	x2 := make(chan string ,1)
	go func(){
		time.Sleep(time.Second * 2)
		x2 <- "get resource x2"
	}()
	select {
	case msg := <- x2:
		fmt.Println("get result x2:", msg)
	case <- time.After(time.Second * 3):
		fmt.Println("get timeout after 3 seconds")
	}
}

func doNotblockChan(){
	messages := make(chan string)
	signals := make(chan bool)

	select {
	case msg := <- messages:
		fmt.Println("received message :", msg)
	default:
		fmt.Println("no message received")
	}

	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

func closeChan(){
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func(){
		for {
			job, more := <- jobs
			if more {
				fmt.Println("received job", job)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()

	for i :=1; i <= 3; i++ {
		jobs <- i
		fmt.Println("sent msg:", i)
	}
	close(jobs)
	<- done
}

func doChanRange(){
	queue := make(chan string, 2)
	queue <- "two"
	queue <- "three"
	close(queue)
	for elem:= range queue {
		fmt.Println("e in queue :", elem)
	}
}

func doConcurrent(){
	doGo()
	doChan()
	doSelect()
	doTimeout()
	doNotblockChan()
	closeChan()
	doChanRange()
}

//If you just wanted to wait, you could have used time.Sleep.
// One reason a timer may be useful is that you can cancel
// the timer before it expires. Here’s an example of that.
func doTime(){
	timer1 := time.NewTimer(time.Second * 1)
	<- timer1.C
	fmt.Println("Timer 1 expired")

	timer2 := time.NewTimer(time.Second * 2)
	go func(){
		<- timer2.C
		fmt.Println("timer2 expired!")
	}()

	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}
}
// Timers are for when you want to do something once in the future - tickers are for
// when you want to do something repeatedly at regular intervals.
func doTicker(){
	ticker := time.NewTicker(time.Millisecond * 500)

	go func(){
		for t := range ticker.C {
			fmt.Println("ticker transaction :", t.Format("2006-01-02 15:00:00.0000"))
		}
	}()
	time.Sleep(time.Millisecond * 1600)
	ticker.Stop()
	fmt.Println("ticker stopped!")
}

func doTask(id int, jobs <- chan int, results chan <-  int){
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		//time.Sleep(time.Millisecond * 100)
		results <- j * 2
	}
}

func doWorkerPool(){
	jobs := make(chan int , 100)
	results := make(chan int, 100)

	for i := 0; i < 3; i++ {
		go doTask(i, jobs, results)
	}
	for j := 1; j < 10; j++ {
		jobs <- j
	}
	close(jobs)
	for a := 1; a <= 9; a++ {
		fmt.Println("result ", <-results)
	}
	//How to get the each job result.
}

//Rate limiting is an important mechanism for controlling resource utilization and maintaining
// quality of service. Go elegantly supports rate limiting with goroutines, channels, and tickers.
func doRateLimiting(){
	requests := make(chan string, 3)
	for i:=1; i < 4; i++ {
		requests <- strconv.Itoa(i)
	}
	close(requests)
	limiter := time.Tick(time.Millisecond * 200)
	for req := range requests {
		<- limiter
		fmt.Println("req", req, ", time:", time.Now().Local().Format("2006-01-02 15:04:05.0000"))
	}

	//upgrade version
	burstyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(time.Millisecond * 200) {
			burstyLimiter <- t
		}
	}()

	//
	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, ", time ", time.Now().Format("2006-01-02 15:04:05.0000"))
	}
}

func doMutex(){
	state := make(map[int]int)
	var mutex  = &sync.Mutex{}
	var ops uint64 = 0

	for i:=0; i < 100; i++ {
		go func(){
			total := 0
			key := rand.Intn(5)
			mutex.Lock()
			total += state[key]
			mutex.Unlock()
			atomic.AddUint64(&ops, 1)
			runtime.Gosched()

		}()
	}
	for w:=0; w < 10; w++ {
		go func(){
			key := rand.Intn(5)
			val := rand.Intn(10)
			mutex.Lock()
			state[key] = val
			mutex.Unlock()
			atomic.AddUint64(&ops, 1)
			runtime.Gosched()
		}()
	}
	time.Sleep(time.Second)

	mutex.Lock()
	fmt.Println("state :", state)
	mutex.Unlock()
	fmt.Println("show me count :", atomic.LoadUint64(&ops))
}
func doSyncCount(){
	var op uint64 = 0
	for i:= 0; i < 60; i++ {
		go func(){
			atomic.AddUint64(&op, 1)
			runtime.Gosched()
		}()
	}
	time.Sleep(time.Second)
	fmt.Println(atomic.LoadUint64(&op))

	doMutex()
}
//Stateful Goroutines
// This channel-based approach aligns with Go’s ideas of sharing memory
// by communicating and having each piece of data owned by exactly 1 goroutine.

type readOp struct {
	key int
	resp chan int
}

type writeOp struct {
	key int
	val int
	resp chan bool
}
func doStatefulGorutines(){
	var ops uint64 = 0
	reads := make(chan *readOp)
	writes := make(chan *writeOp)

	go func(){
		var state = make(map[int]int)
		for {
			select {
			case read := <- reads:
				read.resp <- state[read.key]
			case write := <- writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()
	// 100 reads
	for r:=0 ; r < 10; r++ {
		go func(){
			for {
				read := &readOp{
					key: rand.Intn(5),
					resp: make(chan int),
				}
				reads <- read
				fmt.Println("get read chan response :", <- read.resp)
				atomic.AddUint64(&ops, 1)
			}
		}()
	}

	//10 writers
	for w :=0; w < 1; w++ {
		go func(){
			for{
				write := &writeOp {
					key: rand.Intn(5),
					val: rand.Intn(1000),
					resp:make(chan bool),
				}
				writes <- write
				<- write.resp
				atomic.AddUint64(&ops, 1)
			}

		}()
	}
	time.Sleep(time.Millisecond * 10)
	fmt.Println("ops:", atomic.LoadUint64(&ops))
}
func doTimer(){
	doTime()
	doTicker()
	doWorkerPool()
	doRateLimiting()
	doSyncCount()
	doStatefulGorutines()
}
func doSimpleSort(){
	strs := []string{"c", "a", "z"}
	sort.Strings(strs)
	fmt.Println("sorting []strings", strs)

	ints := []int {3, 1, 9, 12, 0}
	sort.Ints(ints)


	s := sort.IntsAreSorted(ints)
	fmt.Println("soring []ints", ints, ", sorted :", s)
}

type ByLength []string

func (s ByLength)Len() int{
	return len(s)
}

func (s ByLength)Swap(i, j int) {
	s[i], s[j] = s[j], s[j]
}

func (s ByLength)Less(i, j int)bool {
	return len(s[i]) > len(s[j])
}

func sortByLength(){
	fruits := []string{"peach", "banana", "kiwi"}
	sort.Sort(ByLength(fruits))
	fmt.Println(fruits)
}
func doSort(){
	doSimpleSort()
	sortByLength()
	doPanic()
	doDefer()
	doCollectionFunc()
}

//A panic typically means something went unexpectedly wrong.
func doPanic(){
	//panic("a problem")
	//_, err := os.Create("/tmp/file")
	//if err != nil {
	//	panic(err)
	//}
}

//defer is often used where e.g. ensure and finally would be used in other languages.

func doDefer(){
	f := createFile("./defer.txt")
	defer closeFile(f)
	writeFile(f)
}
func createFile(p string) *os.File {
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}
func writeFile(f *os.File) {
	for i := 1; i <101; i++ {
		if i % 10 == 0 {
			fmt.Fprint(f, "*\n")
		}else {
			fmt.Fprint(f, "*")
		}

	}
}
func closeFile(f *os.File) {
	f.Close()
}

// Go does not support generics; in Go it’s common to provide collection functions
// if and when they are specifically needed for your program and data types.

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func include(vs []string, t string) bool {
	return index(vs, t) > -1
}

func any(vs []string, f func(string)bool) bool {
	for _, v := range vs {
		if f(v){
			return true
		}
	}
	return false
}

func all(vs []string, f func(string)bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func filter(vs []string, f func(string)bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func doMap(vs []string, f func(string)string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

//collect function
func doCollectionFunc(){
	var strs = []string{"peach", "apple", "pear", "plum"}

	fmt.Println("get index 'pear'", index(strs, "pear"))
	fmt.Println("inclue 'grape':", include(strs, "grape"))
	fmt.Println("any hasPrefix ," , any(strs, func(v string )bool{
		return strings.HasPrefix(v, "p")
	}))
	fmt.Println("...", all(strs, func(v string)bool {
		return strings.HasPrefix(v, "p")
	}))
	fmt.Println("...filter", filter(strs,  func(v string) bool {
		return strings.Contains(v, "e")
	}))
	fmt.Println("map string, ", doMap(strs, strings.ToUpper))
	opString()
}


func opString(){
	fmt.Println("Contains:  ", strings.Contains("test", "es"))
	fmt.Println("Count:     ", strings.Count("test", "t"))
	fmt.Println("HasPrefix: ", strings.HasPrefix("test", "te"))
	fmt.Println("HasSuffix: ", strings.HasSuffix("test", "st"))
	fmt.Println("Index:     ", strings.Index("test", "e"))
	fmt.Println("Join:      ", strings.Join([]string{"a", "b"}, "-"))
	fmt.Println("Repeat:    ", strings.Repeat("a", 5))
	fmt.Println("Replace:   ", strings.Replace("foo", "o", "0", -1))
	fmt.Println("Replace:   ", strings.Replace("foo", "o", "0", 1))
	fmt.Println("Split:     ", strings.Split("a-b-c-d-e", "-"))
	fmt.Println("ToLower:   ", strings.ToLower("TEST"))
	fmt.Println("ToUpper:   ", strings.ToUpper("test"))
	opStringFormat()
}

type point struct {
	x, y int
}

func opStringFormat(){
	p := point{1, 2}
	fmt.Printf("%v\n", p)
	//If the value is a struct, the %+v variant will include the struct’s field names.
	fmt.Printf("%+v\n", p)
	//he %#v variant prints a Go syntax representation of the value, i.e. the source code snippet that would produce that value.
	fmt.Printf("%#v\n", p)
	//print the type of a value, use %T.
	fmt.Printf("%T\n", p)
	//Formatting booleans is straight-forward.
	fmt.Printf("%t\n", true)

	//There are many options for formatting integers. Use %d for standard, base-10 formatting.
	fmt.Printf("%d\n", 123)

	//This prints a binary representation.
	fmt.Printf("%b\n", 14)

	//This prints the character corresponding to the given integer.
	fmt.Printf("%c\n", 33)

	//%x provides hex encoding.
	fmt.Printf("%x\n", 456)

	//There are also several formatting options for floats. For basic decimal formatting use %f.
	fmt.Printf("%f\n", 78.9)

	//%e and %E format the float in (slightly different versions of) scientific notation.
	fmt.Printf("%e\n", 123400000.0)
	fmt.Printf("%E\n", 123400000.0)

	//For basic string printing use %s.
	fmt.Printf("%s\n", "\"string\"")

	//To double-quote strings as in Go source, use %q.
	fmt.Printf("%q\n", "\"string\"")

	//As with integers seen earlier, %x renders the string in base-16, with two output characters per byte of input.
	fmt.Printf("%x\n", "hex this")

	//To print a representation of a pointer, use %p.
	fmt.Printf("%p\n", &p)

	//When formatting numbers you will often want to control the width and precision of the resulting figure. To
	// specify the width of an integer, use a number after the % in the verb. By default the result will be
	// right-justified and padded with spaces
	fmt.Printf("|%6d|%6d|\n", 12, 345)

	//You can also specify the width of printed floats, though usually you’ll also want to restrict the
	// decimal precision at the same time with the width.precision syntax
	fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)

	//To left-justify, use the - flag.
	fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45)

	//You may also want to control width when formatting strings, especially to ensure that they
	// align in table-like output. For basic right-justified width.
	fmt.Printf("|%6s|%6s|\n", "foo", "b")

	//To left-justify use the - flag as with numbers.
	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

	//So far we’ve seen Printf, which prints the formatted string to os.Stdout. Sprintf formats
	// and returns a string without printing it anywhere.
	fmt.Println(fmt.Sprintf("a %s", "string"))

	//You can format+print to io.Writers other than os.Stdout using Fprintf.
	fmt.Fprintf(os.Stderr, "an %s\n", "error")
}

// Go offers built-in support for regular expressions. Here are some examples of common regexp-related tasks in Go.
func doMatch(){
	isMatch, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	fmt.Println("go match, ", isMatch)
}

func doCompile(){
	r, _ := regexp.Compile("p([a-z]+)ch")
	fmt.Println("r.FindString", r.FindString("paech punch"))

	fmt.Println("find string index", r.FindStringIndex("peach punch peach peach"))
	fmt.Println("find submatch string", r.FindStringSubmatch("peach punch"))
	fmt.Println("submatch string index", r.FindStringSubmatchIndex("peach punch"))
	fmt.Println("find string, ", r.FindAllString("peach punch sinch", -1))
	fmt.Println(r.Match([]byte("peach")))

	r = regexp.MustCompile("p([a-z]+)ch")

	fmt.Println(r.ReplaceAllString("a peach", "<fruit>"))

	_in := []byte("a peach")
	_out := r.ReplaceAllFunc(_in, func(s []byte) []byte {
		// How can I access the capture group here?
		return []byte(`<a href="/view/PageName">PageName</a>`)
	})
	fmt.Println(string(_out))
	//_out = r.ReplaceAllFunc(_in, bytes.ToUpper)
	//fmt.Println(string(_out))
}
func doRegexp(){
	doMatch()
	doCompile()
}

func doBasicTime(){
	now := time.Now()
	secs := now.Unix()
	nanos := now.UnixNano()
	fmt.Println(now)

	millis := nanos / 1000000
	fmt.Println(secs)
	fmt.Println(millis)
	fmt.Println(nanos)

	//You can also convert integer seconds or nanoseconds
	//since the epoch into the corresponding time.
	fmt.Println(time.Unix(secs, 0))
	fmt.Println(time.Unix(0, nanos))

	doRandom()
}

func doRandom(){
	fmt.Println("rand int 1: ", rand.Intn(100))
	fmt.Println("rand int 2:", rand.Intn(100))

	fmt.Println("rand float 1: ", rand.Float64())
	fmt.Println("rand float 1: ", rand.Float64()* 5 + 5)

	//The default number generator is deterministic, so it’ll produce the same sequence of numbers
	//each time by default. To produce varying sequences, give it a seed that changes. Note that
	//this is not safe to use for random numbers you intend to be secret, use crypto/rand for those.
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	//If you seed a source with the same number, it produces the same sequence of random numbers.
	fmt.Print("same source seed :", r1.Intn(100), ",")
	fmt.Print(r1.Intn(100))

	//same seed has the same sequential
	s2 := rand.NewSource(42)
	r2 := rand.New(s2)
	fmt.Print("1: seed 42:", r2.Intn(100), ",")
	fmt.Print(r2.Intn(100))
	fmt.Println()
	s3 := rand.NewSource(42)
	r3 := rand.New(s3)
	fmt.Print("2: seed 42: ", r3.Intn(100), ",")
	fmt.Print(r3.Intn(100))

}
type Response1 struct {
	Page   int
	Fruits []string
}
type Response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}
func doJson(){
	bValue, _:= json.Marshal(true)
	fmt.Println("bool value:", string(bValue))

	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))

	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]int{"apple": 5, "lettuce": 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	res1D := &Response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	//response 2 json output
	res2D := &Response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))

	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println("unmarshal entity :", dat)

	num := dat["num"].(float64)
	fmt.Println(num)

	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := Response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])

	//json output
	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)

	//time opr
	doBasicTime()
}

func grep(re, filename string) error {
	regex, err := regexp.Compile(re)
	if err != nil {
		return  err// there was a problem with the regular expression.
	}

	fh, err := os.Open(filename)
	f := bufio.NewReader(fh)

	if err != nil {
		return err // there was a problem opening the file.
	}
	defer fh.Close()

	buf := make([]byte, 1024)
	for {
		buf, _ , err = f.ReadLine()
		if err != nil && err.Error() == "EOF" {
			return nil
		}else if err != nil && err.Error() != "EOF" {
			return err
		}

		s := string(buf)
		if regex.MatchString(s) {
			fmt.Printf("%s\n", string(buf))
		}
	}
	return nil
}

func replace(fileName string , re string, replaceTxt string , newFile string) error{
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	output := bytes.Replace(input, []byte(re), []byte( replaceTxt), -1)

	if err = ioutil.WriteFile(newFile, output, 0666); err != nil {
		return err
	}
	return nil
}

func rexexpExample(){
	flag.Parse()
	if flag.NArg() == 2 {
		err := grep(flag.Arg(0), flag.Arg(1))
		if err != nil  {
			fmt.Println(err)
		}
	} else if flag.NArg() == 4 {
		//replace("mainsimple.go", "--", "*", "mainsimpel-updated.go")
		err := replace(flag.Arg(0), flag.Arg(1), flag.Arg(2), flag.Arg(3))
		if err != nil {
			fmt.Println(err)
		}
	} else if flag.NArg() == 3 {
		err := replace(flag.Arg(0), flag.Arg(1), flag.Arg(2), flag.Arg(0))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("Wrong number of arguments. arg1 file name\n")
	}
}

/**
	fmt.Println("*******-args desc***********")
	fmt.Println("args : args == 2 , arg 1: full file path, arg 2: search regexp content")
	fmt.Println("args : args == 4 , arg 1: full file path, arg 2: search regexp content \n "+
	                                          " args 3: replace content, args 4 :  new file name")
	fmt.Println("***********************")
 */
//func main() {
//	//doSimple()
//	//doConcurrent()
//	//doTimer()
//	//doSort()
//	//doRegexp()
//	//doJson()
//	//playPigGame()
//	rexexpExample()
//}







