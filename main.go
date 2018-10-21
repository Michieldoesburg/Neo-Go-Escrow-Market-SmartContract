package MyFirstNeoGoContract

import (
	"neo-storm/token/nep5"
	"neo-storm/interop/runtime"
	"neo-storm/interop/engine"
	"neo-storm/interop/storage"
	"neo-storm/interop/util"
)

const (
	decimals   = 8
	multiplier = 100000000
)


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
	productHash := args[0].([]byte)

	if trigger == runtime.Application() {
		if operation == "create_offer" {
			// Key: Hash, Value: Item description, price, and contact info
			if checkArgs(args, 2) {
				return false
			}
			return createOffer(args)
		}
		if operation == "reject_application" {
			//TODO:
			return true
		}
		if operation == "confirm_purchase" {
			//TODO check if sender is the buyer
			buyer := retrieveBuyerAddress(productHash)
			validateCallingAddress(buyer)


			amountPre := []byte("AMOUNT")
			amountKey := append(amountPre, productHash...)
			productCost := storage.Get(storage.GetContext(), amountKey)

			seller := retrieveSellerAddress(productHash)

			token.Transfer(storage.GetContext(), engine.GetExecutingScriptHash(), buyer,  1/4*productCost.(int))
			token.Transfer(storage.GetContext(), engine.GetExecutingScriptHash(), seller,  3/4*productCost.(int))
			return true
		}

		if operation == "apply_to_offer" {
			// Tokens and buyer's contact info.
			if checkArgs(args, 2) {
				return false
			}
			applyToOffer(token, args)
		}
		if operation == "accept_offer" {
			// Gas and buyer's contact info.
			if checkArgs(args, 2) {
				return false
			}
			acceptOffer(token, args)
		}
	}
	

	return false
}


// The seller puts an item on offer on the chain
func createOffer(args []interface{}) bool {
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
func applyToOffer(t nep5.Token, args []interface{}) bool {
	
	productHash := args[0].([]byte)
	amountPre := []byte("AMOUNT")
	amountKey := append(amountPre, productHash...)
	productCost := storage.Get(storage.GetContext(), amountKey).(int)


	// Check if buyer has enough tokens.
	buyer := retrieveBuyerAddress(productHash)


	if t.Transfer(storage.GetContext(), buyer, engine.GetExecutingScriptHash(), 2*productCost) {
		buyerPre := []byte("BUYER")
		buyerKey := append(buyerPre, productHash...)
		storage.Put(storage.GetContext(), buyerKey, engine.GetCallingScriptHash())
		return true
	}
	return false
}

// The seller puts tokens in the escrow of the offer.
func acceptOffer(t nep5.Token, args []interface{}) bool {
	productHash := args[0].([]byte)

	//Validate if sender is the seller
	seller := retrieveSellerAddress(productHash)
	validateCallingAddress(seller)

	amountPre := []byte("AMOUNT")
	amountKey := append(amountPre, productHash...)
	productCost := storage.Get(storage.GetContext(), amountKey).(int)

	//Transfer tokens
	return t.Transfer(storage.GetContext(), seller, engine.GetExecutingScriptHash(), 2*productCost)

}

func checkArgs(args []interface{}, length int) bool {
	if len(args) == length {
		return true
	}

	return false
}

func validateCallingAddress(desiredAddress []byte) bool {
	return util.Equals(engine.GetCallingScriptHash, desiredAddress)
}

func retrieveSellerAddress(productHash []byte) []byte {
	sellerPre := []byte("SELLER")
	sellerKey := append(sellerPre, productHash...)
	return storage.Get(storage.GetContext(), sellerKey).([]byte)
	
}

func retrieveBuyerAddress(productHash []byte) []byte {
	buyerPre := []byte("BUYER")
	buyerKey := append(buyerPre, productHash...)
	return storage.Get(storage.GetContext(), buyerKey).([]byte)
}