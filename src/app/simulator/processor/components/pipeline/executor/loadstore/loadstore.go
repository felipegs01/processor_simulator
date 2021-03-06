package loadstore

import (
	"errors"
	"fmt"

	"app/logger"
	"app/simulator/processor/components/storagebus"
	"app/simulator/processor/consts"
	"app/simulator/processor/models/data"
	"app/simulator/processor/models/operation"
	"app/simulator/processor/models/set"
)

type LoadStore struct {
	bus *storagebus.StorageBus
}

func New(bus *storagebus.StorageBus) *LoadStore {
	return &LoadStore{bus: bus}
}

func (this *LoadStore) Bus() *storagebus.StorageBus {
	return this.bus
}

func (this *LoadStore) Process(operation *operation.Operation) (*operation.Operation, error) {

	instruction := operation.Instruction()
	rdAddress := instruction.Data.(*data.DataI).RegisterD.ToUint32()
	rsAddress := instruction.Data.(*data.DataI).RegisterS.ToUint32()
	immediate := instruction.Data.(*data.DataI).Immediate.ToUint32()

	switch instruction.Info.Opcode {
	case set.OP_LW:
		rsValue := this.Bus().LoadRegister(operation, rsAddress)
		value := this.Bus().LoadData(operation, rsValue+immediate)
		this.Bus().StoreRegister(operation, rdAddress, value)
		logger.Collect(" => [LS][%03d]: [R%d(%#02X) = MEM(%#02X) = %#08X]", operation.Id(), rdAddress, rdAddress*consts.BYTES_PER_WORD, rsValue+immediate, value)
	case set.OP_SW:
		rdValue := this.Bus().LoadRegister(operation, rdAddress)
		rsValue := this.Bus().LoadRegister(operation, rsAddress)
		this.Bus().StoreData(operation, rdValue+immediate, rsValue)
		logger.Collect(" => [LS][%03d]: [MEM(%#02X) = %#08X]", operation.Id(), rdValue+immediate, rsValue)
	case set.OP_LLI:
		this.Bus().StoreRegister(operation, rdAddress, immediate)
		logger.Collect(" => [LS][%03d]: [R%d(%#02X) = %#08X]", operation.Id(), rdAddress, rdAddress*consts.BYTES_PER_WORD, immediate)
	case set.OP_SLI:
		rdValue := this.Bus().LoadRegister(operation, rdAddress)
		this.Bus().StoreData(operation, rdValue, immediate)
		logger.Collect(" => [LS][%03d]: [MEM(%#02X) = %#08X]", operation.Id(), rdValue, immediate)
	case set.OP_LUI:
		this.Bus().StoreRegister(operation, rdAddress, immediate<<16)
		logger.Collect(" => [LS][%03d]: [R%d(%#02X) = %#08X]", operation.Id(), rdAddress, rdAddress*consts.BYTES_PER_WORD, immediate<<16)
	case set.OP_SUI:
		rdValue := this.Bus().LoadRegister(operation, rdAddress)
		this.Bus().StoreData(operation, rdValue, immediate<<16)
		logger.Collect(" => [LS][%03d]: [MEM(%#02X) = %#08X]", operation.Id(), rdValue, immediate<<16)
	default:
		return operation, errors.New(fmt.Sprintf("Invalid operation to process by Data unit. Opcode: %d", instruction.Info.Opcode))
	}
	return operation, nil
}
