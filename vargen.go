package main

import "github.com/jtcaraball/goexpdt/query"

const sep string = string(rune(30))

// varGenBotCount returns a variable with value equal to v with the addition of
// the prefix "botc" separated with the record separator character (ascii 30).
func varGenBotCount(v query.QVar) query.QVar {
	return query.QVar("botc" + sep + string(v))
}

// varGenBotCount returns a variable with value equal to v with the addition of
// the prefix "reach" separated with the record separator character (ascii 30).
func varGenNodeReach(v query.QVar) query.QVar {
	return query.QVar("reach" + sep + string(v))
}

// varGenHammingDistance returns a variable equal to the sorted concatenation
// of variables v1 and v2 with the addition of the prefix "hdist" separated
// using the record separator character (ascii 30).
func varGenHammingDistance(v1, v2 query.QVar) query.QVar {
	if string(v1) < string(v2) {
		return query.QVar("hdist" + sep + string(v1) + sep + string(v2))
	}
	return query.QVar("hdist" + sep + string(v2) + sep + string(v1))
}
