package MyFirstNeoGoContract

import (
	"github.com/CityOfZion/neo-storm/examples/token/nep5"
	"github.com/CityOfZion/neo-storm/interop/runtime"
	"github.com/CityOfZion/neo-storm/interop/storage"
	"github.com/CityOfZion/neo-storm/interop/util"
	"github.com/CityOfZion/neo-storm/interop/runtime"	
)

// Check if the invoker of the contract is the specified owner
var owner = util.FromAddress("Aej1fe4mUgou48Zzup5j8sPrE3973cJ5oz")

var locked := false

func Main(operation string, args []interface{}) bool {
	trigger := runtime.GetTrigger()

	// Log owner upon Verification trigger
	if trigger == runtime.Verification() {
		if runtime.CheckWitness(owner) {
			if !locked {
				runtime.Log("Verified Owner")
				return true
			}
		}
	}

	// Discerns between log and notify for this test
	if trigger == runtime.Application() {
		return handleOperation(operation, args)
	}

	return false
}

func handleOperation(operation string, args []interface{}) bool {
	if operation == "create_offer" {
		// Key: Hash, Value: Item description, price, and contact info
		if checkArgs(args, 2) {
			return false
		}
		create_offer(args)
	}

	if operation == "apply_on_offer" {
		// Gas and buyer's contact info.
		if checkArgs(args, 2) {
			return false
		}
		apply_on_offer()
	}


	if operation == "accept_application" {
		// Gas
		if checkArgs(args, 2) {
			return false
		}
		accept_application()
	}

	if operation == "reject_application" {
		// Gas
		if checkArgs(args, 2) {
			return false
		}
		reject_application()
	}


	return false
}

// The seller puts an item on offer on the chain
func create_offer() bool {
	ctx := storage.GetContext()

	key := args[0].([]byte)
	value := args[1].([]byte)
	storage.Put(ctx, key, value)
	return true
}

// The buyer applies to seller's offer.
func apply_on_offer() bool {
	tx := GetScriptContainer()

	// Check if transaction has currency attached
	references := tx.References
	if len(references) < 1:
		return false
	reference := references[0]


	// Seller's address
	sender := GetScriptHash(reference)

	// Contract's address	
	contract_script_hash = GetExecutingScriptHash()

	context = GetContext()

	// Handle outputs in the transaction
	for index, output := range tx.Outputs {
		shash := GetScriptHash(output)
		output_asset_id = GetAssetId(output)

		if shash == contract_script_hash {
			output_val := GetValue(output)


		}
	} 
}

// The seller accepts the buyer's application
func accept_application() bool {
	locked = true
}

func reject_application() bool {

}

func getCurrentTimeStamp() {
	height := GetHeight()
	currentBlock := GetHeader(height)
	time := currentBlock.Timestamp 
	return time
}

func checkArgs(args []interface{}, length int) bool {
	if len(args) == length {
		return true
	}

	return false
}

func checkState(state int) bool {

}