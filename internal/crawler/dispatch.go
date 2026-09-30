package crawler

func (c *Crawler) dispatch(submit <-chan task, jobs chan<- task, stop <-chan struct{}) {
	defer close(jobs)

	var queue []task

	for {
		var jobsOut chan<- task
		var next task

		if len(queue) > 0 {
			jobsOut = jobs
			next = queue[0]
		}

		select {
		case t := <-submit:
			queue = append(queue, t)
		case jobsOut <- next:
			queue = queue[1:]
		case <-stop:
			return
		}
	}
}
