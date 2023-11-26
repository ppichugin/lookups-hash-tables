package main

import "fmt"

var employeeNames [100]string

type Employee struct {
	name  string
	phone string
}

type ChainingHashTable struct {
	numBuckets int
	buckets    [][]*Employee
}

func main() {
	// Make some names.
	employees := []Employee{
		{"Ann Archer", "202-555-0101"},
		{"Bob Baker", "202-555-0102"},
		{"Cindy Cant", "202-555-0103"},
		{"Dan Deever", "202-555-0104"},
		{"Edwina Eager", "202-555-0105"},
		{"Fred Franklin", "202-555-0106"},
		{"Gina Gable", "202-555-0107"},
		{"Herb Henshaw", "202-555-0108"},
		{"Ida Iverson", "202-555-0109"},
		{"Jeb Jacobs", "202-555-0110"},
	}

	hashTable := NewChainingHashTable(10)
	for _, employee := range employees {
		hashTable.set(employee.name, employee.phone)
	}
	hashTable.dump()

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

// Initialize a ChainingHashTable and return a pointer to it.
func NewChainingHashTable(numBuckets int) *ChainingHashTable {
	// Create a new ChainingHashTable
	table := &ChainingHashTable{
		numBuckets: numBuckets,
		// Allocate the slice of buckets
		buckets: make([][]*Employee, numBuckets),
	}

	// Return the pointer to the new ChainingHashTable
	return table
}

// Display the hash table's contents.
func (hashTable *ChainingHashTable) dump() {
	for i, bucket := range hashTable.buckets {
		fmt.Printf("Bucket %d:\n", i)
		for _, employee := range bucket {
			fmt.Printf("\t%s: %s\n", employee.name, employee.phone)
		}
	}
}

// Find the bucket and Employee holding this key.
// Return the bucket number and Employee number in the bucket.
// If the key is not present, return the bucket number and -1.
func (hashTable *ChainingHashTable) find(name string) (int, int) {
	// Calculate the hash of the name
	hashValue := hash(name)

	// Use the hash to find the bucket index
	bucketIndex := hashValue % hashTable.numBuckets

	// Search for the employee in the bucket
	for i, employee := range hashTable.buckets[bucketIndex] {
		if employee.name == name {
			return bucketIndex, i
		}
	}

	// If the employee is not found, return -1, -1
	return -1, -1
}

// Add an item to the hash table.
func (hashTable *ChainingHashTable) set(name string, phone string) {
	// Calculate the hash of the name
	hashValue := hash(name)

	// Use the hash to find the bucket index
	bucketIndex := hashValue % hashTable.numBuckets

	// Search for the employee in the bucket
	for i, employee := range hashTable.buckets[bucketIndex] {
		if employee.name == name {
			// If the employee is found, update the phone number
			hashTable.buckets[bucketIndex][i].phone = phone
			return
		}
	}

	// If the employee is not found, add them to the bucket
	hashTable.buckets[bucketIndex] = append(
		hashTable.buckets[bucketIndex],
		&Employee{name: name, phone: phone},
	)
}

// Return an item from the hash table.
func (hashTable *ChainingHashTable) get(name string) string {
	// Call find to get the indices of the bucket and Employee struct
	bucketIndex, employeeIndex := hashTable.find(name)

	// If the Employee index is at least 0, return its phone value
	if employeeIndex >= 0 {
		return hashTable.buckets[bucketIndex][employeeIndex].phone
	}

	// If the Employee index is less than 0, return an empty string
	return ""
}

// Return true if the person is in the hash table.
func (hashTable *ChainingHashTable) contains(name string) bool {
	// Call find to get the indices of the bucket and Employee struct
	_, employeeIndex := hashTable.find(name)

	// If the Employee index is -1, return false. Otherwise, return true.
	return employeeIndex != -1
}

// Delete this key's entry.
func (hashTable *ChainingHashTable) delete(name string) {
	// Call find to get the indices of the bucket and Employee struct
	bucketIndex, employeeIndex := hashTable.find(name)

	// If the Employee index is at least 0, cut that struct out of its bucket
	if employeeIndex >= 0 {
		hashTable.buckets[bucketIndex] = append(
			hashTable.buckets[bucketIndex][:employeeIndex],
			hashTable.buckets[bucketIndex][employeeIndex+1:]...,
		)
	}

	// If the Employee index is less than 0, do nothing
}
