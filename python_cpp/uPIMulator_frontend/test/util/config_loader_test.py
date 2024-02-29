from util.config_loader import ConfigLoader


def is_overlap(offset1: int, size1: int, offset2: int, size2) -> bool:
    if offset2 <= offset1 <= offset2 + size2:
        return True
    elif offset1 <= offset2 <= offset1 + size1:
        return True
    elif offset2 <= offset1 + size1 <= offset2 + size2:
        return True
    elif offset1 <= offset2 + size2 <= offset1 + size1:
        return True
    else:
        return False


def test_overlap():
    atomic_offset = ConfigLoader.atomic_offset()
    atomic_size = ConfigLoader.atomic_size()

    assert atomic_offset >= 0
    assert atomic_size > 0

    iram_offset = ConfigLoader.iram_offset()
    iram_size = ConfigLoader.iram_size()

    assert iram_offset >= 0
    assert iram_size > 0

    wram_offset = ConfigLoader.wram_offset()
    wram_size = ConfigLoader.wram_size()

    assert wram_offset >= 0
    assert wram_size > 0

    mram_offset = ConfigLoader.mram_offset()
    mram_size = ConfigLoader.mram_size()

    assert mram_offset >= 0
    assert mram_size > 0

    assert not is_overlap(atomic_offset, atomic_size, iram_offset, iram_size)
    assert not is_overlap(atomic_offset, atomic_size, wram_offset, wram_size)
    assert not is_overlap(atomic_offset, atomic_size, mram_offset, mram_size)
    assert not is_overlap(iram_offset, iram_size, wram_offset, wram_size)
    assert not is_overlap(iram_offset, iram_size, mram_offset, mram_size)
    assert not is_overlap(wram_offset, wram_size, mram_offset, mram_size)
