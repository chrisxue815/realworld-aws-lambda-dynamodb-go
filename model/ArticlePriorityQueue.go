package model

import "container/heap"

type ArticlePriorityQueue [][]Article

func (pq ArticlePriorityQueue) Len() int { return len(pq) }

func (pq ArticlePriorityQueue) Less(i, j int) bool {
	// Pop empty lists first, to reduce computation complexity
	if len(pq[i]) == 0 {
		return true
	}
	if len(pq[j]) == 0 {
		return false
	}
	// We want Pop to give us the latest, not earliest, article so we use greater than here.
	return pq[i][0].CreatedAt > pq[j][0].CreatedAt
}

func (pq ArticlePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *ArticlePriorityQueue) Push(x interface{}) {
	item := x.([]Article)
	*pq = append(*pq, item)
}

func (pq *ArticlePriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func MergeArticles(pq ArticlePriorityQueue, offset, limit int) []Article {
	merged := make([]Article, 0, limit)
	heap.Init(&pq)
	numVisitedArticles := 0

	for len(pq) > 0 && numVisitedArticles < offset+limit {
		list := pq[0]

		if len(list) == 0 {
			heap.Pop(&pq)
		} else {
			if numVisitedArticles >= offset {
				article := list[0]
				merged = append(merged, article)
			}
			pq[0] = list[1:]
			heap.Fix(&pq, 0)
			numVisitedArticles++
		}
	}

	return merged
}
