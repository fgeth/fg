package common


import(
		"math/big"
)


type MoneySupply struct {
	Dollars		float64						//Current amount of money valued in virtual Dollars that is in circulation
	FG			*big.Int					//Current amount of money valued in Coins that is in circulation
	Max			float64						//Largest amount of money valued in virtual Dollars that has been in circulation (should be more than current value as notes and stacks get burned and delay in replacing them)
	LMax		float64						//Last years Max amount of money in circulation valued in virtual Dollars.
	IR			float64						//Current Inflation Rate (Max - LMax) / LMax
}