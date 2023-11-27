package opts

import (
	"encoding/json"
	"log"
)

type Params struct {
	Iface   string `           long:"interface"      env:"INTERFACE"  description:"Server network interface" default:"eno1"`
	MapFile string `           long:"file"     env:"MAP_FILE"  description:"File with name to mac mappings in json" default:"mapping.yaml"`
}

func (o *Params) GetJson() []byte {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		log.Panic(err)
	}
	return jsonBytes
}
