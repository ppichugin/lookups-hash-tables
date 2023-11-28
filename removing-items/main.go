package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Employee struct {
	name    string
	phone   string
	deleted bool
}

func main() {
	// Make some names.
	employees := []Employee{
		Employee{"Ann Archer", "202-555-0101", false},
		Employee{"Bob Baker", "202-555-0102", false},
		Employee{"Cindy Cant", "202-555-0103", false},
		Employee{"Dan Deever", "202-555-0104", false},
		Employee{"Edwina Eager", "202-555-0105", false},
		Employee{"Fred Franklin", "202-555-0106", false},
		Employee{"Gina Gable", "202-555-0107", false},
	}

	hashTable := NewLinearProbingHashTable(10)
	for _, employee := range employees {
		hashTable.set(employee.name, employee.phone)
	}
	hashTable.dump()

	hashTable.probe("Hank Hardy")
	fmt.Printf("Table contains Sally Owens: %t\n", hashTable.contains("Sally Owens"))
	fmt.Printf("Table contains Dan Deever: %t\n", hashTable.contains("Dan Deever"))
	fmt.Println("Deleting Dan Deever")
	hashTable.delete("Dan Deever")
	fmt.Printf("Table contains Dan Deever: %t\n", hashTable.contains("Dan Deever"))
	fmt.Printf("Sally Owens: %s\n", hashTable.get("Sally Owens"))
	fmt.Printf("Fred Franklin: %s\n", hashTable.get("Fred Franklin"))
	fmt.Println("Changing Fred Franklin")
	hashTable.set("Fred Franklin", "202-555-0100")
	fmt.Printf("Fred Franklin: %s\n", hashTable.get("Fred Franklin"))
	hashTable.dump()

	hashTable.probe("Ann Archer")
	hashTable.probe("Bob Baker")
	hashTable.probe("Cindy Cant")
	hashTable.probe("Dan Deever")
	hashTable.probe("Edwina Eager")
	hashTable.probe("Fred Franklin")
	hashTable.probe("Gina Gable")
	hashTable.set("Hank Hardy", "202-555-0108")
	hashTable.probe("Hank Hardy")

	// Look at clustering.
	fmt.Println(time.Now())                   // Print the time so it will compile if we use a fixed seed.
	random := rand.New(rand.NewSource(12345)) // Initialize with a fixed seed
	// random := rand.New(rand.NewSource(time.Now().UnixNano())) // Initialize with a changing seed
	bigCapacity := 1009
	bigHashTable := NewLinearProbingHashTable(bigCapacity)
	numItems := int(float32(bigCapacity) * 0.9)
	for i := 0; i < numItems; i++ {
		str := fmt.Sprintf("%d-%d", i, random.Intn(1000000))
		bigHashTable.set(str, str)
	}
	bigHashTable.dumpConcise()
	fmt.Printf("Average probe sequence length: %f\n",
		bigHashTable.aveProbeSequenceLength())
}

// djb2 hash function. See http://www.cse.yorku.ca/~oz/hash.html.
func hash(value string) int {
	hash := 5381
	for _, ch := range value {
		hash = ((hash << 5) + hash) + int(ch)
	}

	// Make sure the result is non-negative.
	if hash < 0 {
		hash = -hash
	}
	return hash
}

type LinearProbingHashTable struct {
	capacity  int
	employees []*Employee
}

// NewLinearProbingHashTable Initialize a LinearProbingHashTable and return a pointer to it.
func NewLinearProbingHashTable(capacity int) *LinearProbingHashTable {
	// Create a new LinearProbingHashTable
	table := &LinearProbingHashTable{
		capacity: capacity,
		// Allocate the slice of employees
		employees: make([]*Employee, capacity),
	}

	// Return the pointer to the new LinearProbingHashTable
	return table
}

// Display the hash table's contents.
func (hashTable *LinearProbingHashTable) dump() {
	for i, employee := range hashTable.employees {
		if employee == nil {
			fmt.Printf("%d: ---\n", i)
		} else if employee.deleted {
			fmt.Printf("%d: xxx\n", i)
		} else {
			fmt.Printf("%d: %s\t%s\n", i, employee.name, employee.phone)
		}
	}
}

// Return the key's index or where it would be if present and
// the probe sequence length.
// If the key is not present and the table is full, return -1 for the index.
func (hashTable *LinearProbingHashTable) find(name string) (int, int) {
	// Calculate the hash of the name
	hash := hash(name) % hashTable.capacity

	// Set deletedIndex to -1. This will be the index of the first deleted item we come across (if we find one).
	deletedIndex := -1

	// Follow the probe sequence
	for i := 0; i < hashTable.capacity; i++ {
		// Calculate a position in the probe sequence
		index := (hash + i) % hashTable.capacity

		// If this spot is empty, then the target is not in the table.
		if hashTable.employees[index] == nil {
			// If deletedIndex is greater than or equal to 0, then we found a deleted item earlier. Return its index.
			if deletedIndex >= 0 {
				return deletedIndex, i + 1
			}
			// Otherwise, we did not find a deleted item, so return the current index for the nil entry.
			return index, i + 1
		}

		// (At this point, the current spot is not nil.) If deletedIndex is still -1 and this spot is deleted, update deletedIndex to save this index.
		if deletedIndex == -1 && hashTable.employees[index].deleted {
			deletedIndex = index
		}

		// Otherwise, if this spot contains the target, return this index.
		if hashTable.employees[index].name == name {
			return index, i + 1
		}
	}

	// After the loop ends if we have not returned yet, then the key is not in the table and the table is full.
	// If deletedIndex is greater than or equal to 0, then we found a deleted entry. Return that index.
	if deletedIndex >= 0 {
		return deletedIndex, hashTable.capacity
	}

	// Otherwise, the table is full, the target is not present, there are no deleted spots to reuse, and the world seems like a generally harsh and unforgiving place. Don’t despair! Just return -1.
	return -1, hashTable.capacity
}

// Add an item to the hash table.
func (hashTable *LinearProbingHashTable) set(name string, phone string) {
	// Call find to get the index where the key belongs
	index, _ := hashTable.find(name)

	// If the index is less than 0, the key is not in the table and the table is full
	if index < 0 {
		panic("Hash table is full")
	}

	// If the slice entry at the key's index is nil or deleted, create a new Employee struct
	if hashTable.employees[index] == nil || hashTable.employees[index].deleted {
		hashTable.employees[index] = &Employee{name: name, phone: phone, deleted: false}
	} else {
		// Otherwise, find found the target key. Update its value.
		hashTable.employees[index].phone = phone
	}
}

// Return an item from the hash table.
func (hashTable *LinearProbingHashTable) get(name string) string {
	// Call find to get the index where the key belongs
	index, _ := hashTable.find(name)

	// If the returned index is less than 0, then the target is not present and the hash table is full. Return a blank string.
	if index < 0 {
		return ""
	}

	// Else if the returned spot is nil, return a blank string.
	if hashTable.employees[index] == nil {
		return ""
	}

	// Else if the returned spot is deleted, return a blank string.
	if hashTable.employees[index].deleted {
		return ""
	}

	// Else return the found item’s value.
	return hashTable.employees[index].phone
}

// Return true if the person is in the hash table.
func (hashTable *LinearProbingHashTable) contains(name string) bool {
	// Call find to get the index where the key belongs
	index, _ := hashTable.find(name)

	// If the returned index is less than 0, then the target is not present and the hash table is full. Return false.
	if index < 0 {
		return false
	}

	// If the returned spot is nil, then the target is not present. So return false.
	if hashTable.employees[index] == nil {
		return false
	}

	// If the returned spot is deleted, then the target is not present. So return false.
	if hashTable.employees[index].deleted {
		return false
	}

	// Otherwise, return true.
	return true
}

// Make a display showing whether each array entry is nil.
func (hashTable *LinearProbingHashTable) dumpConcise() {
	// Loop through the array.
	for i, employee := range hashTable.employees {
		if employee == nil {
			// This spot is empty.
			fmt.Printf(".")
		} else if employee.deleted {
			// This spot is deleted.
			fmt.Printf("x")
		} else {
			// Display this entry.
			fmt.Printf("O")
		}
		if i%50 == 49 {
			fmt.Println()
		}
	}
	fmt.Println()
}

// Return the average probe sequence length for the items in the table.
func (hashTable *LinearProbingHashTable) aveProbeSequenceLength() float32 {
	totalLength := 0
	numValues := 0
	for _, employee := range hashTable.employees {
		if employee != nil {
			_, probeLength := hashTable.find(employee.name)
			totalLength += probeLength
			numValues++
		}
	}
	return float32(totalLength) / float32(numValues)
}

// Delete an item from the hash table.
func (hashTable *LinearProbingHashTable) delete(name string) {
	// Call find to get the index where the key belongs
	index, _ := hashTable.find(name)

	// If the returned index is at least 0 and that spot is not nil, then set its deleted value to true.
	if index >= 0 && hashTable.employees[index] != nil {
		hashTable.employees[index].deleted = true
	}
}

// Show this key's probe sequence.
func (hashTable *LinearProbingHashTable) probe(name string) int {
	// Hash the key.
	hash := hash(name) % hashTable.capacity
	fmt.Printf("Probing %s (%d)\n", name, hash)

	// Keep track of a deleted spot if we find one.
	deletedIndex := -1

	// Probe up to hashTable.capacity times.
	for i := 0; i < hashTable.capacity; i++ {
		index := (hash + i) % hashTable.capacity

		fmt.Printf("    %d: ", index)
		if hashTable.employees[index] == nil {
			fmt.Printf("---\n")
		} else if hashTable.employees[index].deleted {
			fmt.Printf("xxx\n")
		} else {
			fmt.Printf("%s\n", hashTable.employees[index].name)
		}

		// If this spot is empty, the value isn't in the table.
		if hashTable.employees[index] == nil {
			// If we found a deleted spot, return its index.
			if deletedIndex >= 0 {
				fmt.Printf("    Returning deleted index %d\n", deletedIndex)
				return deletedIndex
			}

			// Return this index, which holds nil.
			fmt.Printf("    Returning nil index %d\n", index)
			return index
		}

		// If this spot is deleted, remember where it is.
		if hashTable.employees[index].deleted {
			if deletedIndex < 0 {
				deletedIndex = index
			}
		} else if hashTable.employees[index].name == name {
			// If this cell holds the key, return its data.
			fmt.Printf("    Returning found index %d\n", index)
			return index
		}

		// Otherwise continue the loop.
	}

	// If we get here, then the key is not
	// in the table and the table is full.

	// If we found a deleted spot, return it.
	if deletedIndex >= 0 {
		fmt.Printf("    Returning deleted index %d\n", deletedIndex)
		return deletedIndex
	}

	// There's nowhere to put a new entry.
	fmt.Printf("    Table is full\n")
	return -1
}
