package calculator

import (
	"errors"
	"sort"
)

var ErrNoPackSizes = errors.New("no pack sizes configured")
var ErrInvalidOrder = errors.New("order quantity must be positive")

type Calculator struct {
	packSizes []int
}

func New(packSizes []int) *Calculator {
	sizes := make([]int, len(packSizes))
	copy(sizes, packSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	return &Calculator{packSizes: sizes}
}

func (c *Calculator) PackSizes() []int {
	return c.packSizes
}

func (c *Calculator) SetPackSizes(sizes []int) {
	c.packSizes = make([]int, len(sizes))
	copy(c.packSizes, sizes)
	sort.Sort(sort.Reverse(sort.IntSlice(c.packSizes)))
}

type Result struct {
	Packs      map[int]int `json:"packs"`
	TotalItems int         `json:"total_items"`
	TotalPacks int         `json:"total_packs"`
}

func (c *Calculator) Calculate(orderQty int) (*Result, error) {
	if len(c.packSizes) == 0 {
		return nil, ErrNoPackSizes
	}
	if orderQty <= 0 {
		return nil, ErrInvalidOrder
	}

	minPack := c.packSizes[len(c.packSizes)-1]
	maxPack := c.packSizes[0]
	upperBound := orderQty + maxPack

	dp := make([]int, upperBound+1)
	parent := make([]int, upperBound+1)

	for i := range dp {
		dp[i] = -1
		parent[i] = -1
	}
	dp[0] = 0

	for i := minPack; i <= upperBound; i++ {
		for _, packSize := range c.packSizes {
			if packSize > i {
				continue
			}
			prev := i - packSize
			if dp[prev] >= 0 {
				newPacks := dp[prev] + 1
				if dp[i] < 0 || newPacks < dp[i] {
					dp[i] = newPacks
					parent[i] = packSize
				}
			}
		}
	}

	bestTotal := -1
	for total := orderQty; total <= upperBound; total++ {
		if dp[total] >= 0 {
			bestTotal = total
			break
		}
	}

	if bestTotal < 0 {
		numPacks := (orderQty + minPack - 1) / minPack
		return &Result{
			Packs:      map[int]int{minPack: numPacks},
			TotalItems: numPacks * minPack,
			TotalPacks: numPacks,
		}, nil
	}

	packs := make(map[int]int)
	remaining := bestTotal
	for remaining > 0 {
		packSize := parent[remaining]
		packs[packSize]++
		remaining -= packSize
	}

	totalPacks := 0
	for _, qty := range packs {
		totalPacks += qty
	}

	return &Result{
		Packs:      packs,
		TotalItems: bestTotal,
		TotalPacks: totalPacks,
	}, nil
}
