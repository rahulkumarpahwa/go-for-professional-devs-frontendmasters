package data

import "fmt"

type Item struct {
	Name string
	Type string
}

type Player struct {
	Name      string
	Inventory []Item
}

func (p *Player) PickUpItem(item Item) {
	p.Inventory = append(p.Inventory, item)
	fmt.Printf("%s has PickedUp %s/n", p.Name, item.Name)

}

func (p *Player) DropItem(itemName string) {
	var newInventory []Item
	for _, i := range p.Inventory {
		if i.Name != itemName {
			newInventory = append(newInventory, i)
		}
	}
	fmt.Printf("%s has Dropped %s/n", p.Name, itemName)
	p.Inventory = newInventory
}

// Use an item (if potion, remove it after use)
func (p *Player) UseItem(itemName string) {
	// TODO: Implement function to use an item
	var newInventory []Item
	for _, i := range p.Inventory {
		if i.Name != itemName {
			newInventory = append(newInventory, i)
		}
	}
	fmt.Printf("%s has uses %s/n", p.Name, itemName)
	p.Inventory = newInventory
}
