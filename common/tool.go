package common

import (
	"encoding/binary"
	"github.com/kpli0rn/jdwpgo/protocol/vm"
)

func GetMethodByName(methods *vm.AllMethodsReply, name string) *vm.AllMethodsMethod {
	for _, method := range methods.Methods {
		if method.Name.String() == name {
			return &method
		}
	}
	return nil
}

func ParseEvent(buf []byte, eventId int32, idsize *vm.IDSizesReply) (int32, uint64) {
	raw := buf[11:]
	rId := int32(binary.BigEndian.Uint32(raw[6:10]))
	if rId != eventId {
		return 0, 0
	}
	rawtId := raw[10 : 10+idsize.ObjectIDSize]
	tId := binary.BigEndian.Uint64(rawtId)
	return rId, tId
}
