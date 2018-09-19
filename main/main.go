//Copyright Copyright 2018 Andrey Z.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"fmt"
	"github.com/AndreyZWorkAccount/FuzzyTextSearch/extensions"
	search "github.com/AndreyZWorkAccount/FuzzyTextSearch/fuzzySearch"
	"github.com/AndreyZWorkAccount/FuzzyTextSearch/levenshteinAlg"
	"time"
)

func main() {
	//setup
	const testWord = "miztake"
	const dictionaryFileName = "main\\bigdict.txt"
	const dictionarySize = 256
	const topCount = 20
	requestProcessingTime := time.Second * 50
	costs := levenshteinAlg.ChangesCosts{AddCost: 1, RemoveCost: 1, ReplaceCost: 1}

	fmt.Printf("Word to search: %v.\n\n", testWord)

	//read input
	ok, dictionaries := readDictionaries(dictionaryFileName, dictionarySize)
	if !ok {
		return
	}

	//run processor
	processor := search.NewProcessor(dictionaries, requestProcessingTime, costs)
	processor.Start()

	//send request
	processor.Requests() <- search.NewRequest(testWord)
	result := waitForResponse(requestProcessingTime, processor.Responses())

	if result == nil {
		return
	}

	fmt.Println("Most matching:")
	for _, res := range result.GetItems(topCount) {
		fmt.Printf("%v  ( distance:  %v ).\n", res.Word, res.Distance)
	}
}

func waitForResponse(requestProcessingTime time.Duration, responses <-chan search.Response) search.Response {
	defer extensions.TrackTime(time.Now(), "waitForResponse")

	requestBreak := time.After(requestProcessingTime)
	for {
		select {
		case response := <-responses:
			return response
		case <-requestBreak:
			fmt.Println("Processing timeout.")
			return nil
		default:
		}
	}
}
