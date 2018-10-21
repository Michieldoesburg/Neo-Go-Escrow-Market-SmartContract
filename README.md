# NEO NEP-5 Token Smart Contract in Go
KPN - TU Delft - NEO Hackathon, Rotterdam/Delft, 20/21 October 2018

https://www.eventbrite.co.uk/e/neo-blockchain-hackathon-rotterdam-delft-tickets-50526211258

Neo Smart contract to implement a P2P market of real life items with built-in escrow system. 

Sellers can offer items for sale by creating an offer on the blockchain. 

The escrow system works by having both buyer and seller commit 2x the price of the item being sold. The funds are released when the buyer confirms the purchase was succesful: 3 parts will go to the seller, and 1 part returns to the buyer.

The smart contract implementing the market is in the file 'main.go'.

The contract offers the following invocation options:

create_offer: The seller creates an offer on the blockchain of an item with a corresponding price in nep5 tokens. 

apply_to_offer: A buyer can 'apply' to buy the item by sending and commits 2x the price of the item to the contract. 

accept_offer: A seller can accept the buy, and also commits 2x the price of the item to the contract. At this point the funds are locked and irretrievable until the buyer confirms the purchase. The seller sends the item in the real world (food delivery, by mail, etc...)

confirm_purchase: A buyer confirms that the purchase was succesful when the item arrives. The contract releases the funds by sending 3 parts of the stored tokens to the seller (2 parts seller's escrow, 1 part payment for the item), and 1 part of the stored tokens back to the buyer (his escrow).  

reject_application: Before the seller can accept an offer and enter into the escrow, he can still choose to cancel the transaction. The contract refunds the tokens the buyer committed. (Not implemented).

Outside the locked state, the contract will automatically time out after a certain time and refund the tokens the buyer committed. This is not possible after the seller has also entered into the escrow and the contract is in the locked state, at which point only the buyer can choose to release the funds. (Not implemented). 


