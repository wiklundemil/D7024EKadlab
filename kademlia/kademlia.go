package kademlia

type Kademlia struct {
	RoutingTable *RoutingTable
	Network *Network

	// Initialize data storage as a map with string keys and byte slice values
	// This map will be used to store data in the Kademlia instance.
	// - Keys (string): These represent unique identifiers for the data entries.
	// - Values ([]byte): These are byte slices, allowing for flexible storage of various types of data, such as text, files, or other serialized objects.
	// Using a map provides efficient lookups, insertions, and deletions, which is crucial for the performance of the distributed system.
	Data *map[string][]byte 
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
