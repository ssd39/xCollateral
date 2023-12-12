package sender

import "container/heap"

func (i *BroadcastTransactionQueueItem) incrementItemGasFactor() {
	newPayload := i.Provider.SetTransactionGasPrice(i.TransactionData.Transaction, i.GasFactor+1)
	i.TransactionData.Transaction = newPayload
	i.GasFactor = i.GasFactor + 1
}

func (pq BroadcastPriorityQueue) Len() int { return len(pq) }

func (pq BroadcastPriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest nonce so we use less than here.
	return pq[i].Nonce < pq[j].Nonce
}

func (pq BroadcastPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *BroadcastPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*BroadcastTransactionQueueItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *BroadcastPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func ProcessInOutQueue(in <-chan *BroadcastTransactionQueueItem, out chan<- *BroadcastTransactionQueueItem) {
	// Make us a queue!
	pq := make(BroadcastPriorityQueue, 0)
	heap.Init(&pq)

	var currentItem *BroadcastTransactionQueueItem
	var currentIn = in
	var currentOut chan<- *BroadcastTransactionQueueItem

	for {
		select {
		// Read from the input
		case item := <-currentIn:
			// Were we holding something to write? Put it back.
			if currentItem != nil {
				currentItem.Index = len(pq)
				heap.Push(&pq, currentItem)
			}

			// Put our new thing on the queue
			item.Index = len(pq)
			heap.Push(&pq, item)

			currentOut = out
			// Grab our best item. We know there's at least one. We just put it there.
			currentItem = heap.Pop(&pq).(*BroadcastTransactionQueueItem)

			// Write to the output
		case currentOut <- currentItem:
			// OK, we wrote. Is there anything else?
			if len(pq) > 0 {
				// Hold onto it for next time
				currentItem = heap.Pop(&pq).(*BroadcastTransactionQueueItem)
			} else {
				// Turn off the output stream for now until new items
				currentItem = nil
				currentOut = nil
			}
		}
	}
}
