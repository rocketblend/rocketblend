package types

// CyclesDevice represents the available devices for rendering with the Cycles engine.
type CyclesDevice string

const (
	CyclesDeviceCPU    CyclesDevice = "CPU"
	CyclesDeviceCUDA   CyclesDevice = "CUDA"
	CyclesDeviceOPTIX  CyclesDevice = "OPTIX"
	CyclesDeviceHIP    CyclesDevice = "HIP"
	CyclesDeviceONEAPI CyclesDevice = "ONEAPI"
	CyclesDeviceMETAL  CyclesDevice = "METAL"
)
