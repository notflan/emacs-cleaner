package semaphore

type Semaphore struct {
	channel chan bool
}

func New(max int) *Semaphore {
	this:= new(Semaphore)
	
	this.channel = make(chan bool, max)
	
	return this
}

func (sem* Semaphore) Close() {
	close(sem.channel)
}

func (sem *Semaphore) Lock() {
	sem.channel <- true
}

func (sem *Semaphore) Unlock() {
	<- sem.channel
}
