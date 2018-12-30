package main

import (
    "strings"
    "testing"
)

func GetEmptyServerFilterList() []string {
    return make([]string, 0)
}

func GetSingleServerFilterListWith(serverToAdd string) []string {
    filterList := GetEmptyServerFilterList()
    filterList = append(filterList, serverToAdd)
    return filterList
}

func TestWholeTableGeneration(t *testing.T) {
    table := BuildTable("test_inventories/simple_valid_inventory", GetEmptyServerFilterList())

    for _, testCase := range []string {
        "one.example.com",
        "DBSERVERS",
        "WEBSERVERS" } {

        if ! strings.Contains(table, testCase) {
            t.Errorf("Expected %s in output but was not found. Output was: %s\n", testCase, table)
        }
    }
}

func TestFilter(t *testing.T) {
    table := BuildTable("test_inventories/simple_valid_inventory", GetSingleServerFilterListWith("one.example.com"))

    for _, testCase := range []string {
        "two.example.com",
        "foo.example.com" } {

        if strings.Contains(table, testCase) {
            t.Errorf("Expected %s not to be in output but was found. Output was: %s\n", testCase, table)
        }
    }
}
