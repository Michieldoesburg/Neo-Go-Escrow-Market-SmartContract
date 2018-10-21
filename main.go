package MyFirstNeoGoContract

import (
	"github.com/CityOfZion/neo-storm/examples/token/nep5"
	"github.com/CityOfZion/neo-storm/interop/runtime"
	"github.com/CityOfZion/neo-storm/interop/storage"
	"github.com/CityOfZion/neo-storm/interop/util"
	"github.com/CityOfZion/neo-storm/interop/runtime"	
)

const (
	decimals   = 8
	multiplier = 100000000
)

const VERIFICATION_R = 0x01

// CreateToken initializes the Token Interface for the Smart Contract to operate with
func CreateToken() nep5.Token {
	return nep5.Token{
		Name:           "Marketplace Token",
		Symbol:         "MKT",
		Decimals:       decimals,
		Owner:          engine.GetExecutingScriptHash(),
		TotalSupply:    10000000 * multiplier,
		CirculationKey: "MarketplaceTokenCirculation",
	}
}

func Main(operation string, args []interface{}) bool {
	token := CreateToken()
	trigger := runtime.GetTrigger()

	if trigger == runtime.Application() {
		if operation == "create_offer" {
			// Key: Hash, Value: Item description, price, and contact info
			if checkArgs(args, 2) {
				return false
			}
			return create_offer(args)
		}
		if operation == "reject_application" {
			//TODO:
			return true
		}
		if operation == "confirm_purchase" {
			//TODO check if sender is the buyer
			buyerPre := []byte("BUYER")
			buyerKey := append(buyerPre, productHash...)
			buyer := storage.Get(storage.GetContext(), buyerKey)
			validateCallingAddress(buyer)

			productHash := args[0].([]byte)

			amountPre := []byte("AMOUNT")
			amountKey := append(amountPre, productHash...)
			productCost := storage.Get(storage.GetContext(), amountKey)

			sellerPre := []byte("SELLER")
			sellerKey := append(sellerPre, productHash...)
			seller := storage.Get(storage.GetContext(), sellerKey)

			token.Transfer(storage.GetContext(), engine.GetExecutingScriptHash(), buyer.([]byte),  1/4*productCost.(int))
			token.Transfer(storage.GetContext(), engine.GetExecutingScriptHash(), seller.([]byte),  3/4*productCost.(int))
			return true
		}

		if operation == "apply_to_offer" {
			// Tokens and buyer's contact info.
			if checkArgs(args, 2) {
				return false
			}
			apply_to_offer(token, args)
		}
		if operation == "accept_offer" {
			// Gas and buyer's contact info.
			if checkArgs(args, 2) {
				return false
			}
			accept_offer(token, args)
		}
	}
	

	return false
}


// The seller puts an item on offer on the chain
func create_offer(args []interface{}) bool {
	ctx := storage.GetContext()

	productHash := args[0].([]byte)

	descriptionPre := []byte("DESCRIPTION")
	descriptionKey := append(descriptionPre, productHash...)
	storage.Put(ctx, descriptionKey, args[0].([]byte))

	amountPre := []byte("AMOUNT")
	amountKey := append(amountPre, productHash...)
	storage.Put(ctx, amountKey, args[1].([]byte))

	sellerPre := []byte("SELLER")
	sellerKey := append(sellerPre, productHash...)
	storage.Put(ctx, sellerKey, engine.GetCallingScriptHash())

	return true
}

// The buyer applies to seller's offer.
func apply_to_offer(t Token, args []interface{}) bool {
	amountSent := args[0]
	
	productHash := args[0].([]byte)
	amountPre := []byte("AMOUNT")
	amountKey := append(amountPre, productHash...)
	productCost := storage.Get(storage.GetContext(), amountKey)


	// Check if buyer has enough tokens.
	buyerPre := []byte("BUYER")
	buyerKey := append(buyerPre, productHash...)
	buyer := storage.Get(storage.GetContext(), buyerKey)

	// Not enough tokens
	if !token.canTransfer(storage.getContext, buyer([]byte), engine.GetCallingScriptHash(), amountSent) {
		return false
	}

	// Tokens don't match price + escrow
	if amountSent != 2 * productCost.(int) { 
		return false
	}

	token.tranfser(storage.getContext, buyer([]byte), engine.GetExecutingScriptHash(), amountSent)

	buyerPre := []byte("BUYER")
	buyerKey := append(buyerPre, productHash...)
	storage.Put(storage.GetContext(), buyerKey, engine.GetCallingScriptHash())

	//TODO put buyers address on the chain (args[1])

	return true

}

// The seller puts tokens in the escrow of the offer.
func accept_offer(t Token, args []interface{}) bool {
	amountSent := args[0]

	//Validate if sender is the seller
	sellerPre := []byte("SELLER")
	sellerKey := append(sellerPre, productHash...)
	seller := storage.Get(storage.GetContext(), sellerKey)
	validateCallingAddress(seller)

	// Not enough tokens
	if !token.canTransfer(storage.getContext, seller([]byte), engine.GetExecutingScriptHash(), amountSent) {
		return false
	}

	productHash := args[0].([]byte)
	amountPre := []byte("AMOUNT")
	amountKey := append(amountPre, productHash...)
	productCost := storage.Get(storage.GetContext(), amountKey)

	// Tokens don't match price + escrow
	if amountSent != 2 * productCost.(int) { 
		return false
	}

	//Transfer tokens
	token.tranfser(storage.getContext, seller([]byte), engine.GetExecutingScriptHash(), amountSent)

	return true

}

func checkArgs(args []interface{}, length int) bool {
	if len(args) == length {
		return true
	}

	return false
}

func validateCallingAddress(desiredAddress []byte) bool {
	return engine.GetCallingScriptHash == desiredAddress
}

func retrieveSellerAddress() interface{} {
	sellerPre := []byte("SELLER")
	sellerKey := append(sellerPre, productHash...)
	seller := storage.Get(storage.GetContext(), sellerKey)
}

func retrieveBuyerAddress() interface{} {
	buyerPre := []byte("BUYER")
	buyerKey := append(buyerPre, productHash...)
	buyer := storage.Get(storage.GetContext(), buyerKey)
}