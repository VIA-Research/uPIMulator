package misc

import (
	"errors"
)

type ConfigValidator struct {
	config_loader *ConfigLoader
}

func (this *ConfigValidator) Init(config_loader *ConfigLoader) {
	this.config_loader = config_loader
}

func (this *ConfigValidator) Validate() {
	if this.config_loader.AddressWidth() <= 0 {
		err := errors.New("address width <= 0")
		panic(err)
	}

	if this.config_loader.AtomicDataWidth() <= 0 {
		err := errors.New("atomic data width <= 0")
		panic(err)
	}

	if this.config_loader.IramDataWidth() <= 0 {
		err := errors.New("IRAM data width <= 0")
		panic(err)
	}

	if this.config_loader.WramDataWidth() <= 0 {
		err := errors.New("WRAM data width <= 0")
		panic(err)
	}

	if this.config_loader.MramDataWidth() <= 0 {
		err := errors.New("MRAM data width <= 0")
		panic(err)
	}

	if this.config_loader.AtomicOffset() < 0 {
		err := errors.New("atomic offset < 0")
		panic(err)
	}

	if this.config_loader.IramOffset() < 0 {
		err := errors.New("IRAM offset < 0")
		panic(err)
	}

	if this.config_loader.WramOffset() < 0 {
		err := errors.New("WRAM offset < 0")
		panic(err)
	}

	if this.config_loader.MramOffset() < 0 {
		err := errors.New("MRAM offset < 0")
		panic(err)
	}

	if this.config_loader.AtomicSize() <= 0 {
		err := errors.New("atomic size <= 0")
		panic(err)
	}

	if this.config_loader.IramSize() <= 0 {
		err := errors.New("IRAM size <= 0")
		panic(err)
	}

	if this.config_loader.WramSize() <= 0 {
		err := errors.New("WRAM size <= 0")
		panic(err)
	}

	if this.config_loader.MramSize() <= 0 {
		err := errors.New("MRAM size <= 0")
		panic(err)
	}

	if this.AreOverlapped(
		this.config_loader.AtomicOffset(),
		this.config_loader.AtomicSize(),
		this.config_loader.IramOffset(),
		this.config_loader.IramSize(),
	) {
		err := errors.New("atomic and IRAM are overlapped")
		panic(err)
	}

	if this.AreOverlapped(
		this.config_loader.AtomicOffset(),
		this.config_loader.AtomicSize(),
		this.config_loader.WramOffset(),
		this.config_loader.WramSize(),
	) {
		err := errors.New("atomic and WRAM are overlapped")
		panic(err)
	}

	if this.AreOverlapped(
		this.config_loader.AtomicOffset(),
		this.config_loader.AtomicSize(),
		this.config_loader.MramOffset(),
		this.config_loader.MramSize(),
	) {
		err := errors.New("atomic and MRAM are overlapped")
		panic(err)
	}

	if this.AreOverlapped(
		this.config_loader.IramOffset(),
		this.config_loader.IramSize(),
		this.config_loader.WramOffset(),
		this.config_loader.WramSize(),
	) {
		err := errors.New("IRAM and WRAM are overlapped")
		panic(err)
	}

	if this.AreOverlapped(
		this.config_loader.IramOffset(),
		this.config_loader.IramSize(),
		this.config_loader.MramOffset(),
		this.config_loader.MramSize(),
	) {
		err := errors.New("IRAM and MRAM are overlapped")
		panic(err)
	}

	if this.AreOverlapped(
		this.config_loader.WramOffset(),
		this.config_loader.WramSize(),
		this.config_loader.MramOffset(),
		this.config_loader.MramSize(),
	) {
		err := errors.New("WRAM and MRAM are overlapped")
		panic(err)
	}

	if this.config_loader.StackSize() <= 0 {
		err := errors.New("stack size <= 0")
		panic(err)
	}

	if this.config_loader.HeapSize() <= 0 {
		err := errors.New("heap size <= 0")
		panic(err)
	}

	if this.config_loader.NumGpRegisters() <= 0 {
		err := errors.New("num gp registers <= 0")
		panic(err)
	}

	if this.config_loader.MaxNumTasklets() <= 0 {
		err := errors.New("max num tasklets <= 0")
		panic(err)
	}
}

func (this *ConfigValidator) AreOverlapped(
	offset1 int64,
	size1 int64,
	offset2 int64,
	size2 int64,
) bool {
	if offset1 <= offset2 && offset2 <= offset1+size1 {
		return true
	} else if offset1 <= offset2+size2 && offset2+size2 <= offset1+size1 {
		return true
	} else if offset2 <= offset1 && offset1 <= offset2+size2 {
		return true
	} else if offset2 <= offset1+size1 && offset1+size1 <= offset2+size2 {
		return true
	} else {
		return false
	}
}
