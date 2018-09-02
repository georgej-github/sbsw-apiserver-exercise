package main

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
	brothers    map[string][]*Node
	sisters     map[string][]*Node
)

const (
	// MALE - gender constants
	MALE = "male"
	// FEMALE - gender constants
	FEMALE     = "female"
	rootFamily = "0"
)

// Populate - retrieve data from CSV and build family tree data structure
func Populate() {
	dataFile, err := os.OpenFile("data/familytree.csv", os.O_RDONLY, os.ModePerm)
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
		//fmt.Println("Hello", family.ID, family.Name, family.Gender)
		familyTree[family.ID] = family
	}

	//set up tree
	for _, family := range familyTree {
		if family.ParentID != rootFamily {
			_, ok := familyTree[family.ParentID]
			if ok {
				familyTree[family.ParentID].children = append(familyTree[family.ParentID].children, family)
			}
		} else {
			root = family //upper-most family found
		}
	}

	//set up brothers and sisters
	brothers = make(map[string][]*Node)
	sisters = make(map[string][]*Node)
	for _, family := range familyTree {
		if len(family.children) > 0 {
			for _, child := range family.children {
				if child.Gender == "male" {
					brothers[family.ID] = append(brothers[family.ID], child)
				} else {
					sisters[family.ID] = append(sisters[family.ID], child)
				}
			}
		}
	}

	if root.ID == "" {
		panic("Upper-most / highest parent family not found")
	}
}

func unclesAndaunts(member string, parentGender string, siblingsGender string) []string {
	fmt.Printf("in uncleAndaunts |%s|\n", member)
	var list []string
	for _, family := range familyTree {
		//fmt.Printf("|%s|, |%s|\n", family.Name, family.Partner)
		if member == family.Name || member == family.Partner {
			fmt.Println("match found")
			fmt.Printf("name: %s, gender: %s, partner: %s, partnergender: %s, id: %s, parentid: %s\n",
				family.Name, family.Gender, family.Partner, family.PartnerGender, family.ID, family.ParentID)
			if family.ParentID != rootFamily {
				fmt.Printf("%s not from root family (%s, %s)\n", member, family.ParentID, family.Gender)
				if _, ok := familyTree[family.ParentID]; !ok {
					fmt.Printf("Could not find parent family of %s!\n", member)
					return nil
				}
				//eg if paternal (ie direct descendant parent is male)
				//(can skip if matching parent is not a direct descendant of family, as tree of 'in-law'/spouse not available)
				if familyTree[family.ParentID].Gender == parentGender {
					parent := familyTree[family.ParentID].Name
					upperParentID := familyTree[family.ParentID].ParentID
					if siblingsGender == MALE { //get 'all' brothers of father/mother
						// eg if paternal get father's brothers
						if _, ok := brothers[upperParentID]; ok {
							for _, brother := range brothers[upperParentID] {
								if parentGender == MALE && brother.Name != parent {
									list = append(list, brother.Name)
								}
							}
						}
						fmt.Printf("member's parent's brothers: %v\n", list)
						// eg if paternal get father's sister's husbands
						if _, ok := sisters[upperParentID]; ok {
							for _, sister := range sisters[upperParentID] {
								if sister.PartnerGender == MALE {
									if parentGender == MALE && sister.Partner != parent {
										list = append(list, sister.Partner)
									}
								}
							}
						}
						fmt.Printf("member's parent's brotherinLaws (sisters' husbands): %v\n", list)
					} else if siblingsGender == FEMALE { //get 'all' sisters of father/mother
						// eg if paternal get father's sisters
						if _, ok := sisters[upperParentID]; ok {
							for _, sister := range sisters[upperParentID] {
								if parentGender == FEMALE && sister.Name != parent {
									list = append(list, sister.Name)
								}
							}
						}
						fmt.Printf("member's parent's sisters: %v\n", list)
						// eg if paternal get father's brother's wives
						if _, ok := brothers[upperParentID]; ok {
							for _, brother := range brothers[upperParentID] {
								if brother.PartnerGender == FEMALE {
									if parentGender == FEMALE && brother.Partner != parent {
										list = append(list, brother.Partner)
									}
								}
							}
						}
						fmt.Printf("member's parent's sisterinLaws: %v\n", list)
					}
				}
			}
		}
	}
	fmt.Printf("final list: %v\n", list)
	return list
}

func inLaws(member string, siblingsGender string) []string {
	var list []string
	for _, family := range familyTree {
		if family.ParentID != rootFamily {
			if member == family.Name { //direct descendant
				//capture ID of member's family
				memberFamilyID := family.ID
				//only look at spouses of siblings (partner's family tree not available)
				if _, ok := familyTree[family.ParentID]; !ok {
					fmt.Printf("Could not find parent family of %s!\n", member)
					return nil
				}
				upperParentID := familyTree[family.ParentID].ParentID
				if siblingsGender == MALE {
					// Looking for brother-in-laws
					// Check sister list and look for male spouses
					if _, ok := sisters[upperParentID]; ok {
						for _, sister := range sisters[upperParentID] {
							if sister.PartnerGender == MALE && sister.ID != memberFamilyID {
								list = append(list, sister.Partner)
							}
						}
					}
					return list
				}
				if siblingsGender == FEMALE {
					// Looking for sister-in-laws
					// Check brother list and look for female spouses
					if _, ok := brothers[upperParentID]; ok {
						for _, brother := range brothers[upperParentID] {
							if brother.PartnerGender == FEMALE && brother.ID != memberFamilyID {
								list = append(list, brother.Partner)
							}
						}
					}
					return list
				}
			}
			if member == family.Partner { //spouse of direct descendant
				//only look at siblings of spouse
				//capture ID of member's spouse's family
				spouseFamilyID := family.ID
				spouseGender := family.Gender
				upperParentID := familyTree[family.ParentID].ParentID
				if siblingsGender == MALE {
					// get brother-in-laws
					if _, ok := brothers[upperParentID]; ok {
						for _, brother := range brothers[upperParentID] {
							if spouseGender == MALE && brother.ID != spouseFamilyID {
								list = append(list, brother.Name)
							}
						}
					}
				} else if siblingsGender == FEMALE {
					// get sister-in-laws
					if _, ok := sisters[upperParentID]; ok {
						for _, sister := range sisters[upperParentID] {
							if spouseGender == FEMALE && sister.ID != spouseFamilyID {
								list = append(list, sister.Name)
							}
						}
					}
				}
				return list
			}
		}
	}
	return list
}

func cousins(member string) []string {
	var list []string
	for _, family := range familyTree {
		if member == family.Name { //spouse's cousins cannot be traced in this tree
			fmt.Println("match found")
			fmt.Printf("name: %s, gender: %s, partner: %s, partnergender: %s, id: %s, parentid: %s\n",
				family.Name, family.Gender, family.Partner, family.PartnerGender, family.ID, family.ParentID)
			//capture ID of member's parents family
			memberParentFamilyID := family.ParentID
			//iterate through all children of members grandparents and get list of children belonging to the grandparents' children
			var grandParentFamily *Node
			if _, ok := familyTree[family.ParentID]; !ok {
				if family.ParentID == rootFamily {
					fmt.Println("member's parents are top-most family in tree, no cousins for members")
				}
				fmt.Println("could not find parent family of this member")
				return nil
			}
			if _, ok := familyTree[familyTree[family.ParentID].ParentID]; !ok {
				if familyTree[family.ParentID].ParentID != rootFamily {
					fmt.Println("could not find parent family of this member")
				}
				return nil
			}
			if familyTree[family.ParentID].ParentID == rootFamily {
				grandParentFamily = root
				fmt.Printf("grandparent family is root node, has %v children\n", len(grandParentFamily.children))
			} else {
				grandParentFamily = familyTree[familyTree[family.ParentID].ParentID]
				fmt.Printf("grandparent family found, ID %s, has %v children\n", grandParentFamily.ID, len(grandParentFamily.children))
			}
			for _, parentsSibling := range grandParentFamily.children {
				if len(parentsSibling.children) > 0 && (parentsSibling.ID != memberParentFamilyID) {
					for _, child := range parentsSibling.children {
						list = append(list, child.Name)
					}
				}
			}
			fmt.Println()
		}
	}
	return list
}

func search(member string, relationShip string) []string {
	fmt.Printf("search(%v), relation: %v\n", member, relationShip)
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
		return inLaws(member, MALE)
	case "sisterinLaw":
		return inLaws(member, FEMALE)
	case "cousins":
		return cousins(member)
	default:
		return nil
	}
}

/*
func main() {
	Populate()
	list := search("Jata", "paternalUncle")
	fmt.Printf("main: final list %v\n", list)
	list = search("Superman", "paternalUncle")
	fmt.Printf("main: final list %v\n", list)
	list = search("Lavnya", "paternalUncle")
	fmt.Printf("main: final list %v\n", list)
	list = search("Drita", "cousins")
	fmt.Printf("main: final list %v\n", list)
}
*/
