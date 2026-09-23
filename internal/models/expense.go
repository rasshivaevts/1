package models

import (
	"sync"
	"time"
)

// Expense — структура одной траты
type Expense struct {
	ID          int
	Amount      float64
	Description string
	Date        string
}

// ExpenseStore — хранилище трат в памяти (этап 2)
// Потокобезопасно через sync.Mutex
type ExpenseStore struct {
	mu       sync.Mutex
	nextID   int
	expenses []Expense
}

// NewExpenseStore создаёт новое хранилище
func NewExpenseStore() *ExpenseStore {
	return &ExpenseStore{
		nextID:   1,
		expenses: []Expense{},
	}
}

// Add добавляет новую трату и возвращает её с присвоенным ID
func (s *ExpenseStore) Add(amount float64, description, date string) Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	exp := Expense{
		ID:          s.nextID,
		Amount:      amount,
		Description: description,
		Date:        date,
	}
	s.nextID++
	s.expenses = append(s.expenses, exp)
	return exp
}

// GetAll возвращает все траты
func (s *ExpenseStore) GetAll() []Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]Expense, len(s.expenses))
	copy(result, s.expenses)
	return result
}

// GetByID возвращает трату по ID
func (s *ExpenseStore) GetByID(id int) (Expense, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range s.expenses {
		if e.ID == id {
			return e, true
		}
	}
	return Expense{}, false
}

// Delete удаляет трату по ID
func (s *ExpenseStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range s.expenses {
		if e.ID == id {
			s.expenses = append(s.expenses[:i], s.expenses[i+1:]...)
			return true
		}
	}
	return false
}

// Total возвращает общую сумму трат
func (s *ExpenseStore) Total() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	var sum float64
	for _, e := range s.expenses {
		sum += e.Amount
	}
	return sum
}
