package misc

type ConfigLoader struct {
}

func (this *ConfigLoader) Init() {
}

func (this *ConfigLoader) AddressWidth() int {
	return 32
}

func (this *ConfigLoader) AtomicDataWidth() int {
	return 32
}

func (this *ConfigLoader) AtomicOffset() int64 {
	return 0
}

func (this *ConfigLoader) AtomicSize() int64 {
	return 256
}

func (this *ConfigLoader) IramDataWidth() int {
	return 96
}

func (this *ConfigLoader) IramOffset() int64 {
	return 384 * 1024
}

func (this *ConfigLoader) IramSize() int64 {
	return 48 * 1024
}

func (this *ConfigLoader) WramDataWidth() int {
	return 32
}

func (this *ConfigLoader) WramOffset() int64 {
	return 512
}

func (this *ConfigLoader) WramSize() int64 {
	return 128 * 1024
}

func (this *ConfigLoader) MramDataWidth() int {
	return 32
}

func (this *ConfigLoader) MramOffset() int64 {
	return 512 * 1024
}

func (this *ConfigLoader) MramSize() int64 {
	return 64 * 1024 * 1024
}

func (this *ConfigLoader) StackSize() int64 {
	return 2 * 1024
}

func (this *ConfigLoader) HeapSize() int64 {
	return 4 * 1024
}

func (this *ConfigLoader) NumGpRegisters() int {
	return 24
}

func (this *ConfigLoader) MaxNumTasklets() int {
	return 24
}

func (this *ConfigLoader) VmBankOffset() int64 {
	return 512
}

func (this *ConfigLoader) VmBankSize() int64 {
	return 128 * 1024 * 1024
}

func (this *ConfigLoader) VmBg0() int {
	return 6
}

func (this *ConfigLoader) VmBg1() int {
	return 17
}

func (this *ConfigLoader) VmBank() int {
	return 18
}

func (this *ConfigLoader) VmMemorySize() int64 { return 1024 }

func (this *ConfigLoader) GarbageCollectionThreshold() int64 {
	return 100
}
