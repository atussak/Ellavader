package backup

import(
	"os"
	def "../definitions"
	"io/ioutil"
	"../elevio"
	"fmt"
	"encoding/json"
)

func WriteCabOrdersToFile(requests [][] bool) {

	data := make([]bool, def.NUM_FLOORS)


	for floor := 0; floor < def.NUM_FLOORS; floor++ {
		data[floor] = (requests[floor][elevio.BT_Cab])
	}
	
	file, err := os.Create("caborders.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	k, _ := json.Marshal(data)

	_, err = file.Write(k)

	if err != nil {
		panic(err)
	}

	file.Sync()
}

func ReadCabOrdersFromFile(receiver chan<- elevio.Order) {
	raw_data, err := ioutil.ReadFile("caborders.txt")


	data := make([]bool, def.NUM_FLOORS)
	err = json.Unmarshal(raw_data, &data)	

	if err != nil {
		panic(err)
	}

	for floor := 0; floor < def.NUM_FLOORS; floor++ {
		if (data[floor]) {
			fmt.Printf("Found order")
			receiver <- elevio.Order{floor, elevio.BT_Cab}
		}
	}
}