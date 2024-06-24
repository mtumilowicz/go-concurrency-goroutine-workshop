# go-concurrency-goroutine-workshop

* references
    * https://www.oreilly.com/library/view/learning-go/9781492077206/
    * https://go.dev/tour/concurrency/5
    * https://en.wikipedia.org/wiki/Communicating_sequential_processes
    * https://www.sciencedirect.com/topics/computer-science/communicating-sequential-process
    * https://slikts.github.io/concurrency-glossary/?id=communicating-sequential-processes-csp
    * https://medium.com/@arturkulig/communicating-sequential-processes-in-js-intro-e620688b6175
    * https://dev.to/karanpratapsingh/csp-vs-actor-model-for-concurrency-1cpg
    * https://aiochan.readthedocs.io/en/latest/csp.html
    * https://chatgpt.com/

## preface
* goals of this workshop
    * introduction to concurrency approach taken by golang
        * Communicating Sequential Processes (CSP)
        * how it differs from actor model
    * understanding of basic components
        * goroutine
        * channel and select
        * context
    * peek into concurrency primitives
        * `WaitGroup`s, mutexes, atomics
* workshop plan
    1. task1
        1. execute two requests in parallel (`fetchCustomer`, `fetchProduct`)
        1. add timeouts than can terminate http requests
        1. return both customer and product to enable further processing
    1. task2
        * encapsulate heavy computation in a function that supports
            1. contextual timeout
            1. declarative timeout
    1. task3
        * divide and conquer: split array to smaller sums, sum them in goroutines and sum partial results

## prerequisite
* Communicating Sequential Processes (CSP)
    * mathematical theory of concurrency based on message passing via channels
        * combines the sequential thread (threaded state) abstraction with synchronous message passing communication
            * collection of independent processes with a distinct memory space
                * example: two processes may run on two cores of the same multiprocessor chip
            * data is exchanged between processors by the sending and receiving of messages
    * gradation
        * processes
            * group of codes that can be considered as an independent entity
            * example: function
        * sequential processes
            * deterministically equivalent to sequential execution
            * it is possible to predict what happens next by knowing the current state
            * disclaimer: it does not mean execution from top to bottom
                * control statements like while, for, etc., are useful because they disrupt the sequential flow
        * communicating sequential processes
            * sequential processes + IO by rendezvous
                * rendezvous = synchronization mechanism where two processes come together to exchange data directly and simultaneously
    * vs actor model
        * actors have identities
            * processes in CSP are anonymous
        * actor model: asynchronous + indirect communication
            * CSP: synchronous + direct communication
        * actor model: no ordering guarantees across different senders
            * CSP messages are delivered in the order they were sent

## components
* goroutine
    * lightweight thread, managed by the Go runtime
    * set of kernel-level threads, each managing a local run queue (LRQ) of goroutines
        * number of threads = `GOMAXPROCS`
        * global run queue (GRQ) for goroutines not assigned yet to a kernel-level thread
        * work stealing
            * blocking operations
                * old thread with its goroutine waiting is descheduled by the OS
                    * LRQ is moved to other thread (new or from the idle pool)
            * balance queues
    * based on Communicating Sequential Processes (CSP)
        * processes communicate with each other by sending and receiving messages through channels
    * Go program starts => Go runtime creates a number of threads and launches a single goroutine to run program
    * stack sizes can grow as needed
    * launched by placing the `go` keyword before a function invocation
    * any values returned by the function are ignored
        * goroutines communicate using channels
    * most of the time has no parameters
        * captures values from the environment
        * if variable might change => pass by parameter
* channel
    * are message boxes — stacks of messages with no specified receiver
    * act as synchronization points
        * reading and writing to a channel are synchronized operations
        * makes them safe even if they run in parallel
    * example
        ```
        ch := make(chan int)
        a := <-ch // read
        ch <- b // write
        ```
    * passing a channel to a function = passing a pointer to the channel
        * similar to how `map` is implemented
    * zero value: `nil`
        * similar to map
        * read from a nil channel never returns
            * not done inside a case in a select statement => program will hang
    * value written to a channel can be read only once
        * in particular: multiple goroutines reading same channel => value read by only one of them
    * convention for passing / assigning
        * `ch <-chan int` - only reads from the channel
        * `ch chan<- int` - only writes to the channel
        * allows the Go compiler to ensure that a channel is only read from or written to by a function
    * default: unbuffered - only one slot
        * similar to promise
            * write => blocks until empty
            * read => blocks until non-empty
        * buffered: `make(chan int, 10)`
            * send operation will block only if the buffer is full
            * use cases
                * gather data back from a set of goroutines
                * limit concurrent usage
                * backpressure
    * read from a channel by using a for-range loop: `for v := range ch`
        * loop continues until the channel is closed or break / return statement is reached
    * `close(ch)`
        * built-in function
        * calling twice => panic
        * write => panic
        * read => always succeeds
            * closed channel always immediately returns its zero value
            * buffered => remaining values will be returned in order
        * required only if a goroutine is reading
            * Go’s runtime can detect channels that are no longer referenced
    * `select`
        * allows a goroutine to read/write to one of a set of multiple channels
            * picks randomly from any of its cases that can go forward
        * blocks until one of its cases can run
        * prevents acquiring locks in an inconsistent order
            * select checks whether any of its cases can proceed
            * every goroutine deadlocked => runtime kills program
                * `fatal error: all goroutines are asleep - deadlock!`
        * for-select loop ~ communicating over a number of channels
            ```
            for {
                select {
                case msg1 := <-ch1:
                    fmt.Println(msg1)
                case msg2 := <-ch2:
                    fmt.Println(msg2)
                }
            }
            ```
        * handling closed channels
            * for-select loop + nil channel
                ```
                case v, ok := <-in:
                    if !ok {
                        in = nil // the case will never succeed again!
                    continue
                ```
* context
    * threadlocal doesn’t work in Go
        * goroutine can be rescheduled to different thread
    * example: `ctx := context.Background()`
    * treated as an immutable instance
        * adding information => wrapping an existing parent context with a child context
        * allows to pass information into deeper layers of the code
            * not the other way around
        * `Value` method checks whether a value is in a context or any of its parent contexts
            * linear search
    * important to make it impossible for context keys to collide
        ```
        type productKey struct{}
        func ContextWithProduct(ctx context.Context, product string) context.Context {
            return context.WithValue(ctx, productKey{}, product)
        }
        func ProductFromContext(ctx context.Context) (string, bool) {
            product, ok := ctx.Value(productKey{}).(string)
            return product, ok
        }
        ```
    * can be cancellable
        * example
            ```
            ctx, cancelFunc := context.WithCancel(context.Background())
            defer cancelFunc() // must be called, otherwise resources leak
            ```
        * tells all the code that’s listening for cancellation that it’s time to stop processing
            * `context.Done` method returns a channel of type `struct{}`
                * empty struct uses no memory
                * channel is closed when the cancel function is invoked
                * returns nil if context not cancellable
            * example
                ```
                for {
                    select {
                    case // reading other channels
                    case <-ctx.Done():
                    }
                }
                ```
            * std HTTP client respects cancellation
            * use cases
                * timeouts
                    ```
                    ctx, cancel := context.WithTimeout(context.Background(), limit) // reaching the timeout cancels the context
                    ```
                * coordinate concurrent goroutines
                    * example: cancel other goroutines when one of them errored
        * cause
            * specifies a cause for the cancellation
            * example
                ```
                ctx, cancelFunc := context.WithCancelCause(context.Background())
                defer cancelFunc(nil) // creating nil cause

                context.Cause(ctx) // reading cause
                ```
            * added in Go 1.20
    * convention: explicitly passed as the first parameter of a function

## primitives
* `WaitGroup`
    * similar to `CountDownLatch` in java
    * use case
        * used to wait for a collection of goroutines to finish executing
        * channel being written by multiple goroutines can be closed only once
* `Once`
    * use case
        * lazy load
        * call some initialization exactly once
* mutex
    * use case
        * sharing access to a field in a struct
        * nearly any other case => use channels
    * not reentrant
        * trying to acquire the same lock twice => deadlock
    * option: readers–writer mutex
        * blocks concurrency when writing
* atomics
    * sync/atomic package
