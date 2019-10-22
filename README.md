# statssimple
Simple thread-safe running time calculation (min/max/avg/total meter)

Suitable for calculating the runtime of a large number of fast calls in different threads,
when the time of collection of more detailed statistics has a significant impact on the results.
 
Typical use is to find the answer to the question "Is it profitable to increase the number of threads".


Подходит для подсчета времени исполнения большого количества быстрых вызовов в разных потоках,
когда время сбора более детальной статистки значимо влияет на результаты.

Типичное использование - поиск ответа на вопрос "выгодно ли увеличивать количество потоков".

Usage:
```go
	testsCnt := 10
	dt := time.Millisecond

	stats := NewStatsSimple() // main instance

	wg := sync.WaitGroup{} // used for this test
	wg.Add(testsCnt)
	
    for cnt := 0; cnt < testsCnt; cnt++ { // run 10 goproc
		go func() { 
			statsn := NewStatsSimple() // alloc local instance
            defer statsn.Shutdown() // shutdown channel & read loop
            for ... { // 
			    statsn.RunOne(func() { // meter function, not thread safe!
				    time.Sleep(dt)
			    })
                // the code below does the same
                statsn.StartOne()
			    time.Sleep(dt)
	            statsn.DoneOne()
            }
			stats.Append(statsn) // appends stats to main instance. thread safe
			wg.Done()
		}()
	}
	wg.Wait() 
	stats.Wait() // to make sure all the data from all the streams are processed
	// stats
    min, max, avg, cnt := stats.GetStatsNs() // get all results
```