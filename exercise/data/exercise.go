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
	for i, item := range p.Inventory {
		if item.Name == itemName {
			p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
			fmt.Printf("%s has Dropped %s/n", p.Name, itemName)
			return
		}
	}
	fmt.Printf("%s does not has %s in inventory./n", p.Name, itemName)
}

// Use an item (if potion, remove it after use)
func (p *Player) UseItem(itemName string) {
	for i, item := range p.Inventory {
		if item.Name == itemName {
			if item.Type == "potion" {
				fmt.Printf("%s has used %s and feels rejuninated!/n", p.Name, itemName)
				p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
			} else {
				fmt.Printf("%s uses %s/n", p.Name, itemName)
			}
			return
		}
	}
	fmt.Printf("%s does not has %s in inventory./n", p.Name, itemName)
}
