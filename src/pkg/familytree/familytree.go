package main

//partially adapted from https://blog.serverbooter.com/post/parsing-nested-json-in-go/

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

// Node - family node struct
type Node struct {
	ID            string `csv:"FamilyID"`
	Name          string `csv:"Name"`
	Gender        string `csv:"Gender"`
	Partner       string `csv:"Partner"`
	PartnerGender string `csv:"PartnerGender"`
	ParentID      string `csv:"ParentID"`
	children      []*Node
}

var (
	familyTree  map[string]*Node
	familyTable []*Node
	root        *Node
	brothers    map[string][]string
	sisters     map[string][]string
)

const MALE = "male"
const FEMALE = "female"

// Populate - retrieve data from CSV and build family tree data structure
func Populate() {
	dataFile, err := os.OpenFile("familytree.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	familyTable = []*Node{}
	if err := gocsv.UnmarshalFile(dataFile, &familyTable); err != nil { // Load clients from file
		panic(err)
	}

	familyTree = make(map[string]*Node)
	for _, family := range familyTable {
		fmt.Println("Hello", family.ID, family.Name, family.Gender)
		familyTree[family.ID] = family
	}

	//set up tree
	for _, family := range familyTree {
		if family.ParentID != "0" {
			_, ok := familyTree[family.ParentID]
			if ok {
				familyTree[family.ParentID].children = append(familyTree[family.ParentID].children, family)
			}
		} else {
			root = family //upper-most family found
		}
	}

	//set up brothers and sisters
	brothers = make(map[string][]string)
	sisters = make(map[string][]string)
	for _, family := range familyTree {
		if len(family.children) > 0 {
			for _, child := range family.children {
				if child.Gender == "male" {
					brothers[family.ID] = append(brothers[family.ID], child.Name)
				} else {
					sisters[family.ID] = append(sisters[family.ID], child.Name)
				}
			}
		}
	}

	if root.ID == "" {
		panic("Upper-most / highest parent family not found")
	}
}

func unclesAndaunts(member string, expectedParentGender string, siblingsGender string) []string {
	for _, family := range familyTree {
		if member == family.Name || member == family.Partner {
			if family.ParentID != "0" {
				if familyTree[family.ParentID].ParentID != "0" {
					if familyTree[family.ParentID].Gender == expectedParentGender {
						upperParentID := familyTree[family.ParentID].ParentID
						if siblingsGender == MALE {
							if _, ok := brothers[upperParentID]; ok {
								return brothers[upperParentID]
							}
						} else {
							if _, ok := sisters[upperParentID]; ok {
								return sisters[upperParentID]
							}
						}
					}
				}
			}
		}
		return nil
	}
	return nil
}

func brothersAndsisters(member string, siblingsGender string) []string {
	for _, family := range familyTree {
		if member == family.Name || member == family.Partner {
			if family.ParentID != "0" {
				var sibLings []string
				for _, child := range familyTree[family.ParentID].children {
					if child.Name != member && child.Partner != member {
						if child.Partner != "" {
							if siblingsGender == child.PartnerGender {
								sibLings = append(sibLings, child.Partner)
							}
						}
						if siblingsGender == child.Gender {
							sibLings = append(sibLings, child.Gender)
						}
					}
				}
				return sibLings
			}
		}
	}
	return nil
}

func search(member string, relationShip string) []string {
	switch relationShip {
	case "paternalUncle":
		return unclesAndaunts(member, MALE, MALE)
	case "paternalAunt":
		return unclesAndaunts(member, MALE, FEMALE)
	case "maternalAunt":
		return unclesAndaunts(member, FEMALE, FEMALE)
	case "maternalUncle":
		return unclesAndaunts(member, FEMALE, MALE)
	case "brotherInLaw":
		return brothersAndsisters(member, MALE)
	case "sisterinLaw":
		return brothersAndsisters(member, FEMALE)
	default:
		return nil
	}
}

func main() {
	Populate()
}
